package flag

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type flag interface {
	Names() string
	Parse(n, v string, more []string) ([]string, error)
}

var (
	flags []flag
	args  []string
)

func Parse() {
	args := os.Args
	nomore := false
	var err error

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		switch {
		case arg == "--" && !nomore:
			nomore = true
		case len(arg) != 0 && arg[0] == '-' && !nomore:
			if len(arg) > 1 && arg[1] == '-' {
				arg = arg[2:]
			} else {
				arg = arg[1:]
			}
			var k, v string
			if p := strings.IndexByte(arg, '='); p != -1 {
				k, v = arg[:p], arg[p:]
			} else {
				k = arg
			}
			f := getflag(k)
			if f == nil {
				panic(fmt.Errorf("no such flag: %v", arg))
			}
			args, err = f.Parse(k, v, args)
			if err != nil {
				panic(err)
			}
		default:
			args = append(args, arg)
		}
	}
}

func getflag(n string) flag {
	for _, f := range flags {
		ns := strings.Split(f.Names(), ",")
		for _, fn := range ns {
			if fn == n {
				return f
			}
		}
	}
	return nil
}

func Int(name string, val int, help string) *int {
	f := &IntFlag{Name: name, Value: val, Help: help}
	flags = append(flags, f)
	return &f.Value
}

func String(name, val, help string) *string {
	f := &StringFlag{Name: name, Value: val, Help: help}
	flags = append(flags, f)
	return &f.Value
}

func Duration(name string, val time.Duration, help string) *time.Duration {
	f := &DurationFlag{Name: name, Value: val, Help: help}
	flags = append(flags, f)
	return &f.Value
}

func Bool(name string, val bool, help string) *bool {
	f := &BoolFlag{Name: name, Value: val, Help: help}
	flags = append(flags, f)
	return &f.Value
}
