// Package log implements a generic interface to log
package log

import (
	"encoding/json"
	"io"
)

// Logger represents a log implementation
type Logger interface {
	AddContext(string, interface{}) Logger
	Info(string) error
	Debug(string) error
	DebugEnabled(bool)
	DebugStatus() bool
}

// NewLogger returns a Logger implementation
func NewLogger(logWriter io.Writer) Logger {
	return newGoKitLogger(logWriter)
}

type loggerWriter struct {
	Logger
}

func (l *loggerWriter) Write(p []byte) (n int, err error) {
	var fields map[string]interface{}

	if err = json.Unmarshal(p, &fields); err != nil {
		l.Logger.Info(string(p))
		return
	}

	delete(fields, "time")

	var message string
	if m, ok := fields["message"]; ok {
		message = m.(string)
		delete(fields, "message")
	}

	if len(fields) == 0 {
		l.Logger.Info(message)
		return
	}

	l.Logger.AddContext("fields", fields).Info(message)
	return
}

// NewWriter returns io.Writer implementation based on a logger
func NewWriter(l Logger) io.Writer {
	return &loggerWriter{l}
}
