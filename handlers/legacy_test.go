package handlers

import (
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/cyverse-de/event-recorder/common"
	"github.com/streadway/amqp"
	"gopkg.in/cyverse-de/messaging.v7"
)

// MockMessagingClient provides mock implementations of the functions we need from messaging.Client.
type MockMessagingClient struct {
	PublishedNotificationMessage *messaging.NotificationMessage
	PublishedEmailRequest        *messaging.EmailRequest
}

// PublishNotificationMessage simply stores a copy of the notification message for later inspection.
func (c *MockMessagingClient) PublishNotificationMessage(msg *messaging.NotificationMessage) error {
	c.PublishedNotificationMessage = msg
	return nil
}

// PublishEmailRequest simply stores a copy of the email request for later inspection.
func (c *MockMessagingClient) PublishEmailRequest(req *messaging.EmailRequest) error {
	c.PublishedEmailRequest = req
	return nil
}

// NewMockMessagingClient creates a new mock messaging client for testing.
func NewMockMessagingClient() *MockMessagingClient {
	return &MockMessagingClient{
		PublishedNotificationMessage: nil,
		PublishedEmailRequest:        nil,
	}
}

// FakeNotificationID is the identifier that will be assigned to notifications by the mock database
// client.
const FakeNotificationID = "46ae63be-7030-4cdd-8eb9-66aa49fcf38b"

// MockDatabaseClient provides mock implementations of functions that handlers call to interact with the
// database.
type MockDatabaseClient struct {
	BeginCalled          bool
	CommitCalled         bool
	RollbackCalled       bool
	SavedNotification    *common.Notification
	savedOutgoingMessage *messaging.NotificationMessage
}

// Begin records the fact that it was called.
func (c *MockDatabaseClient) Begin() (*sql.Tx, error) {
	c.BeginCalled = true
	return nil, nil
}

// Commit records the fact that it was called.
func (c *MockDatabaseClient) Commit(*sql.Tx) error {
	c.CommitCalled = true
	return nil
}

// Rollback records the fact that it was called.
func (c *MockDatabaseClient) Rollback(*sql.Tx) error {
	c.RollbackCalled = true
	return nil
}

// SaveNotification records a copy of the notification that was saved.
func (c *MockDatabaseClient) SaveNotification(tx *sql.Tx, notification *common.Notification) error {
	notification.ID = FakeNotificationID
	c.SavedNotification = notification
	return nil
}

// SaveOutgoingNotification records a copy of the notification message that was saved.
func (c *MockDatabaseClient) SaveOutgoingNotification(
	tx *sql.Tx,
	outgoingNotification *messaging.NotificationMessage,
) error {
	c.savedOutgoingMessage = outgoingNotification
	return nil
}

// NewMockDatabaseClient creates a new mock database client for testing.
func NewMockDatabaseClient() *MockDatabaseClient {
	return &MockDatabaseClient{
		BeginCalled:          false,
		CommitCalled:         false,
		RollbackCalled:       false,
		SavedNotification:    nil,
		savedOutgoingMessage: nil,
	}
}

// getLegacyNotificationRequest returns a map that can be used as a request
// for a legacy notification request.
func getLegacyNotificationRequest() map[string]interface{} {
	return map[string]interface{}{
		"type":      "analysis",
		"user":      "sarahr",
		"subject":   "some job status changed",
		"message":   "This is a test message",
		"timestamp": "2020-07-07T17:59:59-07:00",
		"payload": map[string]interface{}{
			"action":                "job_status_change",
			"analysisname":          "some job",
			"analysisdescription":   "some job description",
			"analysisstatus":        "Completed",
			"analysisstartdate":     "2020-07-07T17:59:59-07:00",
			"analysisresultsfolder": "/iplant/home/foo/analyses",
			"description":           "some job description",
			"email_address":         "sarahr@cyverse.org",
			"name":                  "some job",
			"resultfolderid":        "/iplant/home/foo/analyses",
			"startdate":             "2020-07-07T17:59:59-07:00",
			"status":                "Completed",
			"user":                  "sarahr",
		},
		"email_template": "analysis_status_change",
		"email":          true,
	}
}

// timestampFormatCorrect returns true if the format of the timestamp in the
// given message appears to be corect.
func timestampFormatCorrect(timestamp string) bool {
	re := regexp.MustCompile("^\\d+$")
	return re.MatchString(timestamp)
}

