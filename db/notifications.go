package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cyverse-de/event-recorder/common"
	"github.com/pkg/errors"
	"gopkg.in/cyverse-de/messaging.v7"

	sq "github.com/Masterminds/squirrel"
)

// SaveNotification sabves a single notification into the database.
func SaveNotification(tx *sql.Tx, notification *common.Notification) error {
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
		Columns(
			"notification_type_id",
			"user_id",
			"subject",
			"seen",
			"deleted",
			"time_created",
			"incoming_json",
			"routing_key").
		Values(
			notificationTypeID,
			userID,
			notification.Subject,
			notification.Seen,
			notification.Deleted,
			notification.TimeCreated,
			notification.Message,
			notification.RoutingKey).
		Suffix("RETURNING id").
		ToSql()
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

// SaveOutgoingNotification adds the outgoing notification JSON to the notification in the database.
func SaveOutgoingNotification(tx *sql.Tx, outgoingNotification *messaging.NotificationMessage) error {
	wrapMsg := "unable to save outgoing notification JSON"

	// Marshal the outgoing notification message.
	outgoingJSON, err := json.Marshal(outgoingNotification)
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}

	// Build the statement to add the notification.
	statement, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("notifications").
		Set("outgoing_json", outgoingJSON).
		Where(sq.Eq{"id": outgoingNotification.Message["id"].(string)}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}

	// Execute the update statement and verify that the correct number of rows was affected.
	result, err := tx.Exec(statement, args...)
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("%s: unexpected number of rows affected: %d", wrapMsg, rowsAffected)
	}

	return nil
}
