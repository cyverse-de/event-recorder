package handlers

import (
	"database/sql"

	"github.com/cyverse-de/event-recorder/db"

	"github.com/cyverse-de/event-recorder/common"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gopkg.in/cyverse-de/messaging.v8"
)

// MessageHandler describes the interface used to handle AMQP messages.
type MessageHandler interface {
	HandleMessage(updateType string, delivery amqp.Delivery) error
}

// MessagingClient is a subset of messaging.Client. Its purpose is to limit the number of mock
// functions we have to implement for unit testing. During unit tests, a mock messaging client
// will be used. Otherwise, messaging.Client will be used directly.
type MessagingClient interface {
	PublishEmailRequest(*messaging.EmailRequest) error
	PublishNotificationMessage(*messaging.NotificationMessage) error
}

// DatabaseClient provides a wrapper around functions that handlers might call in order to interact
// with the database.
type DatabaseClient interface {
	Begin() (*sql.Tx, error)
	Commit(*sql.Tx) error
	Rollback(*sql.Tx) error
	SaveNotification(*sql.Tx, *common.Notification) error
	SaveOutgoingNotification(*sql.Tx, *messaging.NotificationMessage) error
}

// DatabaseClientImpl provides the default implementation of DatabaseClient.
type DatabaseClientImpl struct {
	db *sql.DB
}

// Begin starts a new database transaction.
func (c *DatabaseClientImpl) Begin() (*sql.Tx, error) {
	return c.db.Begin()
}

// Commit commits an existing database transaction.
func (c *DatabaseClientImpl) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

// Rollback rolls back an existing database transaction.
func (c *DatabaseClientImpl) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

// SaveNotification saves a notification in the database.
func (c *DatabaseClientImpl) SaveNotification(tx *sql.Tx, notification *common.Notification) error {
	return db.SaveNotification(tx, notification)
}

// SaveOutgoingNotification adds the outbound notification JSON to the notification record in the database.
func (c *DatabaseClientImpl) SaveOutgoingNotification(
	tx *sql.Tx,
	outgoingNotification *messaging.NotificationMessage,
) error {
	return db.SaveOutgoingNotification(tx, outgoingNotification)
}

// NewDatabaseClient creates a new default database client implementation.
func NewDatabaseClient(db *sql.DB) DatabaseClient {
	return &DatabaseClientImpl{db: db}
}

// createMessagingClient creates a new AMQP messaging client and sets up publishing on that client.
func createMessagingClient(amqpSettings *common.AMQPSettings) (*messaging.Client, error) {
	wrapMsg := "unable to create the messaging client"

	// Create the messaging client.
	client, err := messaging.NewClient(amqpSettings.URI, false)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Set up publishing on the messaging client.
	err = client.SetupPublishing(amqpSettings.ExchangeName)
	if err != nil {
		client.Close()
		return nil, errors.Wrap(err, wrapMsg)
	}

	return client, nil
}

// InitMessageHandlers returns a map from category name to message handler.
func InitMessageHandlers(db *sql.DB, amqpSettings *common.AMQPSettings) (map[string]MessageHandler, error) {
	wrapMsg := "unable to initialize message handlers"

	// Create the database client.
	databaseClient := NewDatabaseClient(db)

	// Create the messaging client.
	messagingClient, err := createMessagingClient(amqpSettings)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Create the message handlers.
	messageHandlers := map[string]MessageHandler{
		"notification": NewLegacy(databaseClient, messagingClient),
	}

	return messageHandlers, nil
}
