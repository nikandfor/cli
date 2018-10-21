package app

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
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

	HelpFlag = &F{Name: "help", Aliases: []string{"h"},
		After:       DefaultHelpFlagAction,
		Description: "shows this message"}
	HelpHiddenFlag = F{Name: "hidden", Aliases: []string{"h"}, Hidden: true}.NewLevel(0)

	HelpCommand = &Command{
		Name:        "help",
		Action:      HelpAction,
		Description: "shows this message",
		Flags:       []Flag{HelpHiddenFlag},
	}
)

func NewActionChain(a ...Action) Action {
	return func(c *Command) error {
		for _, a := range a {
			if err := a(c); err != nil {
				return err
			}
		}
		return nil
	}
}

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

var DefaultHelpAction = func(c *Command) error {
	if Complete {
		return nil
	}
	if c.Help != nil {
		return c.Help(c)
	}

	l := HelpHiddenFlag.VInt()

	if len(c.Commands) == 0 {
		fmt.Fprintf(Writer, "usage: %s [OPTION...] [ARGS...]\n", c.Name)
	} else {
		fmt.Fprintf(Writer, "usage: %s [OPTION...] [COMMAND...]\n", c.Name)
	}

	if c.Description != "" {
		fmt.Fprintf(Writer, "  %s\n", c.Description)
	}

	if len(c.Commands) != 0 {
		fmt.Fprintf(Writer, "\nCOMMANDS\n")
		for _, c := range c.Commands {
			if c.Hidden && l == 0 {
				continue
			}
			if c.Name[0] == '_' && l < 2 {
				continue
			}
			fmt.Fprintf(Writer, "  %10s %-10s - %s\n", c.Name, fmt.Sprintf("%s", c.Aliases), c.Description)
		}
	}

	if len(c.Flags) != 0 {
		fmt.Fprintf(Writer, "\nFLAGS\n")
		for _, f := range c.Flags {
			b := f.Base()
			if b.Hidden && l == 0 {
				continue
			}
			if b.Name[0] == '_' && l < 2 {
				continue
			}
			fmt.Fprintf(Writer, "  %10s %s %s - %s (default %v)\n", b.Name, b.Aliases, f.Type(), f.Base().Description, f.VAny())
		}
	}

	return nil
}

func HelpAction(c *Command) error {
	for c != nil && c.parent != nil {
		// we need to go deeper
		c = c.parent
	}
	if c == nil {
		return errors.New("help for nil")
	}
	return DefaultHelpAction(c)
}

func AddHelpCommandAndFlag() {
	have := false
	for _, c := range App.Commands {
		if c == HelpCommand {
			have = true
			break
		}
	}
	if !have {
		App.Commands = append(App.Commands, HelpCommand)
	}

	have = false
	for _, f := range App.Flags {
		if f == HelpFlag {
			have = true
			break
		}
	}
	if !have {
		App.Flags = append(App.Flags, HelpFlag)
	}
}
