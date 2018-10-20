package app

import (
	"log"
	"strings"
)

type (
	Args   []string
	Action func(*Command) error

	Command struct {
		Name       string
		Aliases    []string
		Action     Action
		Flags      []Flag
		Commands   []*Command
		Completion Action
		Hidden     bool

		run         Action
		args        Args
		noMoreFlags bool
		parent      *Command
	}
)

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
		return a
	}
	return a[1:]
}

func (c *Command) Args() Args {
	return c.args
}

var num int

func (c *Command) Run(args []string) error {
	log.Printf("cmd: %v", c.String())
	log.Printf("run: %v", args)
	arg0 := args[0]
	args = args[1:]
	var err error
	var f Flag

	comp, nl := ifCompletion()

	for len(args) > nl {
		arg := args[0]

		//	log.Printf("arg: %q %v %v", arg, c.noMoreFlags, len(arg))
		if !c.noMoreFlags && len(arg) > 1 && arg[0] == '-' && !(len(arg) == 2 && arg[1] == '-') {
			var val string
			if len(arg) > 2 && arg[1] == '-' {
				arg = arg[2:]
				if p := strings.Index(arg, "="); p != -1 {
					val = arg[p:]
					arg = arg[:p]
				}
			} else {
				val = arg[2:]
				arg = arg[1:2]
			}

			f = c.Flag(arg)
			if f == nil {
			}

			args, err = f.parse(arg, val, args)
			if err != nil {
				return err
			}

			comp, nl = ifCompletion()

			continue
		}

		if arg == "--" {
			c.noMoreFlags = true
			args = args[1:]
			f = nil
			continue
		}

		if c.args == nil {
			sub := c.getSubcommand(arg)
			if sub != nil {
				sub.parent = c
				return sub.Run(args)
			}
			c.args = append(c.args, arg0)
		}
		c.args = append(c.args, arg)
		args = args[1:]
		f = nil
	}

	log.Printf("cmd: %v", c.String())
	log.Printf("act: %v %v (comp %v %v) %v", c.Name, c.args, comp, nl, f.Base().String())

	if comp {
		var last string
		if nl != 0 && len(args) != 0 {
			last = args[0]
		}

		if false && f != nil {
			return f.complete(last)
		}

		c.args = append(c.args, last)

		if c.Completion != nil {
			return c.Completion(c)
		}
		return DefaultCommandCompletion(c)
	}

	if c.Action != nil {
		return c.Action(c)
	}
	return DefaultCommandAction(c)
}

func (c *Command) Flag(n string) Flag {
	f := c.flag(n)
	return f
}

func (c *Command) StringFlag(n string) *StringFlag {
	f := c.flag(n)
	sf, _ := f.(*StringFlag)
	return sf
}

func (c *Command) BoolFlag(n string) *BoolFlag {
	f := c.flag(n)
	sf, _ := f.(*BoolFlag)
	return sf
}

func (c *Command) IntFlag(n string) *IntFlag {
	f := c.flag(n)
	sf, _ := f.(*IntFlag)
	return sf
}

func (c *Command) flag(n string) Flag {
	//	log.Printf("at %-10s: Flag %v", path.Base(c.Name), n)
	num++
	if num > 10 {
		panic(num)
	}
	for _, f := range c.Flags {
		if n == f.Base().Name {
			return f
		}
		for _, a := range f.Base().Aliases {
			if n == a {
				return f
			}
		}
	}
	if c.parent != nil {
		return c.parent.flag(n)
	}
	return nil
}

func (c *Command) FlagsSet() []string {
	var r []string
	if c.parent != nil {
		r = c.parent.FlagsSet()
	}
	for _, f := range c.Flags {
		if !f.IsSet() {
			continue
		}
		r = append(r, f.Base().Name)
	}
	return r
}

func (c *Command) VisibleCommands() []*Command {
	var r []*Command
	for _, c := range c.Commands {
		if c.Hidden {
			continue
		}
		r = append(r, c)
	}
	return r
}

func (c *Command) VisibleFlags() []Flag {
	var r []Flag
	for _, f := range c.Flags {
		if f.Base().Hidden {
			continue
		}
		r = append(r, f)
	}
	return r
}

func (c *Command) getSubcommand(arg string) *Command {
	for _, sub := range c.Commands {
		if sub.Name == arg {
			return sub
		}
		for _, a := range sub.Aliases {
			if a == arg {
				return sub
			}
		}
	}
	return nil
}
