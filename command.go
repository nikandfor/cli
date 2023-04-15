package cli

import (
	stderrors "errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/nikandfor/errors"
)

type (
	Command struct {
		// Set by the lib

		Parent *Command

		OSArgs []string
		OSEnv  []string

		Arg0 string   // command name
		Args Args     // must be initialized to cli.Args{} if arguments expected
		Env  []string // env vars not used for local flags

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

		Stdout io.Writer
		Stderr io.Writer
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
	if false {
		c.Env = env

		if err := beforeComplete(c); err != nil {
			panic(err)
		}
	}

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
		if match(sub.Name, n) {
			sub.Parent = c

			return sub
		}
	}

	return nil
}

func (c *Command) Flag(n string) *Flag {
	for q := c; q != nil; q = q.Parent {
		for _, f := range q.Flags {
			if !match(f.Name, n) {
				continue
			}

			if f.Local && q != c {
				return nil
			}

			return f
		}
	}

	return nil
}

func (c *Command) run(args, env []string) (err error) {
	defer func() {
		if errors.Is(err, ErrExit) {
			err = nil
		}
	}()

	c.OSArgs = args
	c.OSEnv = env

	c.Arg0 = args[0]
	args = args[1:]

	err = c.setup()
	if err != nil {
		return errors.WrapNoCaller(err, "setup")
	}

	c.Env, err = c.parseEnv(env)
	if err != nil {
		return errors.WrapNoCaller(err, "parse env")
	}

	last, comp := c.completeIndex()
	if comp && c.Parent == nil {
		args = args[:last-1]
	}

	for len(args) != 0 {
		arg := args[0]

		if arg != "" && arg[0] == '-' && arg != "-" && arg != "--" {
			args, err = c.parseFlag(arg, args[1:])
			if err != nil {
				return errors.WrapNoCaller(err, "parse flag: %v", arg)
			}

			continue
		}

		if sub := c.Command(arg); sub != nil {
			err = sub.run(args, c.Env)

			return errors.WrapNoCaller(err, MainName(arg))
		}

		if c.Args == nil {
			return fmt.Errorf("%w, got %v", ErrNoArgsExpected, arg)
		}

		if arg == "--" {
			c.Args = append(c.Args, args[1:]...)
			args = nil
		} else {
			c.Args = append(c.Args, arg)
			args = args[1:]
		}
	}

	if err = c.check(); err != nil {
		return errors.WrapNoCaller(err, "check flags")
	}

	err = c.runBefore()
	if err != nil {
		return errors.WrapNoCaller(err, "before")
	}

	defer func() {
		e := c.runAfter()
		if err == nil {
			err = errors.WrapNoCaller(e, "after")
		}
	}()

	if comp {
		return c.complete()
	}

	if c.Action == nil {
		_, err = defaultHelp(nil, "", nil)
		return errors.WrapNoCaller(err, "help")
	}

	return c.Action(c)
}

func (c *Command) setup() error {
	if c.Stdout == nil {
		if c.Parent != nil {
			c.Stdout = c.Parent.Stdout
		} else {
			c.Stdout = os.Stdout
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
	if c.ParseFlag != nil {
		return c.ParseFlag(c, arg, more)
	}

	return DefaultParseFlag(c, arg, more)
}

func (c *Command) parseEnv(env []string) (rest []string, err error) {
	for c := c; c != nil; c = c.Parent {
		if c.ParseEnv != nil {
			return c.ParseEnv(c, env)
		}
	}

	return DefaultParseEnv(c, env)
}

func (c *Command) completeIndex() (int, bool) {
	x, ok := c.LookupEnv("CLI_COMP_WORDS_INDEX")
	if !ok {
		return 0, false
	}

	i, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		return 0, false
	}

	return int(i), true
}

func (c *Command) complete() error {
	if c.Complete != nil {
		return c.Complete(c)
	}

	return DefaultComplete(c)
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
			return errors.WrapNoCaller(err, "%v", MainName(c.Parent.Name))
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
	if c.Parent != nil {
		defer func() {
			e := c.Parent.runAfter()
			if err == nil {
				err = errors.WrapNoCaller(e, "%v", MainName(c.Parent.Name))
			}
		}()
	}

	if c.After != nil {
		if err = c.After(c); err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) check() (err error) {
	if c.Parent != nil {
		err = c.Parent.check()
		if err != nil {
			return err
		}
	}

	for _, f := range c.Flags {
		if err = f.Check(); err != nil {
			return errors.WrapNoCaller(err, f.MainName())
		}
	}

	return nil
}

func match(name, sub string) bool {
	ns := strings.Split(name, ",")
	for _, sn := range ns {
		if sn == sub {
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

func (a Args) Get(i int) string {
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

func FullName(c *Command) (r []string) {
	return fullName(c, r)
}

func fullName(c *Command, r []string) []string {
	if c == nil {
		return r
	}

	r = fullName(c.Parent, r)

	return append(r, MainName(c.Name))
}
