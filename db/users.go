package db

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// AddUser adds a user to the `users` table in the notifications database, returning
// the ID assigned to the user.
func AddUser(ctx context.Context, tx *sql.Tx, user string) (string, error) {
	wrapMsg := fmt.Sprintf("unable to add `%s` to the users table", user)

	// Build the query.
	statement, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("users").Columns("username").
		Values(user).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return "", errors.Wrap(err, wrapMsg)
	}

	// Execute the statement.
	var id string
	row := tx.QueryRowContext(ctx, statement, args...)
	err = row.Scan(&id)
	if err != nil {
		return "", errors.Wrap(err, wrapMsg)
	}

	return id, nil
}

// GetUserID obtains the user ID for `user`, adding the user to the `users` table in
// the notifications database if necessary.
func GetUserID(ctx context.Context, tx *sql.Tx, user string) (string, error) {
	wrapMsg := fmt.Sprintf("unable to get the user ID for `%s`", user)

	// Build the query.
	statement, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id").From("users").
		Where(sq.Eq{"username": user}).
		ToSql()
	if err != nil {
		return "", errors.Wrap(err, wrapMsg)
	}

	// Query the databse.
	var id string
	row := tx.QueryRowContext(ctx, statement, args...)
	err = row.Scan(&id)

	// If The error is nil then we've got the ID already.
	if err == nil {
		return id, nil
	}

	// If the error is ErrNoRows then we need to add the user to the database.
	if err == sql.ErrNoRows {
		return AddUser(ctx, tx, user)
	}

	// If we get here then all we can do is return the error.
	return "", errors.Wrap(err, wrapMsg)
}
