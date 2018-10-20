package app

import (
	"io"
	"os"
	"path"
)

var (
	Reader    io.Reader = os.Stdin
	Writer    io.Writer = os.Stdout
	ErrWriter io.Writer = os.Stderr

	App = &Command{
		Name: path.Base(os.Args[0]),
	}
)
