package common

import (
	"time"

	"github.com/mcnijman/go-emailaddress"
)

// AMQPSettings represents the settings that we require in order to connect to the AMQP exchange.
type AMQPSettings struct {
	URI          string
	ExchangeName string
	ExchangeType string
}

// Notification represents a single notification to be recorded in the database.
type Notification struct {
	ID               string
	NotificationType string
	User             string
	Subject          string
	Seen             bool
	Deleted          bool
	TimeCreated      time.Time
	Message          string
}

// ValidateEmailAddress returns an error if the format of an email address is invalid.
func ValidateEmailAddress(emailAddress string) error {
	_, err := emailaddress.Parse(emailAddress)
	return err
}
