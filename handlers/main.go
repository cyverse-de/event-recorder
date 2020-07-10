package handlers

import (
	"database/sql"

	"github.com/cyverse-de/event-recorder/common"
	"github.com/cyverse-de/messaging"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// MessageHandler describes the interface used to handle AMQP messages.
type MessageHandler interface {
	HandleMessage(updateType string, delivery amqp.Delivery) error
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

	// Create the messaging client.
	messagingClient, err := createMessagingClient(amqpSettings)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Create the message handlers.
	messageHandlers := map[string]MessageHandler{
		"notification": NewLegacy(db, messagingClient),
	}

	return messageHandlers, nil
}
