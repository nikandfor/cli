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
		Usage:       "=[hidden]",
		Description: "print command help end exit",
		Action:      defaultHelp,
	}
}

func defaultHelp(c *Command, f *Flag, arg string, args []string) (rest []string, err error) {
	_, v, rest, err := ParseFlagArg(arg, args, false, true)
	if err != nil {
		return
	}

	hidden := v == "hidden"

	b := new(bytes.Buffer)

	//	full := FullName(c.Parent)
	//	if full != nil {
	//		fmt.Fprintf(b, "%s ", strings.Join(full, " "))
	//	}

	fmt.Fprintf(b, "%s", c.Name)

	if c.Usage != "" {
		fmt.Fprintf(b, " %s", c.Usage)
	} else if c.Args != nil {
		fmt.Fprintf(b, " [flags_and_args]")
	} else {
		fmt.Fprintf(b, " [flags]")
	}

	if c.Description != "" {
		fmt.Fprintf(b, " - %s", c.Description)
	}

	fmt.Fprintf(b, "\n")

	if c.Help != "" {
		fmt.Fprintf(b, "\n%s\n", c.Help)
	}

	if len(c.Commands) != 0 {
		cnt := 0
		for _, sub := range c.Commands {
			if sub.Hidden && !hidden {
				continue
			}

			cnt++
		}

		if cnt != 0 {
			fmt.Fprintf(b, "\nSubcommands\n")
		}

		for _, sub := range c.Commands {
			if sub.Hidden && !hidden {
				continue
			}

			fmt.Fprintf(b, "    %-20s", sub.Name)

			if sub.Description != "" {
				fmt.Fprintf(b, " - %s", sub.Description)
			}

			fmt.Fprintf(b, "\n")
		}
	}

	for cc := c; cc != nil; cc = cc.Parent {
		cnt := 0
		for _, f := range cc.Flags {
			if f.Hidden && !hidden || cc != c && f.Local {
				continue
			}

			cnt++
		}

		if cnt == 0 {
			continue
		}

		if cc == c {
			fmt.Fprintf(b, "\nFlags\n")
		} else {
			fmt.Fprintf(b, "\nFlags of parent command %v\n", cc.MainName())
		}

		for _, f := range cc.Flags {
			if f.Hidden && !hidden || cc != c && f.Local {
				continue
			}

			name := f.Name

			if f.Usage != "" {
				name += f.Usage
			}

			fmt.Fprintf(b, "    %-20s", name)

			if f.Description != "" {
				fmt.Fprintf(b, " - %s", f.Description)
			}

			if f.Value != nil && f.Value != "" {
				fmt.Fprintf(b, " (default %v)", f.Value)
			}

			fmt.Fprintf(b, "\n")
		}
	}

	_, err = b.WriteTo(c.Stdout)
	if err != nil {
		return nil, errors.Wrap(err, "write")
	}

	return nil, ErrExit
}
