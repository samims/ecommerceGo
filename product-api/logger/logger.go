package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(logLevel logrus.Level) *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logLevel)
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	return l
}
