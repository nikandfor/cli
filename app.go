package app

import (
	"fmt"
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
		Flags: []Flag{
			HelpFlag,
		},
	}

	HelpFlag = &F{Name: "help", Aliases: []string{"h"}, After: DefaultHelpFlagAction}
)

func Run(args []string) error {
	return App.Run(args)
}

func RunAndExit(args []string) {
	err := Run(args)
	if err == nil {
		return
	}

	fmt.Fprintf(ErrWriter, "error: %v\n", err)
	os.Exit(1)
}

func DefaultHelpFlagAction(f Flag, c *Command) error {
	err := DefaultHelpAction(c)
	if err == nil {
		return ErrFlagExit
	}
	return err
}

func DefaultHelpAction(c *Command) error {
	if len(c.Commands) == 0 {
		fmt.Fprintf(Writer, "usage: %s [OPTION...] [ARGS...]\n", c.Name)
	} else {
		fmt.Fprintf(Writer, "usage: %s [OPTION...] [COMMAND...]\n", c.Name)
	}

	if len(c.Commands) != 0 {
		fmt.Fprintf(Writer, "\nCOMMANDS\n")
		for _, c := range c.Commands {
			fmt.Fprintf(Writer, "  %s %s\n", c.Name, c.Aliases)
		}
	}

	if len(c.Flags) != 0 {
		fmt.Fprintf(Writer, "\nFLAGS\n")
		for _, f := range c.Flags {
			b := f.Base()
			fmt.Fprintf(Writer, "  %s %s %T (default %v)\n", b.Name, b.Aliases, f, f.VAny())
		}
	}

	return nil
}
