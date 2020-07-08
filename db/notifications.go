package db

import (
	"database/sql"

	"github.com/cyverse-de/event-recorder/model"
	"github.com/cyverse-de/logcabin"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

// SaveNotification sabves a single notification into the database.
func SaveNotification(tx *sql.Tx, notification *model.Notification) error {
	wrapMsg := "unable to save notification"

	// Get the notification type ID.
	notificationTypeID, err := GetNotificationTypeID(tx, notification.NotificationType)
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}

	// Get the user ID.
	userID, err := GetUserID(tx, notification.User)
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}

	// Build the statement to insert the notifications.
	statement, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("notifications").
		Columns("notification_type_id", "user_id", "subject", "seen", "deleted", "time_created", "message").
		Values(
			notificationTypeID,
			userID,
			notification.Subject,
			notification.Seen,
			notification.Deleted,
			notification.TimeCreated,
			notification.Message).
		Suffix("RETURNING id").
		ToSql()
	logcabin.Error.Println(statement)
	logcabin.Error.Println(args)
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}

	// Execute the insert statement, scanning the ID into the notification structure.
	row := tx.QueryRow(statement, args...)
	err = row.Scan(&notification.ID)
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}

	return nil
}
