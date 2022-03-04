package cli

import (
	"bytes"
	"fmt"

	"github.com/nikandfor/errors"
)

var HelpFlag *Flag

func init() {
	HelpFlag = &Flag{
		Name:        "help,h",
		Description: "print command help end exit",
		Action:      defaultHelp,
	}
}

func defaultHelp(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	b := new(bytes.Buffer)

	fmt.Fprintf(b, "usage: %s", c.Name)

	if c.Usage != "" {
		fmt.Fprintf(b, " %s", c.Usage)
	}

	if c.Description != "" {
		fmt.Fprintf(b, " - %s", c.Description)
	}

	fmt.Fprintf(b, "\n")

	if c.Help != "" {
		fmt.Fprintf(b, "\n%s\n", c.Help)
	}

	if len(c.Commands) != 0 {
		fmt.Fprintf(b, "\nSubcommands\n\n")

		for _, sub := range c.Commands {
			fmt.Fprintf(b, "    %-20s", sub.Name)

			if sub.Description != "" {
				fmt.Fprintf(b, " - %s", sub.Description)
			}

			fmt.Fprintf(b, "\n")
		}
	}

	if len(c.Flags) != 0 {
		fmt.Fprintf(b, "\nFlags\n\n")

		for c := c; c != nil; c = c.Parent {
			for _, f := range c.Flags {
				fmt.Fprintf(b, "    %-20s", f.Name)

				if f.Description != "" {
					fmt.Fprintf(b, " - %s", f.Description)
				}

				if f.Value != nil && f.Value != "" {
					fmt.Fprintf(b, " (default %v)", f.Value)
				}

				fmt.Fprintf(b, "\n")
			}
		}
	}

	_, err := b.WriteTo(c)
	if err != nil {
		return nil, errors.Wrap(err, "write")
	}

	return nil, ErrFlagExit
}
