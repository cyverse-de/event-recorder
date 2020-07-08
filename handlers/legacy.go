package handlers

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/cyverse-de/event-recorder/db"

	"github.com/cyverse-de/event-recorder/model"

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
type Legacy struct {
	db *sql.DB
}

// NewLegacy returns a new legacy event handler.
func NewLegacy(db *sql.DB) *Legacy {
	return &Legacy{db: db}
}

// HandleMessage handles a single AMQP delivery.
func (lh *Legacy) HandleMessage(updateType string, delivery amqp.Delivery) error {

	// Parse the message body.
	var request LegacyRequest
	err := json.Unmarshal(delivery.Body, &request)
	if err != nil {
		return NewUnrecoverableError("unable to parse message body: %s", err.Error())
	}

	// Parse the timestamp.
	timeCreated, err := time.Parse(time.RFC3339Nano, request.Timestamp)
	if err != nil {
		return NewUnrecoverableError("unable to parse timestamp: %s", err.Error())
	}

	// Begin a database transaction.
	tx, err := lh.db.Begin()
	if err != nil {
		return NewRecoverableError("uanble to begin a database transaction: %s", err.Error())
	}
	defer tx.Rollback()

	// Store the message in the database.
	storableRequest := &model.Notification{
		NotificationType: updateType,
		User:             request.User,
		Subject:          request.Subject,
		Seen:             false,
		Deleted:          false,
		TimeCreated:      timeCreated,
		Message:          string(delivery.Body),
	}
	err = db.SaveNotification(tx, storableRequest)
	if err != nil {
		return NewUnrecoverableError(err.Error())
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return NewRecoverableError("unable to commit the database transaction: %s", err.Error())
	}

	return nil
}
