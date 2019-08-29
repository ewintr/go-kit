package log

import (
	"runtime"
	"strconv"
)

// Contexter ensures type is intentionally a log context
type Contexter interface {
	ContextName() string
}

// Caller represents a runtime file:line caller for log context
type Caller func() string

// ContextName returns the key for the log context
func (c Caller) ContextName() string { return "caller" }

// NewCaller returns a log context for runtime file caller with full path
func NewCaller(depth int) Caller {
	return func() string {
		_, file, line, _ := runtime.Caller(depth)
		return file + ":" + strconv.Itoa(line)
	}
}

// Add adds a contexter interface to a Logger
func Add(l Logger, cc ...Contexter) Logger {
	for _, c := range cc {
		if caller, ok := c.(Caller); ok {
			l = l.AddContext(c.ContextName(), caller())
			continue
		}
		l = l.AddContext(c.ContextName(), c)
	}
	return l
}
