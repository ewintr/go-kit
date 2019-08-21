package herror

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/xerrors"
)

// Err represents an error
type Err struct {
	error   string
	wrapped *Err
	details string
	stack   *Stacktrace
}

type errJSON struct {
	E string      `json:"error"`
	W *Err        `json:"wrapped"`
	D string      `json:"details"`
	S *Stacktrace `json:"stack"`
}

// New returns a new instance for Err type with assigned error
func New(err interface{}) *Err {
	newerror := new(Err)

	switch e := err.(type) {
	case string:
		newerror.error = e

	case error:
		if castErr, ok := e.(*Err); ok {
			return castErr
		}
		newerror.error = e.Error()
	}

	return newerror
}

// Wrap set an error that is wrapped by Err
func Wrap(err, errwrapped error) error {
	newerr := New(err)
	return newerr.Wrap(errwrapped)
}

// Unwrap returns a wrapped error if present
func Unwrap(err error) error {
	return xerrors.Unwrap(err)
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return xerrors.Is(err, target)
}

// Wrap set an error that is wrapped by Err
func (e *Err) Wrap(err error) *Err {
	wrapped := New(err)

	if deeper := xerrors.Unwrap(err); deeper != nil {
		Wrap(wrapped, deeper)
	}

	newerr := &Err{
		error:   e.error,
		wrapped: e.wrapped,
		details: e.details,
		stack:   e.stack,
	}
	newerr.wrapped = wrapped
	return newerr
}

// Unwrap returns wrapped error
func (e *Err) Unwrap() error {
	if e.wrapped == nil {
		return nil
	}
	return e.wrapped
}

// Is reports whether an error matches.
func (e *Err) Is(err error) bool {
	if e.wrapped != nil {
		return e.error == err.Error() || e.wrapped.Is(err)
	}
	return e.error == err.Error()
}

// CaptureStack sets stack traces when the method is called
func (e *Err) CaptureStack() *Err {
	e.stack = NewStacktrace()
	return e
}

// Stack returns full stack traces
func (e *Err) Stack() *Stacktrace {
	return e.stack
}

// AddDetails records variable info to the error mostly for debugging purposes
func (e *Err) AddDetails(v ...interface{}) *Err {
	buff := new(bytes.Buffer)
	fmt.Fprintln(buff, e.details)
	spew.Fdump(buff, v...)
	e.details = buff.String()

	return e
}

// Details returns error's details
func (e *Err) Details() string {
	return e.details
}

// Errors return a composed message of the assigned error e wrapped error
func (e *Err) Error() string {
	if e.wrapped == nil {
		return e.error
	}

	return fmt.Sprintf("%s\n-> %s", e.error, e.wrapped.Error())
}

// UnmarshalJSON
func (e *Err) UnmarshalJSON(b []byte) error {
	var errJSON errJSON
	if err := json.Unmarshal(b, &errJSON); err != nil {
		return err
	}

	*e = Err{
		error:   errJSON.E,
		wrapped: errJSON.W,
		details: errJSON.D,
		stack:   errJSON.S,
	}

	return nil
}

// MarshalJSON
func (e *Err) MarshalJSON() ([]byte, error) {
	return json.Marshal(errJSON{
		E: e.error,
		W: e.wrapped,
		D: e.details,
		S: e.stack,
	})
}
