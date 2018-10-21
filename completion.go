package app

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var CompletionScript = `
# %[1]s bash completion function
_%[1]s_complete() {
	local base=$1 cur=$2 prev=$3
	cmd=$(${base} --_comp-bash --_comp-word ${COMP_CWORD} --_comp-point ${COMP_POINT} --_comp-line "${COMP_LINE}" "${COMP_WORDS[@]:1:$COMP_CWORD}")
	eval "$cmd"
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

func CompletionLast(args []string) (bool, string) {
	if !CompleteBash.Value {
		return false, ""
	}
	return true, args[len(args)-NLastComplete]
}

func NoArgumentsExpectedCompletion(c *Command) error {
	fmt.Fprintf(Writer, `COMPREPLY=("%s" " "); compopt -o nosort`, "no arguments expected")
	return nil
}

func AlternativeCompletion(list []string) Action {
	return func(c *Command) error {
		fmt.Printf(`mapfile -t COMPREPLY < <(grep "^$cur" <<EOF
%s
EOF
)
[ ${#COMPREPLY[@]} -eq 1 ] && { a="${COMPREPLY[0]}"; [[ "$a" =~ " " ]] && COMPREPLY[0]=$(printf '"%%s"' "$a"); }`, strings.Join(list, "\n"))
		return nil
	}
}

func DefaultBashCompletion(last string) error {
	fmt.Fprintf(Writer, "_longopt")
	return nil
}

var DefaultCommandCompletion = func(c *Command) error {
	pref := c.Args().Last()
	if len(c.Commands) == 0 && !strings.HasPrefix(pref, "-") {
		return DefaultBashCompletion(pref)
	}

	var names []string

	if !strings.HasPrefix(pref, "-") {
		for _, s := range c.Commands {
			if pref == "" && s.Name[0] == '_' {
				continue
			}
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
			if strings.TrimLeft(pref, "-") == "" && b.Name[0] == '_' {
				continue
			}
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

	fmt.Fprintf(Writer, `COMPREPLY=($(compgen -W "%s")); compopt -o nosort`, strings.Join(names, " "))

	return nil
}

var DefaultFlagCompletion = func(f Flag, c *Command, last string) error {
	switch f.(type) {
	case *FileFlag:
		fmt.Fprintf(Writer, `_longopt`)
	}
	if last != "" {
		return nil
	}

	var msg string
	if h := f.Base().CompletionHelp; h != "" {
		msg = h
	} else {
		tp := ""
		switch f.(type) {
		case *StringFlag:
			tp = "string"
		case *StringSliceFlag:
			tp = "string"
		case *IntFlag:
			tp = "int"
		default:
			tp = fmt.Sprintf("%T", f)
		}
		msg = fmt.Sprintf("%s argument expected", tp)
	}
	fmt.Fprintf(Writer, `COMPREPLY=(%s); compopt -o nosort`, msg)
	return nil
}

func PrintCompletionScript(c *Command) error {
	fmt.Printf(CompletionScript, path.Base(os.Args[0]), strings.Join(os.Args, " "))
	return nil
}
