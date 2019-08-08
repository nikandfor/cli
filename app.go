package cli

import (
	"fmt"
	"strings"
)

type (
	Args   []string
	Action func(c *Command) error

	Command struct {
		parent *Command

		Name     string
		Action   Action
		Flags    []FlagDev
		Commands []*Command
		Before   Action
		After    Action

		Args Args

		noMoreFlags bool
	}

	Flag interface{}

	FlagDev interface {
		Flag
		Names() string
		Parse(name, val string, more []string) (rest []string, err error)
	}
)

var App Command

func (c *Command) run(args []string) (err error) {
	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		switch {
		case arg == "--" && !c.noMoreFlags:
			c.noMoreFlags = true
		case len(arg) != 0 && arg[0] == '-' && !c.noMoreFlags:
			var k, v string
			if len(arg) != 1 && arg[1] == '-' {
				arg = arg[2:]
				if p := strings.IndexByte(arg, '='); p != -1 {
					k, v = arg[:p], arg[p:]
				} else {
					k = arg
				}
			} else {
				arg = arg[1:]
				if len(arg) != 0 {
					k, v = arg[:1], arg[1:]
				}
			}
			f := c.flag(k)
			if f == nil {
				return fmt.Errorf("no such flag: %v", arg)
			}
			args, err = f.Parse(k, v, args)
			if err != nil {
				return err
			}
		case c.Args == nil:
			sub := c.sub(arg)
			if sub != nil {
				return sub.run(args)
			}
			fallthrough
		default:
			c.Args = append(c.Args, arg)
		}
	}

	if c.Before != nil {
		err = c.Before(c)
		if err != nil {
			return err
		}
	}
	defer func() {
		if c.After != nil {
			e := c.After(c)
			if err == nil {
				err = e
			}
		}
	}()

	return c.Action(c)
}

func (c *Command) sub(n string) *Command {
	for _, sub := range c.Commands {
		ns := strings.Split(sub.Name, ",")
		for _, sn := range ns {
			if sn == n {
				sub.parent = c
				return sub
			}
		}
	}
	return nil
}

func (c *Command) flag(n string) FlagDev {
	for _, f := range c.Flags {
		ns := strings.Split(f.Names(), ",")
		for _, fn := range ns {
			if fn == n {
				return f
			}
		}
	}
	return nil
}
