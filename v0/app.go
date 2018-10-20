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

	App Command = Command{
		Name:   path.Base(os.Args[0]),
		Action: Help,
	}

	HelpFlag = &F{
		Name:    "help",
		Aliases: []string{"h"},
	}
)

var DefaultCommandAction = Help

func Run(args []string) error {
	args[0] = path.Base(args[0])
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

func Help(c *Command) error {
	base := path.Base(c.Args().First())
	if len(c.Commands) != 0 {
		fmt.Fprintf(Writer, "usage: %s [OPTIONS]... [SUBCOMMAND]\n", base)
	} else {
		fmt.Fprintf(Writer, "usage: %s [OPTIONS]... [ARGS]...\n", base)
	}

	fmt.Fprintf(Writer, "\nOPTIONS\n")
	for _, f := range c.Flags {
		b := f.Base()
		if b.Hidden {
			continue
		}
		fmt.Fprintf(Writer, "  %s %s %T\n", b.Name, b.Aliases, f)
	}

	fmt.Fprintf(Writer, "\nCOMMANDS\n")
	for _, c := range c.Commands {
		if c.Hidden {
			continue
		}
		fmt.Fprintf(Writer, "  %s\n", c.Name)
	}

	return nil
}

func (c *Command) String() string {
	var subs []string
	for _, c := range c.Commands {
		subs = append(subs, c.Name)
	}
	var flag []string
	for _, f := range c.Flags {
		flag = append(flag, f.Base().String())
	}
	return fmt.Sprintf("{%v %v act %v subs %v flags %v}", c.Name, c.Aliases, c.Action != nil, subs, flag)
}
