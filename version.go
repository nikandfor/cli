package cli

import (
	"fmt"

	"github.com/nikandfor/cli/flag"
)

func Version(ver, commit, date string) *Command {
	return &Command{
		Name:        "version",
		Description: "print version, commit hash and build date",
		Action: func(c *Command) (err error) {
			if c.Bool("short") {
				fmt.Printf("%v\n", ver)

				return nil
			}

			if c.Bool("commit") {
				fmt.Printf("%v\n", commit)

				return nil
			}

			if c.Bool("date") {
				fmt.Printf("%v\n", date)

				return nil
			}

			fmt.Printf("%v %v %v\n", ver, commit, date)

			return nil
		},
		Flags: []*Flag{
			flag.New("short", false, "prints only version tag"),
			flag.New("commit", false, "prints only commit hash"),
			flag.New("date", false, "prints only date"),
		},
	}
}
