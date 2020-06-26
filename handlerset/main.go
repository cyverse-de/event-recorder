package handlerset

import (
	"fmt"

	"github.com/cyverse-de/event-recorder/handlers"
	"github.com/cyverse-de/logcabin"
	"github.com/cyverse-de/messaging"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const queueName = "event_listener"
const queueKey = "events.*.update.*"

// AMQPSettings represents the settings that we require in order to connect to the AMQP exchange.
type AMQPSettings struct {
	URI          string
	ExchangeName string
	ExchangeType string
}

// HandlerSet represents a set of AMQP message handlers.
type HandlerSet struct {
	amqpClient   *messaging.Client
	amqpSettings *AMQPSettings
	handlerFor   map[string]handlers.MessageHandler
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
		amqpClient:   amqpClient,
		amqpSettings: amqpSettings,
		handlerFor:   handlerFor,
	}
	return &handlerSet, nil
}

// handleMessage handles an incoming AMQP message.
func (hs *HandlerSet) handle(delivery amqp.Delivery) {
	// For now, we're just giong to exclaim that we got the message then acknowledge it.
	fmt.Println("We can haz messagez!!!!")
	err := delivery.Ack(false)
	if err != nil {
		logcabin.Error.Printf("unable to acknowledge delivery: %s", err.Error())
	}
}

// Listen waits for incoming AMQP messages and dispatches any messages that it recieves to a handler.
func (hs *HandlerSet) Listen() error {
	wrapMsg := "error encountered while listening for incoming events"

	// Set up publishing on the AMQP client in case we need to publish any messages.
	err := hs.amqpClient.SetupPublishing(hs.amqpSettings.ExchangeName)
	if err != nil {
		return errors.Wrap(err, wrapMsg)
	}

	// Start listening for incoming messages.
	go hs.amqpClient.Listen()

	// Listen for incoming messages.
	hs.amqpClient.AddConsumer(
		hs.amqpSettings.ExchangeName,
		hs.amqpSettings.ExchangeType,
		queueName,
		queueKey,
		hs.handle,
		100,
	)

	return nil
}

// Close closes a message handler set.
func (hs *HandlerSet) Close() {
	hs.amqpClient.Close()
}
