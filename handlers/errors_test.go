package handlers

import (
	"testing"
)

func TestRecoverableError(t *testing.T) {
	var err error
	err = NewRecoverableError("this is a test %s", "of the Emergency Broadcast System")

	// Verify that we go the expected error message.
	if err.Error() != "this is a test of the Emergency Broadcast System" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	// Verify that the type is still considered to be an error.
	_, ok := err.(error)
	if !ok {
		t.Errorf("RecoverableError doesn't appear to be an error")
	}

	// Verify that a RecoverableError was actually returned.
	_, ok = err.(RecoverableError)
	if !ok {
		t.Errorf("The error doesn't appear to be a RecoverableError")
	}

	// The type must be distinct from an uncrecoverable error.
	_, ok = err.(UnrecoverableError)
	if ok {
		t.Errorf("The error appears to be an UnrecoverableError")
	}
}

func TestUnrecoverableError(t *testing.T) {
	var err error
	err = NewUnrecoverableError("testing %s %s", "check", "1...2...3")

	// Verify that w get the expected error message.
	if err.Error() != "testing check 1...2...3" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	// Verify that the type is still considered to be an error.
	_, ok := err.(error)
	if !ok {
		t.Errorf("UnrecoverableError doesn't apper to be an error")
	}

	// Verify that an UnrecoverableError was actually returned.
	_, ok = err.(UnrecoverableError)
	if !ok {
		t.Errorf("The error doesn't appear to be an UnrecoverableError")
	}

	// The type must be distinct from a RecoverableError
	_, ok = err.(RecoverableError)
	if ok {
		t.Errorf("The error appears to be a RecoverableError")
	}
}
