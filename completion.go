package app

import (
	"fmt"
	"strings"
)

var (
	Complete      bool
	NLastComplete int

	CompleteBash  = F{Name: "_comp-bash"}.NewBool(false)
	CompleteLine  = F{Name: "_comp-line"}.NewString("")
	CompleteWord  = F{Name: "_comp-word"}.NewInt(0)
	CompletePoint = F{Name: "_comp-point"}.NewInt(0)

	CompleteCommand = &Command{
		Name: "_comp-bash",
	}
)

func init() {
	CompleteBash.After = CheckCompletionAction
	CompleteLine.After = CheckCompletionAction
	CompleteWord.After = CheckCompletionAction
	CompletePoint.After = CheckCompletionAction
}

func EnableCompletion() {
	AddCompletionToApp(App)
}

func AddCompletionToApp(app *Command) {
	app.Commands = append(app.Commands,
		CompleteCommand,
	)
	app.Flags = append(app.Flags,
		CompleteBash,
		CompleteLine,
		CompleteWord,
		CompletePoint,
	)
}

func CheckCompletionAction(f Flag) error {
	if !CompleteBash.Value {
		return nil
	}
	Complete = CompleteBash.Value
	if CompleteWord.Value == len(strings.Fields(CompleteLine.Value)) {
		NLastComplete = 1
		return nil
	}
	NLastComplete = 0
	return nil
}

func CompleteLast(args []string) (bool, string) {
	if !CompleteBash.Value {
		return false, ""
	}
	if CompleteWord.Value == len(strings.Fields(CompleteLine.Value)) {
		return true, args[len(args)-1]
	}

	return true, ""
}

func CompleteDefault(last string) error {
	fmt.Fprintf(Writer, "compgen -o default \"%s\"", last)
	return nil
}

func DefaultCommandComplete(c *Command) error {
	pref := c.Args().Last()
	if len(c.Commands) == 0 && !strings.HasPrefix(pref, "-") {
		return CompleteDefault(pref)
	}

	var names []string

	if !strings.HasPrefix(pref, "-") {
		for _, s := range c.Commands {
			if strings.HasPrefix(s.Name, pref) {
				names = append(names, s.Name)
				continue
			}
			if pref == "" {
				continue
			}
			for _, a := range s.Aliases {
				if !strings.HasPrefix(a, pref) {
					continue
				}
				names = append(names, a)
			}
		}
	}

	if pref == "" || strings.HasPrefix(pref, "-") {
		addflag := func(n string) bool {
			if len(n) == 1 {
				n = "-" + n
			} else {
				n = "--" + n
			}
			if !strings.HasPrefix(n, pref) {
				return false
			}
			names = append(names, n)
			return true
		}
		for _, f := range c.Flags {
			b := f.Base()
			if addflag(b.Name) {
				continue
			}
			if pref == "" {
				continue
			}
			for _, a := range b.Aliases {
				addflag(a)
			}
		}
	}

	fmt.Fprintf(Writer, `compgen -W "%s"`, strings.Join(names, " "))

	return nil
}
