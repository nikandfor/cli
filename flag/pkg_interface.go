package flag

import (
	"os"
	"path/filepath"
	"time"

	"github.com/nikandfor/cli"
)

var (
	cmd = cli.Command{
		Name:   filepath.Base(os.Args[0]),
		Action: cli.NoAction,
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
	if cmd.Flag("help") == nil {
		cmd.Flags = append(cmd.Flags, cli.HelpFlag)
	}

	cli.RunAndExit(&cmd, os.Args, os.Environ())
}

func Bool(name string, val bool, help string) *bool {
	f := cli.NewFlag(name, &val, help)
	cmd.Flags = append(cmd.Flags, f)
	return &val
}

func Int(name string, val int, help string) *int {
	f := cli.NewFlag(name, &val, help)
	cmd.Flags = append(cmd.Flags, f)
	return &val
}

func String(name, val, help string) *string {
	f := cli.NewFlag(name, &val, help)
	cmd.Flags = append(cmd.Flags, f)
	return &val
}

func Duration(name string, val time.Duration, help string) *time.Duration {
	f := cli.NewFlag(name, &val, help)
	cmd.Flags = append(cmd.Flags, f)
	return &val
}

func StringSlice(name string, val []string, help string) *[]string {
	f := cli.NewFlag(name, &val, help)
	cmd.Flags = append(cmd.Flags, f)
	return &val
}

func Usage(name, usage string) {
	if name != "" {
		cmd.Name = name
	}
	cmd.Usage = usage

	if cmd.Flag("help") == nil {
		cmd.Flags = append(cmd.Flags, cli.HelpFlag)
	}

	cli.RunAndExit(&cmd, append(os.Args, "-h"), os.Environ())
}
