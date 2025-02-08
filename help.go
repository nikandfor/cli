package cli

import (
	"bytes"
	"fmt"
	"strings"

	"nikand.dev/go/cli/flag"
)

var HelpFlag = &Flag{
	Name:        "help,h",
	Usage:       "=[hidden]",
	Description: "print command help end exit",
	Action:      defaultHelp,
}

func defaultHelp(f *Flag, arg string, args []string) (rest []string, err error) {
	const minNameW, maxNameW = 20, 40

	c := f.CurrentCommand.(*Command)

	_, v, rest, err := flag.ParseArg(arg, args, false, true)
	if err != nil {
		return
	}

	hidden := v == "hidden"

	b := new(bytes.Buffer)

	pline := func(name, usage, desc string, w int, val interface{}) {
		name += usage

		fmt.Fprintf(b, "    %-*s", w, name)

		if len(name) > w {
			fmt.Fprintf(b, "\n    %-*s", w, "")
		}

		lines := strings.Split(desc, "\n")

		for i, l := range lines {
			if i == 0 {
				fmt.Fprintf(b, " - ")
			} else {
				fmt.Fprintf(b, "\n    %-*s   ", w, "")
			}

			fmt.Fprintf(b, "%s", l)
		}

		if val != nil && val != "" {
			if len(lines) != 0 && lines[len(lines)-1] == "" {
				fmt.Fprintf(b, "default %v", val)
			} else {
				fmt.Fprintf(b, " (default %v)", val)
			}
		}

		fmt.Fprintf(b, "\n")
	}

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
		namew := minNameW

		for _, sub := range c.Commands {
			if sub == nil {
				// spacing
			} else if sub.Hidden && !hidden {
				continue
			} else if w := len(sub.Name) + len(sub.Usage); w > namew {
				namew = w
			}

			cnt++
		}

		if cnt != 0 {
			fmt.Fprintf(b, "\nSubcommands\n")
		}

		headernl := false

		if namew > maxNameW {
			namew = maxNameW
		}

		for _, sub := range c.Commands {
			switch {
			case sub == nil:
				fmt.Fprintf(b, "\n")
				continue
			case sub.Name == "":
				if headernl {
					fmt.Fprintf(b, "\n")
					headernl = false
				}

				fmt.Fprintf(b, "    %*s # %s\n", namew, "", sub.Description)
				continue
			case sub.Hidden && !hidden:
				continue
			}

			headernl = true

			pline(sub.Name, "", sub.Description, namew, nil)
		}
	}

	for cc := c; cc != nil; cc = cc.Parent {
		cnt := 0
		namew := minNameW

		for _, f := range cc.Flags {
			if f == nil {
				// spacing
			} else if f.Hidden && !hidden || cc != c && f.Local {
				continue
			} else if w := len(f.Name) + len(f.Usage); w > namew {
				namew = w
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

		headernl := false

		if namew > maxNameW {
			namew = maxNameW
		}

		for _, f := range cc.Flags {
			switch {
			case f == nil:
				fmt.Fprintf(b, "\n")
				continue
			case f.Name == "":
				if headernl {
					fmt.Fprintf(b, "\n")
					headernl = false
				}

				fmt.Fprintf(b, "    %*s # %s\n", namew, "", f.Description)
				continue
			case f.Hidden && !hidden || cc != c && f.Local:
				continue
			}

			headernl = true

			pline(f.Name, f.Usage, f.Description, namew, f.Value)
		}
	}

	_, err = b.WriteTo(c.Stdout)
	if err != nil {
		return nil, wrap(err, "write")
	}

	return nil, ErrExit
}