func TestNotification(t *testing.T) {

	// Create the AMQP delivery for testing.
	requestBody, err := json.Marshal(getLegacyNotificationRequest())
	if err != nil {
		t.Fatalf("unable to marshal the notification request: %s", err.Error())
	}
	delivery := amqp.Delivery{Body: requestBody}

	// The database and messaging clients along with the handler.
	databaseClient := NewMockDatabaseClient()
	messagingClient := NewMockMessagingClient()
	handler := NewLegacy(databaseClient, messagingClient)

	// Pass the delivery to the handler.
	err = handler.HandleMessage("analysis", delivery)
	if err != nil {
		t.Fatalf("unxpected error returned by legacy handler: %s", err.Error())
	}

	// Verify that a transaction was created and committed.
	if !databaseClient.BeginCalled {
		t.Errorf("no database transaction was started")
	}
	if !databaseClient.CommitCalled {
		t.Errorf("the database transaction was not committed")
	}

	// Verify that a notification was saved and spot-check a couple of fields.
	savedNotification := databaseClient.SavedNotification
	if savedNotification == nil {
		t.Fatalf("no notification was saved")
	}
	notificationType := savedNotification.NotificationType
	if notificationType != "analysis" {
		t.Errorf("incorrect type in notifiation: got '%s'; expected 'analysis'", notificationType)
	}
	user := savedNotification.User
	if user != "sarahr" {
		t.Errorf("incorrect user in notification: got '%s'; expected 'sarahr'", user)
	}

	// Verify that the outgoing notification was saved in the database and spot-check a couple of fields.
	savedOutgoingMessage := databaseClient.savedOutgoingMessage
	if savedOutgoingMessage == nil {
		t.Fatalf("the outbound notification message was not recorded in the database")
	}
	if savedOutgoingMessage.Message["id"] != FakeNotificationID {
		t.Errorf("incorrect ID in notification message: %s", savedOutgoingMessage.Message["id"])
	}
	if !timestampFormatCorrect(savedOutgoingMessage.Message["timestamp"].(string)) {
		t.Errorf("incorrect timestamp format: %s", savedOutgoingMessage.Message["timestamp"].(string))
	}

	// Verify that an email request was sent and spot-check a couple of fields.
	emailRequest := messagingClient.PublishedEmailRequest
	if emailRequest == nil {
		t.Fatalf("no email request was published")
	}
	if emailRequest.Subject != "some job status changed" {
		t.Errorf("incorrect subject in email request: %s", emailRequest.Subject)
	}
	if emailRequest.ToAddress != "sarahr@cyverse.org" {
		t.Errorf("incorrect address in email request: %s", emailRequest.ToAddress)
	}

	// Verify that the notification was published and spot-check a couple of fields.
	notification := messagingClient.PublishedNotificationMessage
	if notification == nil {
		t.Fatalf("no notification was published")
	}
	if notification.Message["id"] != FakeNotificationID {
		t.Errorf("incorrect ID in notification message: %s", notification.Message["id"])
	}
	if !timestampFormatCorrect(notification.Message["timestamp"].(string)) {
		t.Errorf("incorrect timestamp format: %s", notification.Message["timestamp"].(string))
	}

	// Spot-check some fields in the payload.
	payload, ok := notification.Payload.(*LegacyRequest)
	if !ok {
		t.Fatal("payload doesn't appear to be a LegacyRequest")
	}
	if !timestampFormatCorrect(payload.Payload["startdate"].(string)) {
		t.Errorf("incorrect timestamp format: %s", payload.Payload["startdate"].(string))
	}
	_, ok = payload.Payload["enddate"]
	if ok {
		t.Error("enddate was found in the payload when it wasn't expected")
	}
}

func TestNotificationWithoutEmail(t *testing.T) {

	// Disable emails for this notification.
	req := getLegacyNotificationRequest()
	req["email"] = false

	// Create the AMQP delivery for testing.
	requestBody, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("unable to marshal the notification request: %s", err.Error())
	}
	delivery := amqp.Delivery{Body: requestBody}

	// The database and messaging clients along with the handler.
	databaseClient := NewMockDatabaseClient()
	messagingClient := NewMockMessagingClient()
	handler := NewLegacy(databaseClient, messagingClient)

	// Pass the delivery to the handler.
	err = handler.HandleMessage("analysis", delivery)
	if err != nil {
		t.Fatalf("unxpected error returned by legacy handler: %s", err.Error())
	}

	// Verify that a transaction was created and committed.
	if !databaseClient.BeginCalled {
		t.Errorf("no database transaction was started")
	}
	if !databaseClient.CommitCalled {
		t.Errorf("the database transaction was not committed")
	}

	// Verify that a notification was saved.
	savedNotification := databaseClient.SavedNotification
	if savedNotification == nil {
		t.Fatalf("no notification was saved")
	}

	// Verify that an email request was not sent.
	emailRequest := messagingClient.PublishedEmailRequest
	if emailRequest != nil {
		t.Fatalf("an email request was published when none was expected")
	}

	// Verify that the notification was published and spot-check a couple of fields.
	notification := messagingClient.PublishedNotificationMessage
	if notification == nil {
		t.Fatalf("no notification was published")
	}
}
