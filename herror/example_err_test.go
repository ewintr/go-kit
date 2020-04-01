package herror_test

import (
	"fmt"

	"git.sr.ht/~ewintr/go-kit/herror"
)

var ErrTaskFailed = herror.New("task has failed")

func step() error {
	return fmt.Errorf("cannot move")
}

func performTask() error {
	if err := step(); err != nil {
		return ErrTaskFailed.Wrap(err)
	}
	return nil
}

func Example() {
	if err := performTask(); err != nil {
		fmt.Print(err)
		return
	}
	// Output: task has failed
	//-> cannot move
}
