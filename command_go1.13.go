// +build go1.13

package cli

import "fmt"

func NewNoSuchFlagError(f string) error {
	return fmt.Errorf("%w: %v", ErrNoSuchFlag, f)
}

func NewNoSuchCommandError(c string) error {
	return fmt.Errorf("%w: %v", ErrNoSuchCommand, c)
}
