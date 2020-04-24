// +build !go1.13

package cli

func NewNoSuchFlagError(f string) error {
	return ErrNoSuchFlag
}

func NewNoSuchCommandError(c string) error {
	return ErrNoSuchCommand
}
