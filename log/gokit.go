package log

import (
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"

	kitlog "github.com/go-kit/kit/log"
)

type GoKit struct {
	debugEnabled bool
	info         kitlog.Logger
	debug        kitlog.Logger
}

func caller(depth int) kitlog.Valuer {
	return func() interface{} {
		_, file, line, _ := runtime.Caller(depth)
		return file + ":" + strconv.Itoa(line)
	}
}

func newGoKitLogger(logWriter io.Writer) Logger {
	w := kitlog.NewSyncWriter(logWriter)
	t := kitlog.TimestampFormat(time.Now, time.RFC3339)
	info := kitlog.With(kitlog.NewJSONLogger(w), "time", t)
	debug := kitlog.With(info, "debug", true, "caller", caller(4))

	return &GoKit{
		info:  info,
		debug: debug,
	}
}

// AddContext attaches a key-value information to the log message
func (gk *GoKit) AddContext(contextKey string, contextValue interface{}) Logger {
	return &GoKit{
		debugEnabled: gk.debugEnabled,
		info:         kitlog.With(gk.info, contextKey, contextValue),
		debug:        kitlog.With(gk.debug, contextKey, contextValue),
	}
}

// Info writes out the log message
func (gk *GoKit) Info(message string) error {
	return gk.info.Log("message", normalizeString(message))
}

// Debug writes out the log message when debug is enabled
func (gk *GoKit) Debug(message string) error {
	if gk.debugEnabled {
		return gk.debug.Log("message", normalizeString(message))
	}
	return nil
}

// DebugEnabled sets debug flag to enable or disabled
func (gk *GoKit) DebugEnabled(enable bool) {
	gk.debugEnabled = enable
}

// DebugStatus returns whether or not debug is enabled
func (gk *GoKit) DebugStatus() bool {
	return gk.debugEnabled
}

func normalizeString(s string) string {
	ss := strings.Fields(s)
	if len(ss) == 0 {
		return "(MISSING)"
	}
	return strings.Join(ss, " ")
}
