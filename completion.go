package app

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var CompletionScript = `
# %[1]s bash completion function
_%[1]s_complete() {
	cmd=$(${COMP_WORDS[0]} --_comp-bash --_comp-word ${COMP_CWORD} --_comp-point ${COMP_POINT} --_comp-line "${COMP_LINE}" "${COMP_WORDS[@]:1:$COMP_CWORD}")
	echo $cmd >comp-debug2
	eval $cmd
	#COMPREPLY=($(eval $cmd))
	echo ${COMPREPLY[@]} >>comp-debug2
}

complete -F _%[1]s_complete %[1]s

# to persist bash completion add this to the end of your ~/.bashrc file by command:
#   %[2]s >>~/.bashrc
# or alternatively to enable it to only current session use command:
#   . <(%[2]s)
`

var (
	Complete      bool
	NLastComplete int

	CompleteBash  = F{Name: "_comp-bash"}.NewBool(false)
	CompleteLine  = F{Name: "_comp-line"}.NewString("")
	CompleteWord  = F{Name: "_comp-word"}.NewInt(0)
	CompletePoint = F{Name: "_comp-point"}.NewInt(0)

	CompleteCommand = &Command{
		Name:   "_comp-bash",
		Action: PrintCompletionScript,
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

func CheckCompletionAction(f Flag, _ *Command) error {
	if !CompleteBash.Value {
		return nil
	}
	Complete = CompleteBash.Value
	NLastComplete = 1
	return nil
}

func CompleteLast(args []string) (bool, string) {
	if !CompleteBash.Value {
		return false, ""
	}
	return true, args[len(args)-NLastComplete]
}

func CompleteDefault(last string) error {
	fmt.Fprintf(Writer, "COMPREPLY=($(compgen -o default \"%s\"))", last)
	return nil
}

var debug bytes.Buffer

func DefaultCommandComplete(c *Command) error {
	defer func() {
		_ = ioutil.WriteFile("comp-debug", debug.Bytes(), 0644)
	}()

	fmt.Fprintf(&debug, "%q\n", os.Args)
	fmt.Fprintf(&debug, "command: %s\n", c.Name)
	fmt.Fprintf(&debug, "NLastComplete: %d\n", NLastComplete)
	char := ">>>"
	if CompletePoint.Value < len(CompleteLine.Value) {
		char = fmt.Sprintf("'%c'", CompleteLine.Value[CompletePoint.Value])
	}
	fmt.Fprintf(&debug, "complete word %v pos %d %v line %q last %q\n", CompleteWord.Value, CompletePoint.Value, char, CompleteLine.Value, c.Args())

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

	mw := io.MultiWriter(Writer, &debug)

	fmt.Fprintf(mw, `COMPREPLY=($(compgen -W "%s"))`, strings.Join(names, " "))

	return nil
}

func FileFlagCompleteFunc(f Flag, _ *Command, last string) error {
	fmt.Fprintf(Writer, `compgen -o default "%s"`, last)
	return nil
}

func PrintCompletionScript(c *Command) error {
	fmt.Printf(CompletionScript, path.Base(os.Args[0]), strings.Join(os.Args, " "))
	return nil
}
