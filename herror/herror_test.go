package herror_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"dev-git.sentia.com/go/kit/herror"
	"dev-git.sentia.com/go/kit/test"
)

func TestHError(t *testing.T) {

	t.Run("new error", func(t *testing.T) {
		errDefault := "this is an error"
		for _, tc := range []struct {
			m        string
			input    interface{}
			expected string
		}{
			{
				m: "empty",
			},
			{
				m:        "string",
				input:    errDefault,
				expected: errDefault,
			},
			{
				m:        "error",
				input:    fmt.Errorf(errDefault),
				expected: errDefault,
			},
			{
				m:        "herror.Err",
				input:    herror.New(errDefault),
				expected: errDefault,
			},
			{
				m:        "invalid type",
				input:    123456789,
				expected: "",
			},
		} {
			t.Run(tc.m, func(t *testing.T) {
				test.Equals(t, tc.expected, herror.New(tc.input).Error())
			})
		}
	})

	t.Run("wrap", func(t *testing.T) {
		errmain := herror.New("MAIN ERROR")
		errfmt := fmt.Errorf("ERROR FORMATTED")
		errA := herror.New("ERR A")
		errB := herror.New("ERR B")
		errC := herror.New("ERR C")
		errD := herror.New("ERR D")
		errNested := errmain.Wrap(
			errA.Wrap(
				errB.Wrap(
					errC.Wrap(errD),
				),
			),
		)

		for _, tc := range []struct {
			m        string
			err      error
			expected []error
		}{
			{
				m:   "error",
				err: errfmt,
				expected: []error{
					errfmt,
				},
			},
			{
				m:   "deeper nested wrap",
				err: errNested,
				expected: []error{
					errA, errB, errC, errD,
				},
			},
		} {
			t.Run(tc.m, func(t *testing.T) {
				newerr := errmain.Wrap(tc.err)

				for _, e := range tc.expected {
					test.Equals(t, true, newerr.Is(e))
				}
			})
		}
	})

	t.Run("json marshalling", func(t *testing.T) {
		hError := herror.New("this is an error").
			Wrap(fmt.Errorf("this is another error")).
			CaptureStack()
		marshalled, err := json.Marshal(hError)
		test.OK(t, err)

		var unmarshalled *herror.Err
		test.OK(t, json.Unmarshal(marshalled, &unmarshalled))
		test.Equals(t, hError, unmarshalled)
	})
}

func ExampleErr_Wrap() {
	errA := herror.New("something went wrong")
	errB := fmt.Errorf("because of this error")
	newerr := herror.Wrap(errA, errB)

	fmt.Print(herror.Unwrap(newerr), "\n", newerr)
	// Output: because of this error
	// something went wrong
	// -> because of this error
}

func ExampleErr_Is() {
	errA := herror.New("something went wrong")
	errB := func() error {
		return errA
	}()

	fmt.Print(herror.Is(errA, errB))
	// Output: true
}

func ExampleErr_CaptureStack() {
	err := herror.New("something went wrong")
	err.CaptureStack()

	fmt.Print(err, "\n", err.Stack().Frames[2].Function)
	// Output: something went wrong
	// ExampleErr_CaptureStack
}

func ExampleErr_AddDetails() {
	err := herror.New("something went wrong")
	err.AddDetails(struct {
		number int
	}{123})

	fmt.Print(err, err.Details())
	// Output: something went wrong
	// (struct { number int }) {
	//  number: (int) 123
	// }
}
