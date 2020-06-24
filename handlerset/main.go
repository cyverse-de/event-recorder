package handlerset

import (
	"github.com/cyverse-de/event-recorder/handlers"
	"github.com/cyverse-de/messaging"
	"github.com/pkg/errors"
)

// AMQPSettings represents the settings that we require in order to connect to the AMQP exchange.
type AMQPSettings struct {
	URI          string
	ExchangeName string
	ExchangeType string
}

// HandlerSet represents a set of AMQP message handlers.
type HandlerSet struct {
	amqpClient *messaging.Client
	handlerFor map[string]handlers.MessageHandler
}

// New creates a new handler set.
func New(amqpSettings *AMQPSettings, handlerFor map[string]handlers.MessageHandler) (*HandlerSet, error) {
	wrapMsg := "unable to create the message handler set"

	// Create the AMQP client.
	amqpClient, err := messaging.NewClient(amqpSettings.URI, false)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Build and return the handler set.
	handlerSet := HandlerSet{
		amqpClient: amqpClient,
		handlerFor: handlerFor,
	}
	return &handlerSet, nil
}

// Close closes a message handler set.
func (hs *HandlerSet) Close() {
	hs.amqpClient.Close()
}
