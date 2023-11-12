package cli

import (
	stderrors "errors"

	"nikand.dev/go/cli/flag"
)

type (
	Flag       = flag.Flag
	FlagAction = flag.Action
)

var (
	ErrExit       = stderrors.New("exit") // to stop command execution from flag or before handler
	ErrNoSuchFlag = stderrors.New("no such flag")
)

var DefaultFlags = []*Flag{
	FlagfileFlag,
	EnvfileFlag,
	HelpFlag,
}

func NewFlag(name string, val interface{}, help string, opts ...flag.Option) (f *Flag) {
	return flag.New(name, val, help, opts...)
}

func DefaultParseFlag(c *Command, arg string, args []string) (nextArgs []string, err error) {
	name := flagName(arg)

	f := c.Flag(name)
	if f == nil {
		return nil, ErrNoSuchFlag
	}

	f.CurrentCommand = c

	return f.Action(f, arg, args)
}

func flagName(arg string) string {
	st := 0
	for st < 2 && st < len(arg) && arg[st] == '-' {
		st++
	}

	end := st
	for end < len(arg) && arg[end] != '=' && arg[end] != ' ' {
		end++
	}

	return arg[st:end]
}
