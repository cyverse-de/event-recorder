package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetNotificationTypeID(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.NoError(err, "unable to open the mock database connection")
	defer db.Close()

	// Set up the expectations.
	mock.ExpectBegin()
	testID := "a6a97fd2-74c5-42af-ab22-0549a63d3abd"
	rows := sqlmock.NewRows([]string{"id"}).AddRow(testID)
	mock.ExpectQuery("SELECT id::text FROM notification_types WHERE name =").
		WithArgs("test").
		WillReturnRows(rows)
	mock.ExpectRollback()

	// Look up a notification type.
	tx, err := db.Begin()
	assert.NoError(err, "unable to begin a transaction")
	id, err := GetNotificationTypeID(tx, "test")
	assert.NoError(err, "unexpected error occurred while looking up the notification type ID")
	assert.Equal(testID, id)
	tx.Rollback()

	// Verify that all mock expectations were met.
	assert.NoError(err, "not all mock expectations were met")
}
