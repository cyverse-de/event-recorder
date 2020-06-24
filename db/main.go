package db

import (
	"database/sql"

	"github.com/cyverse-de/dbutil"
	"github.com/pkg/errors"
)

// InitDatabase establishes a database connection and verifies tha the database can be reached.
func InitDatabase(driverName, databaseURI string) (*sql.DB, error) {
	wrapMsg := "unable to initialize the database"

	// Create a database connector to establish the connection.
	connector, err := dbutil.NewDefaultConnector("1m")
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Establish the database connection.
	db, err := connector.Connect(driverName, databaseURI)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return db, nil
}
