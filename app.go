package cli

import (
	"fmt"
	"os"
	"strings"
)

type (
	Args   []string
	Action func(c *Command) error

	Command struct {
		parent *Command
		Arg0   string // command name
		Args   Args

		Name        string
		Description string
		HelpText    string
		Action      Action
		Flags       []Flag
		Commands    []*Command
		Before      Action
		After       Action
		Complete    Action
	}

	Flag interface {
		Base() *F
		Parse(name, val string, more []string) (rest []string, err error)
	}
)

var (
	stdout = os.Stdout
	stderr = os.Stderr
)

var App = Command{
	Name: os.Args[0],
}

func RunAndExit(args []string) {
	err := App.run(args)
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%v\n", err)

	os.Exit(1)
}

func (c *Command) Bool(f string) bool {
	ff := c.Flag(f)
	if f, ok := ff.(*Bool); ok {
		return f.Value
	}
	return false
}

func (c *Command) String(f string) string {
	ff := c.Flag(f)
	if sf, ok := ff.(*String); ok {
		return sf.Value
	}
	return ""
}

func (c *Command) Int(f string) int {
	ff := c.Flag(f)
	if f, ok := ff.(*Int); ok {
		return f.Value
	}
	return 0
}

func (c *Command) run(args []string) (err error) {
	c.Arg0 = args[0]
	args = args[1:]
	noMoreFlags := false

	for len(args) > 0 {
		arg := args[0]

		//	tlog.Printf("arg %v %v", arg, args)

		switch {
		case arg == "--" && !noMoreFlags:
			noMoreFlags = true
			args = args[1:]
		case len(arg) != 0 && arg[0] == '-' && !noMoreFlags:
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
			f := c.Flag(k)
			if f == nil {
				return fmt.Errorf("no such flag: %v", arg)
			}
			args, err = f.Parse(k, v, args[1:])
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
			args = args[1:]
		}
	}

	if err = c.check(); err != nil {
		return err
	}

	if err = c.runBefore(); err != nil {
		return err
	}
	defer func() {
		if e := c.runAfter(); err == nil {
			err = e
		}
	}()

	if c.Action == nil {
		return defaultHelp(c)
	}

	return c.Action(c)
}

func (c *Command) runBefore() (err error) {
	if c.parent != nil {
		if err = c.parent.runBefore(); err != nil {
			return err
		}
	}
	if c.Before != nil {
		if err = c.Before(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) runAfter() (err error) {
	if c.After != nil {
		if err = c.After(c); err != nil {
			return err
		}
	}
	if c.parent != nil {
		if err = c.parent.runAfter(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) check() (err error) {
	for _, f := range c.Flags {
		if err = f.Base().check(); err != nil {
			return
		}
	}
	for _, sub := range c.Commands {
		if err = sub.check(); err != nil {
			return
		}
	}
	return nil
}

func (c *Command) sub(n string) *Command {
	for _, sub := range c.Commands {
		if sub.match(n) {
			sub.parent = c
			return sub
		}
	}
	return nil
}

func (c *Command) match(n string) bool {
	ns := strings.Split(c.Name, ",")
	for _, sn := range ns {
		if sn == n {
			return true
		}
	}
	return false
}

func (c *Command) Flag(n string) Flag {
	if c == nil {
		return nil
	}
	for _, f := range c.Flags {
		ns := strings.Split(f.Base().Name, ",")
		for _, fn := range ns {
			if fn == n {
				return f
			}
		}
	}
	return c.parent.Flag(n)
}

func (c *Command) Parent(n int) *Command {
	if c == nil || n == 0 {
		return c
	}
	return c.parent.Parent(n - 1)
}

func (a Args) Len() int { return len(a) }

func (a Args) String() string { return strings.Join(a, " ") }

func (a Args) First() string {
	if len(a) == 0 {
		return ""
	}
	return a[0]
}

func (a Args) Tail() []string {
	if len(a) == 0 {
		return nil
	}
	return a[1:]
}

func (a Args) Pop() (string, []string) {
	if len(a) == 0 {
		return "", nil
	}
	return a[0], a[1:]
}
