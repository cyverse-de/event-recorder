package handlers

import "database/sql"

// MessageHandler describes the interface used to handle AMQP messages.
type MessageHandler interface {
	HandleMessage(messageBody map[string]interface{}) error
}

// InitMessageHandlers returns a map from category name to message handler.
func InitMessageHandlers(_ *sql.DB) map[string]MessageHandler {
	return map[string]MessageHandler{}
}
