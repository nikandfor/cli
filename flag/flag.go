package flag

import (
	"flag"
	"os"
	"sort"
	"time"

	"github.com/nikandfor/cli"
)

var CommandLine = &cli.Command{
	Name: os.Args[0],
	Args: cli.Args{},
}

var parsed bool

func Parse() {
	parsed = true

	err := cli.Run(CommandLine, os.Args, os.Environ())
	if err != nil {
		panic(err)
	}
}

func Parsed() bool {
	return parsed
}

func NArg() int { return CommandLine.Args.Len() }

func Arg(i int) string {
	return CommandLine.Args.SafeGet(i)
}

func Args() []string {
	return CommandLine.Args
}

func NFlag() int { return len(CommandLine.Flags) }

func Lookup(name string) *cli.Flag {
	return CommandLine.Flag(name)
}

func Bool(name string, value bool, usage string) *bool {
	p := new(bool)

	BoolVar(p, name, value, usage)

	return p
}

func BoolVar(p *bool, name string, value bool, usage string) {
	f := &cli.Flag{
		Name:  name,
		Value: value,
		Action: func(c *cli.Command, f *cli.Flag, arg string, args []string) (_ []string, err error) {
			args, err = cli.ParseFlagBool(c, f, arg, args)
			if err != nil {
				return
			}

			*p = f.Value.(bool)

			return args, nil
		},
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Duration(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)

	DurationVar(p, name, value, usage)

	return p
}

func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	f := &cli.Flag{
		Name:  name,
		Value: value,
		Action: func(c *cli.Command, f *cli.Flag, arg string, args []string) (_ []string, err error) {
			args, err = cli.ParseFlagDuration(c, f, arg, args)
			if err != nil {
				return
			}

			*p = f.Value.(time.Duration)

			return args, nil
		},
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Float64(name string, value float64, usage string) *float64 {
	p := new(float64)

	Float64Var(p, name, value, usage)

	return p
}

func Float64Var(p *float64, name string, value float64, usage string) {
	f := &cli.Flag{
		Name:  name,
		Value: value,
		Action: func(c *cli.Command, f *cli.Flag, arg string, args []string) (_ []string, err error) {
			args, err = cli.ParseFlagFloat64(c, f, arg, args)
			if err != nil {
				return
			}

			*p = f.Value.(float64)

			return args, nil
		},
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Int(name string, value int, usage string) *int {
	p := new(int)

	IntVar(p, name, value, usage)

	return p
}

func IntVar(p *int, name string, value int, usage string) {
	f := &cli.Flag{
		Name:  name,
		Value: value,
		Action: func(c *cli.Command, f *cli.Flag, arg string, args []string) (_ []string, err error) {
			args, err = cli.ParseFlagInt(c, f, arg, args)
			if err != nil {
				return
			}

			*p = f.Value.(int)

			return args, nil
		},
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Int64(name string, value int64, usage string) *int64 {
	p := new(int64)

	Int64Var(p, name, value, usage)

	return p
}

func Int64Var(p *int64, name string, value int64, usage string) {
	f := &cli.Flag{
		Name:  name,
		Value: value,
		Action: func(c *cli.Command, f *cli.Flag, arg string, args []string) (_ []string, err error) {
			args, err = cli.ParseFlagInt64(c, f, arg, args)
			if err != nil {
				return
			}

			*p = f.Value.(int64)

			return args, nil
		},
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func String(name string, value string, usage string) *string {
	p := new(string)

	StringVar(p, name, value, usage)

	return p
}

func StringVar(p *string, name string, value string, usage string) {
	f := &cli.Flag{
		Name:  name,
		Value: value,
		Action: func(c *cli.Command, f *cli.Flag, arg string, args []string) (_ []string, err error) {
			args, err = cli.ParseFlagString(c, f, arg, args)
			if err != nil {
				return
			}

			*p = f.Value.(string)

			return args, nil
		},
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Uint(name string, value uint, usage string) *uint {
	p := new(uint)

	UintVar(p, name, value, usage)

	return p
}

func UintVar(p *uint, name string, value uint, usage string) {
	f := &cli.Flag{
		Name:  name,
		Value: value,
		Action: func(c *cli.Command, f *cli.Flag, arg string, args []string) (_ []string, err error) {
			args, err = cli.ParseFlagUint(c, f, arg, args)
			if err != nil {
				return
			}

			*p = f.Value.(uint)

			return args, nil
		},
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Uint64(name string, value uint64, usage string) *uint64 {
	p := new(uint64)

	Uint64Var(p, name, value, usage)

	return p
}

func Uint64Var(p *uint64, name string, value uint64, usage string) {
	f := &cli.Flag{
		Name:  name,
		Value: value,
		Action: func(c *cli.Command, f *cli.Flag, arg string, args []string) (_ []string, err error) {
			args, err = cli.ParseFlagUint64(c, f, arg, args)
			if err != nil {
				return
			}

			*p = f.Value.(uint64)

			return args, nil
		},
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Var(v flag.Value, name string, usage string) {
	f := &cli.Flag{
		Name:   name,
		Value:  v,
		Action: cli.ParseFlagValue(v, true, false),
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Func(name string, usage string, fn func(string) error) {
	f := &cli.Flag{
		Name: name,
		Action: cli.ParseFlagFunc(func(val string) (x interface{}, err error) {
			err = fn(val)
			return
		}, true, false),
	}

	CommandLine.Flags = append(CommandLine.Flags, f)
}

func Visit(fn func(f *cli.Flag)) {
	l := make([]*cli.Flag, 0, len(CommandLine.Flags))

	for _, f := range CommandLine.Flags {
		if f.IsSet {
			l = append(l, f)
		}
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].MainName() < l[j].MainName()
	})

	for _, f := range l {
		fn(f)
	}
}

func VisitAll(fn func(f *cli.Flag)) {
	l := make([]*cli.Flag, 0, len(CommandLine.Flags))

	for _, f := range CommandLine.Flags {
		l = append(l, f)
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].MainName() < l[j].MainName()
	})

	for _, f := range l {
		fn(f)
	}
}
