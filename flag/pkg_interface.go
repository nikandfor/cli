package flag

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nikandfor/cli"
)

var (
	flags []cli.Flag
	args  cli.Args
)

func Parse() {
	as := os.Args[1:]
	nomore := false
	var err error

	for len(as) > 0 {
		arg := as[0]
		as = as[1:]

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
			as, err = f.Parse(k, v, as)
			if err != nil {
				panic(err)
			}
		default:
			args = append(args, arg)
		}
	}
}

func getflag(n string) cli.Flag {
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

func Bool(name string, val bool, help string) *bool {
	f := cli.NewBool(name, val, help)
	flags = append(flags, f)
	return &f.Value
}

func Int(name string, val int, help string) *int {
	f := cli.NewInt(name, val, help)
	flags = append(flags, f)
	return &f.Value
}

func String(name, val, help string) *string {
	f := cli.NewString(name, val, help)
	flags = append(flags, f)
	return &f.Value
}

func Duration(name string, val time.Duration, help string) *time.Duration {
	f := cli.NewDuration(name, val, help)
	flags = append(flags, f)
	return &f.Value
}
