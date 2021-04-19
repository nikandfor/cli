package cli

import (
	"errors"
	"fmt"
)

var ( // errors
	ErrAliasNotFound = errors.New("alias command not found")
	ErrNoSuchFlag    = errors.New("no such flag")
	ErrNoSuchCommand = errors.New("no such command")
	ErrBadArguments  = errors.New("bad arguments")
)

func NewNoSuchFlagError(f string) error {
	return fmt.Errorf("%w: %v", ErrNoSuchFlag, f)
}

func NewNoSuchCommandError(c string) error {
	return fmt.Errorf("%w: %v", ErrNoSuchCommand, c)
}
