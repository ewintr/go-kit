package log_test

import (
	"bytes"
	"fmt"
	"testing"

	"git.sr.ht/~ewintr/go-kit/log"
	"git.sr.ht/~ewintr/go-kit/test"
)

func TestNewWriter(t *testing.T) {

	defaultMessage := "this is a test"
	for _, tc := range []struct {
		m           string
		message     string
		expected    []string
		notExpected []string
	}{
		{
			m:       "string input",
			message: defaultMessage,
			expected: []string{
				fmt.Sprintf(`"message":%q`, defaultMessage),
			},
		},
		{
			m:       "json map",
			message: fmt.Sprintf(`{"message":%q, "custom": "value"}`, defaultMessage),
			expected: []string{
				fmt.Sprintf(`"message":%q`, defaultMessage),
				`"fields":{`,
				`"custom":"value"`,
			},
		},
		{
			m:       "json map correct time",
			message: fmt.Sprintf(`{"message":%q, "time": "value"}`, defaultMessage),
			expected: []string{
				fmt.Sprintf(`"message":%q`, defaultMessage),
				`"time":"`,
			},
			notExpected: []string{
				`"fields":{`,
				`"time": "value"`,
			},
		},
	} {
		var buf bytes.Buffer
		logger := log.NewLogger(&buf)

		t.Run(tc.m, func(t *testing.T) {
			w := log.NewWriter(logger)
			w.Write([]byte(tc.message))
			for _, e := range tc.expected {
				test.Includes(t, e, buf.String())
			}
			for _, e := range tc.notExpected {
				test.NotIncludes(t, e, buf.String())
			}
		})
	}
}

func TestLog(t *testing.T) {
	t.Run("new-logger", func(t *testing.T) {
		var buf bytes.Buffer
		logger := log.NewLogger(&buf)
		test.NotZero(t, logger)
	})
}
