package log_test

import (
	"errors"
	"testing"

	"git.ewintr.nl/go-kit/log"
	"git.ewintr.nl/go-kit/test"
)

func TestTestLogger(t *testing.T) {
	message := "test line"
	out := log.NewTestOut()
	logger := log.NewTestLogger(out)

	for _, tc := range []struct {
		name    string
		logfunc func(string)
		exp     log.TestLine
	}{
		{
			name:    "debug",
			logfunc: logger.Debug,
			exp: log.TestLine{
				Level:   log.LevelDebug,
				Message: message,
				Fields:  log.Fields{},
			},
		},
		{
			name:    "info",
			logfunc: logger.Info,
			exp: log.TestLine{
				Level:   log.LevelInfo,
				Message: message,
				Fields:  log.Fields{},
			},
		},
		{
			name:    "error",
			logfunc: logger.Error,
			exp: log.TestLine{
				Level:   log.LevelError,
				Message: message,
				Fields:  log.Fields{},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out.Flush()
			tc.logfunc(message)

			test.Equals(t, 1, len(out.Lines))
			test.Equals(t, tc.exp, out.Lines[0])
		})
	}
}

func TestTestLoggerWithField(t *testing.T) {
	key, value := "key", "value"
	message := "message"
	out := log.NewTestOut()
	logger := log.NewTestLogger(out).WithField(key, value)
	logger.Info(message)

	test.Equals(t, 1, len(out.Lines))
	test.Equals(t, log.TestLine{
		Level:   log.LevelInfo,
		Fields:  log.Fields{key: value},
		Message: message,
	}, out.Lines[0])
}

func TestTestLoggerWithErr(t *testing.T) {
	err := errors.New("some err")
	message := "message"
	out := log.NewTestOut()
	logger := log.NewTestLogger(out).WithErr(err)
	logger.Error(message)

	test.Equals(t, 1, len(out.Lines))
	test.Equals(t, log.TestLine{
		Level:   log.LevelError,
		Fields:  log.Fields{"error": err},
		Message: message,
	}, out.Lines[0])
}

func TestTestLoggerWith(t *testing.T) {
	key, value := "key", "value"
	message := "message"
	out := log.NewTestOut()
	logger := log.NewTestLogger(out).With(log.Fields{key: value})
	logger.Info(message)

	test.Equals(t, 1, len(out.Lines))
	test.Equals(t, log.TestLine{
		Level:   log.LevelInfo,
		Fields:  log.Fields{key: value},
		Message: message,
	}, out.Lines[0])
}
