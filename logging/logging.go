package logging

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.WithFields(logrus.Fields{
	"service": "event-recorder",
})
