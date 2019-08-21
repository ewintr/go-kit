package log_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"dev-git.sentia.com/go/kit/log"
	"dev-git.sentia.com/go/kit/test"
)

type testLogWriter struct {
	Logs []string
}

func newLogWriter() *testLogWriter {
	return &testLogWriter{}
}

func (t *testLogWriter) Write(p []byte) (n int, err error) {
	t.Logs = append(t.Logs, string(p))
	return
}

func (t *testLogWriter) count() int {
	return len(t.Logs)
}

func (t *testLogWriter) last() string {
	if len(t.Logs) == 0 {
		return ""
	}

	return t.Logs[len(t.Logs)-1]
}

func TestGoKit(t *testing.T) {
	t.Run("new-logger", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		logger := log.NewLogger(&buf)
		test.NotZero(t, logger)
	})

	t.Run("info", func(t *testing.T) {
		t.Parallel()

		logWriter := newLogWriter()
		logger := log.NewLogger(logWriter)
		test.NotZero(t, logger)
		test.Equals(t, 0, logWriter.count())

		msg := "log this"
		test.OK(t, logger.Info(msg))
		test.Equals(t, 1, logWriter.count())
		testLogLine(t, false, msg, logWriter.last())

		msg = "log again"
		test.OK(t, logger.Info(msg))
		test.Equals(t, 2, logWriter.count())
		testLogLine(t, false, msg, logWriter.last())
	})

	t.Run("debug", func(t *testing.T) {
		t.Parallel()

		logWriter := newLogWriter()
		logger := log.NewLogger(logWriter)
		test.NotZero(t, logger)

		// starts with debug disabled
		test.Equals(t, false, logger.DebugStatus())

		msg := "log this"
		logger.DebugEnabled(true)
		test.Equals(t, true, logger.DebugStatus())
		logger.Debug(msg)
		test.Equals(t, 1, logWriter.count())
		testLogLine(t, true, msg, logWriter.last())

		msg = "log again"
		logger.DebugEnabled(false)
		test.Equals(t, false, logger.DebugStatus())
		logger.Debug(msg)
		test.Equals(t, 1, logWriter.count())
	})

	t.Run("normalize-string", func(t *testing.T) {
		t.Parallel()

		missingMsg := "(MISSING)"
		for _, tc := range []struct {
			context    string
			logMessage string
			expResult  string
		}{
			{
				context:    "empty string",
				logMessage: "",
				expResult:  missingMsg,
			},
			{
				context:    "single whitespace",
				logMessage: " ",
				expResult:  missingMsg,
			},
			{
				context:    "multiple whitespace",
				logMessage: "\t ",
				expResult:  missingMsg,
			},
			{
				context:    "simple line",
				logMessage: "just some text",
				expResult:  "just some text",
			},
			{
				context:    "multiline message",
				logMessage: "one\ntwo\nthree",
				expResult:  "one two three",
			},
		} {
			t.Log(tc.context)

			logWriter := newLogWriter()
			logger := log.NewLogger(logWriter)
			test.NotZero(t, logger)
			logger.DebugEnabled(true)

			test.OK(t, logger.Info(tc.logMessage))
			testLogLine(t, false, tc.expResult, logWriter.last())

			test.OK(t, logger.Debug(tc.logMessage))
			testLogLine(t, true, tc.expResult, logWriter.last())
		}
	})

	t.Run("add-context", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			message      string
			contextKey   string
			contextValue interface{}
		}{
			{
				message:      "empty key",
				contextValue: "value",
			},
			{
				message:    "empty value",
				contextKey: "a key",
			},

			{
				message:      "not empty",
				contextKey:   "a key",
				contextValue: "value",
			},
			{
				message:      "slice",
				contextKey:   "a key",
				contextValue: []string{"one", "two"},
			},
			{
				message: "map",
				contextValue: map[string]interface{}{
					"key1": "value1",
					"key2": 2,
					"key3": nil,
				},
			},
			{
				message: "struct",
				contextValue: struct {
					Key1 string
					Key2 int
					Key3 interface{}
					key4 string
				}{
					Key1: "value1",
					Key2: 2,
					Key3: nil,
					key4: "unexported",
				},
			},
		} {
			t.Log(tc.message)

			logWriter := newLogWriter()
			logger := log.NewLogger(logWriter)
			test.NotZero(t, logger)

			test.OK(t, logger.AddContext(tc.contextKey, tc.contextValue).Info("log message"))
			test.Equals(t, 1, logWriter.count())
			testContext(t, tc.contextKey, tc.contextValue, logWriter.last())

		}
	})
}

func testLogLine(t *testing.T, debug bool, message, line string) {
	test.Equals(t, true, isValidJSON(line))
	test.Includes(t, `"time":`, line)
	test.Includes(t, fmt.Sprintf(`"message":%q`, message), line)

	if debug {
		test.Includes(t, `"debug":true`, line)
		test.Includes(t, `"caller":"`, line)
	}

}

func testContext(t *testing.T, key string, value interface{}, line string) {
	value, err := toJSON(value)
	test.OK(t, err)
	test.Includes(t, fmt.Sprintf(`%q:%s`, key, value), line)
}

func toJSON(i interface{}) (string, error) {
	j, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func isValidJSON(s string) bool {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(s), &js)
	if err != nil {
		return false
	}

	return true
}
