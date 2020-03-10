// +build !go1.13

package cli

func NewNoSuchFlagError(f string) error {
	return ErrNoSuchFlag
}
