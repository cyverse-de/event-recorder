package handlers

import (
	"encoding/json"

	"github.com/cyverse-de/logcabin"
	"github.com/streadway/amqp"
)

// LegacyRequest represents a deserialized request for a backwards compatible notification.
type LegacyRequest struct {
	RequestType   string                 `json:"type"`
	User          string                 `json:"user"`
	Subject       string                 `json:"subject"`
	Timestamp     string                 `json:"timestamp"`
	Email         bool                   `json:"email"`
	EmailTemplate string                 `json:"email_template"`
	Payload       map[string]interface{} `json:"payload"`
	Message       string                 `json:"message"`
}

// Legacy is a message handler for events published by the backwards compatible HTTP API.
type Legacy struct{}

// NewLegacy returns a new legacy event handler.
func NewLegacy() *Legacy {
	return &Legacy{}
}

// HandleMessage handles a single AMQP delivery.
func (lh *Legacy) HandleMessage(updateType string, delivery amqp.Delivery) error {

	// Parse the message body.
	var request LegacyRequest
	err := json.Unmarshal(delivery.Body, &request)
	if err != nil {
		return NewUnrecoverableError("unable to parse message body: %s", err.Error())
	}

	// Just log the update type for now.
	logcabin.Error.Printf("handling message: %+v", request)

	return nil
}
