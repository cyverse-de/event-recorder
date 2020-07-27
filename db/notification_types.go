package db

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

// GetNotificationTypeID obtains the ID of the notification type with the given name. An error
// is returned if the database can't be queried or the notification type doesn't exist.
func GetNotificationTypeID(tx *sql.Tx, notificationType string) (string, error) {
	wrapMsg := fmt.Sprintf("unable to get the notification type ID for `%s`", notificationType)

	// Build the SQL query and arguments.
	query, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id::text").
		From("notification_types").
		Where(sq.Eq{"name": notificationType}).
		ToSql()
	if err != nil {
		return "", errors.Wrap(err, wrapMsg)
	}

	// Query the database.
	var id string
	row := tx.QueryRow(query, args...)
	err = row.Scan(&id)
	if err != nil {
		return "", errors.Wrap(err, wrapMsg)
	}

	return id, nil
}
