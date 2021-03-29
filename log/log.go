package log

import "io"

const (
	LevelDebug = LogLevel("debug")
	LevelInfo  = LogLevel("info")
	LevelError = LogLevel("error")
)

type LogLevel string

type Fields map[string]interface{}

type Logger interface {
	SetLogLevel(loglevel LogLevel)
	WithField(key string, value interface{}) Logger
	WithErr(err error) Logger
	With(fields Fields) Logger
	Debug(message string)
	Info(message string)
	Error(message string)
}

func New(out io.Writer) Logger {
	return NewLogrusLogger(out)
}
