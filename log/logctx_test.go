package log_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"dev-git.sentia.com/go/kit/log"
	"dev-git.sentia.com/go/kit/test"
)

type (
	ContextA string
	ContextB string
)

func (l ContextA) ContextName() string { return "context_a" }
func (l ContextB) ContextName() string { return "context_b" }

func TestLogContext(t *testing.T) {
	t.Run("new caller", func(t *testing.T) {
		caller := log.NewCaller(1)
		s := caller()
		test.Includes(t, "logctx_test.go:", s)
	})

	t.Run("add context", func(t *testing.T) {
		var buff bytes.Buffer
		logger := log.NewLogger(&buff)

		for _, tc := range []struct {
			m  string
			cc []log.Contexter
		}{
			{
				m:  "single context",
				cc: []log.Contexter{ContextA("AA")},
			},
			{
				m:  "multiple context",
				cc: []log.Contexter{ContextA("AA"), ContextB("BB")},
			},
			{
				m:  "with caller context",
				cc: []log.Contexter{ContextA("AA"), log.NewCaller(0)},
			},
		} {
			t.Run(tc.m, func(t *testing.T) {
				log.Add(logger, tc.cc...).Info("something")
				for _, context := range tc.cc {
					switch s := context.(type) {
					case ContextA, ContextB:
						test.Includes(t, fmt.Sprintf("%q:%q", s.ContextName(), s), buff.String())

					case log.Caller:
						file := s()
						i := strings.LastIndexByte(file, ':')
						test.Includes(t, fmt.Sprintf(`%q:"%s`, s.ContextName(), file[:i+1]), buff.String())

					}
				}
			})
		}
	})

	t.Run("attach error", func(t *testing.T) {

		var (
			errOne = fmt.Errorf("error one")
			errTwo = fmt.Errorf("error two")
		)

		for _, tc := range []struct {
			m    string
			errs []error
			err  error
		}{
			{
				m:    "single call",
				errs: []error{errOne},
				err:  errOne,
			},
			{
				m:    "multiple calls overwrite",
				errs: []error{errOne, errTwo},
				err:  errTwo,
			},
		} {
			t.Run(tc.m, func(t *testing.T) {
				var buff bytes.Buffer
				logger := log.NewLogger(&buff)

				currentLogger := logger
				for _, err := range tc.errs {
					currentLogger = log.AttachError(currentLogger, err)
				}
				currentLogger.Info("something")

				test.Includes(t,
					fmt.Sprintf("\"attached_error\":%q", tc.err.Error()), buff.String())
			})
		}
	})
}
