package log

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	fields Fields
	logger *logrus.Logger
}

func NewLogrusLogger(out io.Writer) Logger {
	lr := &logrus.Logger{
		Out:       out,
		Formatter: new(logrus.JSONFormatter),
		Level:     logrus.InfoLevel,
	}

	return &LogrusLogger{
		fields: make(Fields),
		logger: lr,
	}
}

func (ll *LogrusLogger) SetLogLevel(loglevel LogLevel) {
	switch LogLevel(loglevel) {
	case LevelDebug:
		ll.logger.SetLevel(logrus.DebugLevel)
	case LevelInfo:
		ll.logger.SetLevel(logrus.InfoLevel)
	case LevelError:
		ll.logger.SetLevel(logrus.ErrorLevel)
	default:
		ll.logger.SetLevel(logrus.InfoLevel)
	}
}

func (ll *LogrusLogger) WithField(key string, value interface{}) Logger {
	return ll.With(Fields{
		key: value,
	})
}

func (ll *LogrusLogger) WithErr(err error) Logger {
	return ll.With(Fields{
		"error": err,
	})
}

func (ll *LogrusLogger) With(fields Fields) Logger {
	newFields := make(Fields)
	for k, v := range ll.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &LogrusLogger{
		fields: newFields,
		logger: ll.logger,
	}
}

func (ll *LogrusLogger) Debug(message string) {
	ll.logger.WithFields(logrus.Fields(ll.fields)).Debug(message)
}

func (ll *LogrusLogger) Info(message string) {
	ll.logger.WithFields(logrus.Fields(ll.fields)).Info(message)
}

func (ll *LogrusLogger) Error(message string) {
	ll.logger.WithFields(logrus.Fields(ll.fields)).Error(fmt.Errorf(message))
}
