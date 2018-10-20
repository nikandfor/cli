package app

import (
	"strings"

	"github.com/pkg/errors"
)

type (
	Args   []string
	Action func(c *Command) error

	Command struct {
		Name     string
		Aliases  []string
		Action   Action
		Commands []*Command
		Flags    []Flag
		Complete Action

		//	arg0 string
		args Args

		noMoreFlags bool

		parent *Command
	}
)

func (c *Command) Run(args []string) (err error) {
	//	arg0 := args[0]

	for args := args[1:]; len(args) > 0; {
		arg := args[0]
		//	log.Printf("run arg: %v", arg)
		switch {
		case len(arg) > 1 && arg[0] == '-' && arg[1] != '-' && !c.noMoreFlags:
			name, val := arg[1:2], arg[2:]
			f := c.flag(name)
			if f == nil {
				return errors.New("no such flag: " + name)
			}
			args, err = f.Parse(name, val, args)
			if err != nil {
				return errors.Wrap(err, "flag "+f.Base().Name)
			}
		case strings.HasPrefix(arg, "--") && arg != "--" && !c.noMoreFlags:
			name, val := arg[2:], ""
			if p := strings.Index(name, "="); p != -1 {
				val = name[p:]
				name = name[:p]
			}
			f := c.flag(name)
			if f == nil {
				return errors.New("no such flag: " + name)
			}
			args, err = f.Parse(name, val, args)
			if err != nil {
				return errors.Wrap(err, "flag "+f.Base().Name)
			}
		case arg == "--" && !c.noMoreFlags:
			c.noMoreFlags = true
			args = args[1:]
		default:
			if c.args == nil {
				sub := c.sub(arg)
				if sub != nil {
					sub.parent = c
					return sub.Run(args)
				}
				//	c.args = append(c.args, arg0)
			}
			c.args = append(c.args, arg)
			args = args[1:]
		}
	}

	return c.Action(c)
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
