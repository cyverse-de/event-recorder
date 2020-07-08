package model

import "time"

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
