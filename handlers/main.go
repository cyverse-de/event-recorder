package handlers

import (
	"database/sql"

	"github.com/streadway/amqp"
)

// MessageHandler describes the interface used to handle AMQP messages.
type MessageHandler interface {
	HandleMessage(updateType string, delivery amqp.Delivery) error
}

// InitMessageHandlers returns a map from category name to message handler.
func InitMessageHandlers(_ *sql.DB) map[string]MessageHandler {
	return map[string]MessageHandler{}
}
