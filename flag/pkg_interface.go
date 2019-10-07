package flag

import (
	"os"
	"time"

	"github.com/nikandfor/cli"
)

var (
	cmd = cli.Command{
		Name:   os.Args[0],
		Action: cli.NoAction,
		Flags: []*cli.Flag{
			cli.HelpFlag,
		},
	}
)

func Arg(i int) string {
	return cmd.Args[i]
}

func Args() []string {
	return cmd.Args
}

func NArg() int {
	return cmd.Args.Len()
}

func Lookup(n string) *cli.Flag {
	return cmd.Flag(n)
}

func Parse() {
	err := cli.RunCommand(&cmd, os.Args)
	if err != nil {
		panic(err)
	}
}

func Bool(name string, val bool, help string) *bool {
	fv := &cli.Bool{val}
	f := cli.NewFlag(name, fv, help)
	cmd.Flags = append(cmd.Flags, f)
	return &fv.Value
}

func Int(name string, val int, help string) *int {
	fv := &cli.Int{val}
	f := cli.NewFlag(name, fv, help)
	cmd.Flags = append(cmd.Flags, f)
	return &fv.Value
}

func String(name, val, help string) *string {
	fv := &cli.String{val}
	f := cli.NewFlag(name, fv, help)
	cmd.Flags = append(cmd.Flags, f)
	return &fv.Value
}

func Duration(name string, val time.Duration, help string) *time.Duration {
	fv := &cli.Duration{val}
	f := cli.NewFlag(name, fv, help)
	cmd.Flags = append(cmd.Flags, f)
	return &fv.Value
}
