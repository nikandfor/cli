package cli

import (
	stderrors "errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nikandfor/errors"
)

type (
	Command struct {
		// Set by runner

		Parent *Command

		Arg0 string // command name
		Args Args

		// User options

		Name        string
		Group       string
		Usage       string
		Description string
		Help        string

		Before Action
		After  Action
		Action Action

		Complete Action

		Flags    []*Flag
		Commands []*Command

		Hidden bool

		EnvPrefix string

		ParseEnv  func(c *Command, env []string) ([]string, error)
		ParseFlag func(c *Command, arg string, args []string) ([]string, error)

		io.Writer // Strout
		Stderr    io.Writer
	}

	Action func(c *Command) error

	Args []string
)

var (
	ErrNoSuchCommand  = stderrors.New("no such command")
	ErrNoArgsExpected = stderrors.New("no args expected")
)

func Run(c *Command, args, env []string) (err error) {
	return c.run(args, env)
}

func RunAndExit(c *Command, args, env []string) {
	err := c.run(args, env)
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "error: %v\n", err)

	os.Exit(1)
}

func (c *Command) MainName() string {
	return MainName(c.Name)
}

func (c *Command) Command(n string) *Command {
	for _, sub := range c.Commands {
		if sub.match(n) {
			sub.Parent = c

			return sub
		}
	}

	return nil
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

func (c *Command) run(args, env []string) (err error) {
	defer func() {
		if errors.Is(err, ErrFlagExit) {
			err = nil
		}
	}()

	c.Arg0 = args[0]
	args = args[1:]

	err = c.setup()
	if err != nil {
		return errors.WrapNoLoc(err, "setup")
	}

	env, err = c.parseEnv(env)
	if err != nil {
		return errors.WrapNoLoc(err, "parse env")
	}

	for len(args) != 0 {
		arg := args[0]

		switch {
		case arg != "" && arg[0] == '-' && arg != "-":
			args, err = c.parseFlag(arg, args[1:])
			if err != nil {
				return errors.WrapNoLoc(err, "parse flag: %v", arg)
			}
		case len(c.Commands) != 0:
			sub := c.Command(arg)
			if sub == nil {
				return errors.Wrap(ErrNoSuchCommand, "%v", arg)
			}

			err = sub.run(args, env)

			return errors.WrapNoLoc(err, MainName(sub.Name))
		case c.Args == nil:
			return fmt.Errorf("%w, got %v", ErrNoArgsExpected, arg)
		default:
			c.Args = append(c.Args, arg)
			args = args[1:]
		}
	}

	if err = c.check(); err != nil {
		return errors.WrapNoLoc(err, "check")
	}

	err = c.runBefore()
	if err != nil {
		return errors.WrapNoLoc(err, "before")
	}

	defer func() {
		e := c.runAfter()
		if err == nil {
			err = errors.WrapNoLoc(e, "after")
		}
	}()

	if c.Action == nil {
		_, err = defaultHelp(c, nil, "", nil)
		return errors.WrapNoLoc(err, "help")
	}

	return c.Action(c)
}

func (c *Command) setup() error {
	if c.Writer == nil {
		if c.Parent != nil {
			c.Writer = c.Parent.Writer
		} else {
			c.Writer = os.Stdout
		}
	}

	if c.Stderr == nil {
		if c.Parent != nil {
			c.Stderr = c.Parent.Stderr
		} else {
			c.Stderr = os.Stderr
		}
	}

	return nil
}

func (c *Command) parseFlag(arg string, more []string) (rest []string, err error) {
	if c.ParseFlag == nil {
		return DefaultParseFlag(c, arg, more)
	}

	return c.ParseFlag(c, arg, more)
}

func (c *Command) parseEnv(env []string) (rest []string, err error) {
	if c.ParseEnv == nil {
		return ParseEnv(c, env)
	}

	return c.ParseEnv(c, env)
}

func GetEnvPrefix(c *Command) string {
	if c == nil {
		return ""
	}

	if c.EnvPrefix != "" {
		return c.EnvPrefix
	}

	return GetEnvPrefix(c.Parent)
}

func (c *Command) runBefore() (err error) {
	if c.Parent != nil {
		if err = c.Parent.runBefore(); err != nil {
			return errors.WrapNoLoc(err, "%v", MainName(c.Parent.Name))
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
			return errors.WrapNoLoc(err, "%v", MainName(c.Parent.Name))
		}
	}
	return nil
}

func (c *Command) check() (err error) {
	for _, f := range c.Flags {
		if err = f.check(); err != nil {
			return errors.WrapNoLoc(err, "")
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

func (a Args) Len() int { return len(a) }

func (a Args) String() string { return strings.Join(a, " ") }

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

func (a Args) Pop() (string, Args) {
	if len(a) == 0 {
		return "", nil
	}
	return a[0], a[1:]
}

func (a Args) SafeGet(i int) string {
	if i < 0 || i >= len(a) {
		return ""
	}

	return a[i]
}

func MainName(n string) string {
	p := strings.IndexByte(n, ',')
	if p == -1 {
		return n
	}

	return n[:p]
}
