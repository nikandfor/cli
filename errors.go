package cli

import "fmt"

type (
	errWithContext struct {
		err error
		msg string
	}
)

func wrap(err error, msg string, args ...interface{}) error {
	if err == nil {
		panic("wrapping nil err")
	}

	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return &errWithContext{
		err: err,
		msg: msg,
	}
}

func (e *errWithContext) Error() string { return fmt.Sprintf("%v: %v", e.msg, e.err) }
func (e *errWithContext) Unwrap() error { return e.err }
