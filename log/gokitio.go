package log

import (
	"io"

	kitlog "github.com/go-kit/kit/log"
)

type GoKitIOLogger struct {
	fields Fields
	level  LogLevel
	logger kitlog.Logger
}

func NewGoKitIOLogger(out io.Writer) Logger {
	kl := kitlog.NewJSONLogger(out)
	kl = kitlog.With(kl, "time", kitlog.DefaultTimestampUTC)

	return &GoKitIOLogger{
		fields: make(Fields),
		level:  LevelInfo,
		logger: kl,
	}
}

func (kl *GoKitIOLogger) SetLogLevel(loglevel LogLevel) {
	kl.level = loglevel
}

func (kl *GoKitIOLogger) WithField(key string, value interface{}) Logger {
	return kl.With(Fields{
		key: value,
	})
}

func (kl *GoKitIOLogger) WithErr(err error) Logger {
	return kl.With(Fields{
		"error": err,
	})
}

func (kl *GoKitIOLogger) With(fields Fields) Logger {
	newFields := make(Fields)
	for k, v := range kl.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &GoKitIOLogger{
		fields: newFields,
		level:  kl.level,
		logger: kl.logger,
	}
}

func (kl *GoKitIOLogger) Debug(message string) {
	if kl.level == LevelDebug {
		kl.log("debug", message)
	}
}

func (kl *GoKitIOLogger) Info(message string) {
	if kl.level != LevelError {
		kl.log("info", message)
	}
}

func (kl *GoKitIOLogger) Error(message string) {
	kl.log("error", message)
}

func (kl *GoKitIOLogger) log(level, message string) {
	kv := make([]interface{}, 0)
	for k, v := range kl.fields {
		kv = append(kv, k, v)
	}
	kv = append(kv, "level", level, "message", message)

	kl.logger.Log(kv...)
}
