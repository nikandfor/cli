package app

import (
	"strings"

	"github.com/pkg/errors"
)

var ErrFlagExit = errors.New("flag exit")
var ErrNoSuchCommand = errors.New("no such command found")

type (
	Args   []string
	Action func(c *Command) error

	Command struct {
		Name       string
		Aliases    []string
		Action     Action
		Commands   []*Command
		Flags      []Flag
		Before     Action
		After      Action
		Completion Action
		Hidden     bool

		noMoreFlags bool

		//	arg0 string
		args Args

		parent *Command
	}
)

func (c *Command) Run(args []string) (err error) {
	//	arg0 := args[0]
	defer func() {
		if err == ErrFlagExit {
			err = nil
		}
	}()

	for args := args[1:]; len(args) > NLastComplete; {
		arg := args[0]
		//	log.Printf("run arg: %v %v\n", arg, NLastComplete)
		switch {
		case len(arg) > 1 && arg[0] == '-' && arg != "--" && !c.noMoreFlags:
			args, err = c.parseFlag(args)
			if err != nil {
				return err
			}
		case arg == "--" && !c.noMoreFlags:
			c.noMoreFlags = true
			args = args[1:]
		case c.args == nil:
			sub := c.sub(arg)
			if sub != nil {
				sub.parent = c
				return sub.Run(args)
			}
			//	c.args = append(c.args, arg0)
			fallthrough
		default:
			c.args = append(c.args, arg)
			args = args[1:]
		}
	}

	if ok, last := CompleteLast(args); ok {
		c.args = append(c.args, last)
		if c.Completion != nil {
			return c.Completion(c)
		}
		return DefaultCommandCompletion(c)
	}

	if c.Before != nil {
		if err := c.Before(c); err != nil {
			return err
		}
	}
	defer func() {
		if c.After != nil {
			if e := c.After(c); err == nil {
				err = e
			}
		}
	}()

	if c.Action == nil {
		err := DefaultHelpAction(c)
		if err != nil {
			return err
		}
		return ErrNoSuchCommand
	}

	return c.Action(c)
}

func (c *Command) parseFlag(args []string) ([]string, error) {
	var err error
	arg, args := Pop(args)

	var name, val string
	if arg[1] != '-' {
		name, val = arg[1:2], arg[2:]
	} else {
		name, val = arg[2:], ""
		if p := strings.Index(name, "="); p != -1 {
			val = name[p:]
			name = name[:p]
		}
	}

	f := c.flag(name)
	if f == nil {
		return nil, errors.New("no such flag: " + name)
	}

	if a := f.Base().Before; a != nil {
		if err = a(f, c); err != nil {
			return nil, err
		}
	}

	var rep bool
	for {
		more, err := f.Parse(name, val, rep)
		if err != nil {
			return nil, err
		}

		if !more {
			break
		}

		if len(args) > NLastComplete {
			val, args = Pop(args)
			rep = true
			continue
		}

		if ok, last := CompleteLast(args); ok {
			if comp := f.Base().Completion; comp != nil {
				err = comp(f, c, last)
			} else {
				err = DefaultFlagCompletion(f, c, last)
			}
			if err != nil {
				return nil, err
			}
			return nil, ErrFlagExit
		}

		return nil, errors.New("argument expected")
	}

	if a := f.Base().After; a != nil {
		if err = a(f, c); err != nil {
			return nil, err
		}
	}

	return args, nil
}

func (c *Command) Flag(n string) Flag {
	for c := c; c != nil; c = c.parent {
		for _, f := range c.Flags {
			b := f.Base()
			if b.Name == n {
				return f
			}
			for _, a := range b.Aliases {
				if a == n {
					return f
				}
			}
		}
	}
	return nil
}

func (c *Command) Args() Args { return c.args }

func (a Args) First() string {
	if len(a) == 0 {
		return ""
	}
	return a[0]
}

func (a Args) Last() string {
	if len(a) == 0 {
		return ""
	}
	return a[len(a)-1]
}

func (a Args) Tail() Args {
	if len(a) == 0 {
		return nil
	}
	return a[1:]
}

// --

func (c *Command) flag(n string) FlagDev {
	f := c.Flag(n)
	switch f := f.(type) {
	case FlagDev:
		return f
	default:
		return nil
	}
}

func (c *Command) sub(n string) *Command {
	for _, sub := range c.Commands {
		if sub.Name == n {
			return sub
		}
		for _, a := range sub.Aliases {
			if a == n {
				return sub
			}
		}
	}
	return nil
}

func Pop(args []string) (string, []string) {
	return args[0], args[1:]
}
