package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type (
	Args   []string
	Action func(c *Command) error

	Command struct {
		Parent *Command
		Arg0   string // command name
		Args   Args

		Name          string
		Usage         string
		Description   string
		HelpText      string
		Action        Action
		Flags         []*Flag
		Commands      []*Command
		Before        Action
		After         Action
		Complete      Action
		FlagEnvPrefix string
	}
)

var ( // stdout/stderr
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

var ( // App
	App = Command{
		Name: os.Args[0],
	}
)

var ( // errors
	ErrAliasNotFound = errors.New("alias command not found")
	ErrNoSuchFlag    = errors.New("no such flag")
	ErrBadArguments  = errors.New("bad arguments")
)

func Chain(a ...Action) Action {
	return func(c *Command) error {
		for _, a := range a {
			err := a(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func SubcommandAlias(n string) Action {
	return func(c *Command) (err error) {
		sub := c.sub(n)
		if sub == nil {
			return ErrAliasNotFound
		}
		sub.Args = c.Args
		if a := sub.Before; a != nil {
			if err = a(sub); err != nil {
				return
			}
		}
		defer func() {
			if a := sub.After; a != nil {
				if err = a(sub); err != nil {
					return
				}
			}
		}()
		return sub.Action(sub)
	}
}

func Run(args []string) error {
	return App.run(args)
}

func RunAndExit(args []string) {
	err := App.run(args)
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%v\n", err)

	os.Exit(1)
}

func RunCommand(c *Command, args []string) error {
	return c.run(args)
}

func NoAction(c *Command) error { return nil }

func (c *Command) Bool(f string) bool {
	ff := c.Flag(f)
	if ff == nil {
		panic(fmt.Sprintf("no such flag: %v", f))
	}
	return ff.Value.(bool)
}

func (c *Command) String(f string) string {
	ff := c.Flag(f)
	if ff == nil {
		panic(fmt.Sprintf("no such flag: %v", f))
	}
	return ff.Value.(string)
}

func (c *Command) Int(f string) int {
	ff := c.Flag(f)
	if ff == nil {
		panic(fmt.Sprintf("no such flag: %v", f))
	}
	return ff.Value.(int)
}

func (c *Command) Duration(f string) time.Duration {
	ff := c.Flag(f)
	if ff == nil {
		panic(fmt.Sprintf("no such flag: %v", f))
	}
	return ff.Value.(time.Duration)
}

func (c *Command) StringSlice(f string) []string {
	ff := c.Flag(f)
	if ff == nil {
		panic(fmt.Sprintf("no such flag: %v", f))
	}
	return ff.Value.([]string)
}

func (c *Command) run(args []string) (err error) {
	defer func() {
		if err == ErrFlagExit {
			err = nil
		}
	}()

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
		case len(arg) >= 2 && arg[0] == '-' && arg != "--" && !noMoreFlags:
			args, err = c.parseFlag(arg, args)
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
		return defaultHelp(nil, c)
	}

	return c.Action(c)
}

func (c *Command) parseFlag(arg string, args []string) (rest []string, err error) {
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
		return nil, fmt.Errorf("%w: %v", ErrNoSuchFlag, arg)
	}
	if a := f.Before; a != nil {
		if err = a(f, c); err != nil {
			return
		}
	}
	switch fv := f.Value.(type) {
	case FlagValue:
		rest, err = fv.Parse(f, k, v, args[1:])
	case bool:
		rest, err = parseBool(f, k, v, args[1:])
	case int:
		rest, err = parseInt(f, k, v, args[1:])
	case string:
		rest, err = parseString(f, k, v, args[1:])
	case time.Duration:
		rest, err = parseDuration(f, k, v, args[1:])
	case []string:
		rest, err = parseStringSlice(f, k, v, args[1:])
	default:
		panic("unknown flag type")
	}
	if err != nil {
		return
	}
	f.IsSet = true
	if a := f.After; a != nil {
		if err = a(f, c); err != nil {
			return
		}
	}
	return
}

func (c *Command) runBefore() (err error) {
	if c.Parent != nil {
		if err = c.Parent.runBefore(); err != nil {
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
	if c.Parent != nil {
		if err = c.Parent.runAfter(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) check() (err error) {
	for _, f := range c.Flags {
		if err = f.check(); err != nil {
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
			sub.Parent = c
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

func (c *Command) Flag(n string) *Flag {
	if c == nil {
		return nil
	}
	for _, f := range c.Flags {
		ns := strings.Split(f.Name, ",")
		for _, fn := range ns {
			if fn == n {
				return f
			}
		}
	}
	return c.Parent.Flag(n)
}

func (a Args) Len() int { return len(a) }

func (a Args) String() string { return strings.Join(a, " ") }

func (a Args) First() string {
	if len(a) == 0 {
		return ""
	}
	return a[0]
}

func (a Args) Tail() Args {
	if len(a) == 0 {
		return nil
	}
	return a[1:]
}

func (a Args) Pop() (string, Args) {
	if len(a) == 0 {
		return "", nil
	}
	return a[0], a[1:]
}
