package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nikandfor/cli/complete"
	"github.com/nikandfor/errors"
	"github.com/nikandfor/tlog"
)

var ErrCouldNotDetermineShell = errors.New("couldn't determine the shell")

var CompleteCmd = &Command{
	Name:        "_complete",
	Usage:       "[bash|zsh]",
	Description: "print shell completion script",
	Action:      completeAuto,
	Hidden:      true,
	Commands: []*Command{{
		Name:   "bash",
		Action: completeBash,
	}, {
		Name:   "zsh",
		Action: completeZsh,
	}},
}

func beforeComplete(c *Command) (err error) {
	fn := "complete_log_app"
	fn, err = filepath.Abs(fn)
	if err != nil {
		return errors.Wrap(err, "abs path")
	}

	f, err := os.Create(fn)
	if err != nil {
		return errors.Wrap(err, "create log file")
	}

	defer f.Close()

	for _, e := range c.Env {
		if strings.HasPrefix(e, "CLI_COMP_") {
			fmt.Fprintf(f, "%v\n", e)
		}
	}

	//	return errors.New("we were here")
	return nil
}

func Complete(c *Command) (err error) {
	var repl []string

	current := complete.Current(c)

	if false && current == "" {
		_, err = defaultHelp(c, nil, "", nil)
		return
	}

	var dashes string
	cur := current
	{
		i := 0
		for i < len(cur) && cur[i] == '-' {
			i++
		}

		dashes = cur[:i]
		cur = cur[i:]
	}

	if dashes == "" {
		for _, sub := range c.Commands {
		cmd:
			for _, name := range strings.Split(sub.Name, ",") {
				if strings.HasPrefix(name, "_") != strings.HasPrefix(cur, "_") {
					continue
				}

				if strings.HasPrefix(name, cur) {
					repl = append(repl, name)

					break cmd
				}
			}
		}
	} else {
		for _, f := range c.Flags {
		flg:
			for _, name := range strings.Split(f.Name, ",") {
				if (len(dashes) > 1) && (len(name) == 1) {
					continue
				}

				if strings.HasPrefix(name, cur) {
					dd := "--"
					if len(name) == 1 {
						dd = "-"
					}

					repl = append(repl, dd+name)

					break flg
				}
			}
		}
	}

	tlog.Printw("complete", "dashes", dashes, "cur", cur, "opts", repl)

	for i := range repl {
		repl[i] = strconv.Quote(repl[i])
	}

	fmt.Fprintf(c.Stdout, "COMPREPLY=( $(compgen -W '%[2]s' -- %[1]s) )", current, strings.Join(repl, " "))

	//	fmt.Fprintf(c, "COMPREPLY=(%s)", strings.Join(repl, " "))

	return nil
}

func completeAuto(c *Command) error {
	sh, ok := complete.Shell(c)
	if !ok {
		return ErrCouldNotDetermineShell
	}

	root := c
	for root.Parent != nil {
		root = root.Parent
	}

	return complete.ExecTemplate(c.Stdout, sh, root.OSArgs)
}

func completeBash(c *Command) error {
	root := c
	for root.Parent != nil {
		root = root.Parent
	}

	return complete.ExecTemplate(c.Stdout, "bash", root.OSArgs)
}

func completeZsh(c *Command) error {
	root := c
	for root.Parent != nil {
		root = root.Parent
	}

	return complete.ExecTemplate(c.Stdout, "zsh", root.OSArgs)
}
