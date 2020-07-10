package common

import (
	"testing"
	"time"
)

func GetTestTimestamp() time.Time {
	return time.Unix(int64(1594336370), int64(706917000))
}

func GetTestTimestampMillisecondPrecision() string {
	return "1594336370706"
}

func GetTestTimestampSecondPrecision() string {
	return "1594336370000"
}

func TestFormatTimestamp(t *testing.T) {
	timestamp := GetTestTimestamp()
	expected := GetTestTimestampMillisecondPrecision()
	actual := FormatTimestamp(timestamp)
	if actual != expected {
		t.Errorf("unexpected timestamp: got '%s' instead of '%s'", actual, expected)
	}
}

func TestFixTimestampEmptyString(t *testing.T) {
	actual, err := FixTimestamp("")
	if err != nil {
		t.Errorf("FixTimestamp returned an error: %s", err.Error())
	}
	if actual != "" {
		t.Errorf("FixTimestamp returned '%s' instead of an empty string", actual)
	}
}

func TestFixTimestampMillis(t *testing.T) {
	expected := GetTestTimestampMillisecondPrecision()
	actual, err := FixTimestamp(expected)
	if err != nil {
		t.Errorf("FixTimestamp returned an error: %s", err.Error())
	}
	if actual != expected {
		t.Errorf("FixTimestamp returned '%s' instead of '%s'", actual, expected)
	}
}

func TestFixTimestampRFC3339(t *testing.T) {
	timestamp := GetTestTimestamp()
	original := timestamp.Format(time.RFC3339)
	expected := GetTestTimestampSecondPrecision()
	actual, err := FixTimestamp(original)
	if err != nil {
		t.Errorf("FixTimestamp returned an error: %s", err.Error())
	}
	if actual != expected {
		t.Errorf("FixTimestamp returned '%s' instead of '%s'", actual, expected)
	}
}

func TestFixTimestampRFC3339Nano(t *testing.T) {
	timestamp := GetTestTimestamp()
	original := timestamp.Format(time.RFC3339Nano)
	expected := GetTestTimestampMillisecondPrecision()
	actual, err := FixTimestamp(original)
	if err != nil {
		t.Errorf("FixTimestamp returned an error: %s", err.Error())
	}
	if actual != expected {
		t.Errorf("FixTimestamp returned '%s' instead of '%s'", actual, expected)
	}
}
