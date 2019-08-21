package log_test

import (
	"bytes"
	"fmt"

	"dev-git.sentia.com/go/kit/log"
)

// LoggerTestable represents a data structure for a log context
type LoggerTestable struct {
	TestKey string `json:"key"`
}

func (l LoggerTestable) ContextName() string { return "test" }

func Example() {
	var buff bytes.Buffer
	var logger log.Logger

	logger = log.NewLogger(&buff)

	// Please ignore the following line, it was added to allow better
	// assertion of the results when logging.
	logger = logger.AddContext("time", "-")

	logger = log.Add(logger, LoggerTestable{
		TestKey: "value",
	})
	logger.Info("this is an example.")

	fmt.Println(buff.String())
	// Output: {"message":"this is an example.","test":{"key":"value"},"time":"-"}
}
