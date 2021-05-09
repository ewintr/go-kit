package log_test

import (
	"errors"
	"fmt"
	"testing"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/go-kit/test"
)

func TestLogrusLogger(t *testing.T) {
	logline := "test line"
	tw := &testWriter{}
	logger := log.NewLogrusLogger(tw)
	logger.SetLogLevel(log.LevelDebug)

	for _, tc := range []struct {
		name    string
		logfunc func(string)
	}{
		{
			name:    "debug",
			logfunc: logger.Debug,
		},
		{
			name:    "info",
			logfunc: logger.Info,
		},
		{
			name:    "error",
			logfunc: logger.Error,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tw.Flush()
			tc.logfunc(logline)

			test.Includes(t, logline, tw.LogLines...)
		})
	}
}

func TestLogrusLoggerSetLogLevel(t *testing.T) {
	loglines := map[log.LogLevel]string{
		log.LevelDebug: "debug",
		log.LevelInfo:  "info",
		log.LevelError: "error",
	}
	tw := &testWriter{}
	logger := log.NewLogrusLogger(tw)

	for _, tc := range []struct {
		name  string
		level log.LogLevel
		exp   []string
	}{
		{
			name:  "debug",
			level: log.LevelDebug,
			exp: []string{
				loglines[log.LevelDebug],
				loglines[log.LevelInfo],
				loglines[log.LevelError],
			},
		},
		{
			name:  "info",
			level: log.LevelInfo,
			exp: []string{
				loglines[log.LevelInfo],
				loglines[log.LevelError],
			},
		},
		{
			name:  "error",
			level: log.LevelError,
			exp: []string{
				loglines[log.LevelError],
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tw.Flush()
			logger.SetLogLevel(tc.level)

			logger.Debug(loglines[log.LevelDebug])
			logger.Info(loglines[log.LevelInfo])
			logger.Error(loglines[log.LevelError])

			test.Equals(t, len(tc.exp), len(tw.LogLines))
			for _, ll := range tc.exp {
				test.Includes(t, ll, tw.LogLines...)
			}
		})
	}
}

func TestLogrusLoggerWithField(t *testing.T) {
	tw := &testWriter{}
	logger := log.NewLogrusLogger(tw)
	logger.SetLogLevel(log.LevelDebug)

	// the following only tests whether the shortcut to With() works
	// extensive testing of the fields in combination with levels is
	// in the TestLogrusLoggerWith test
	for _, tc := range []struct {
		name  string
		value interface{}
	}{
		{
			name:  "string",
			value: "value",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tw.Flush()
			key, message := "key", "message"
			fieldLogger := logger.WithField(key, tc.value)

			fieldLogger.Info(message)

			test.Equals(t, 1, len(tw.LogLines))
			test.Includes(t, message, tw.LogLines[0])
			test.Includes(t, key, tw.LogLines[0])
			test.Includes(t, fmt.Sprintf("%v", tc.value), tw.LogLines[0])

		})
	}
}
func TestLogrusLoggerWithErr(t *testing.T) {
	tw := &testWriter{}
	logger := log.NewLogrusLogger(tw)
	logger.SetLogLevel(log.LevelDebug)

	// the following only tests whether the shortcut to With() works
	// extensive testing of the fields in combination with levels is
	// in the TestLogrusLoggerWith test
	for _, tc := range []struct {
		name  string
		value error
	}{
		{
			name:  "string",
			value: errors.New("value"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tw.Flush()
			key, message := "key", "message"
			fieldLogger := logger.WithField(key, tc.value)

			fieldLogger.Info(message)

			test.Equals(t, 1, len(tw.LogLines))
			test.Includes(t, message, tw.LogLines[0])
			test.Includes(t, key, tw.LogLines[0])
			test.Includes(t, tc.value.Error(), tw.LogLines[0])

		})
	}
}

func TestLogrusLoggerWith(t *testing.T) {
	tw := &testWriter{}
	logger := log.NewLogrusLogger(tw)
	logger.SetLogLevel(log.LevelDebug)

	for _, tc := range []struct {
		name  string
		value interface{}
	}{
		{
			name:  "string",
			value: "value",
		},
		{
			name:  "int",
			value: 3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			key := "key"
			fieldLogger := logger.With(log.Fields{
				key: tc.value,
			})

			for _, lf := range []struct {
				name string
				log  func(string)
				flog func(string)
			}{
				{
					name: "debug",
					log:  logger.Debug,
					flog: fieldLogger.Debug,
				},
				{
					name: "info",
					log:  logger.Info,
					flog: fieldLogger.Info,
				},
				{
					name: "error",
					log:  logger.Error,
					flog: fieldLogger.Error,
				},
			} {
				t.Run(lf.name, func(t *testing.T) {
					tw.Flush()
					message := "normal"
					fieldMessage := "field"

					lf.log(message)
					lf.flog(fieldMessage)

					test.Equals(t, 2, len(tw.LogLines))
					// first line is normal logger
					test.Includes(t, message, tw.LogLines[0])
					test.NotIncludes(t, key, tw.LogLines[0])
					// second line is logger with fields
					test.Includes(t, fieldMessage, tw.LogLines[1])
					test.Includes(t, key, tw.LogLines[1])
					test.Includes(t, fmt.Sprintf("%v", tc.value), tw.LogLines[1])
				})
			}
		})
	}
}
