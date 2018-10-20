package app

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var BashCompletionTemplate = `
# %[1]s bash completion function
_%[1]s_complete() {
	COMPREPLY=($(${COMP_WORDS[0]} --_complete-bash --_complete-word ${COMP_CWORD} --_complete-point ${COMP_POINT} --_complete-line "${COMP_LINE}" "${COMP_WORDS[@]:0:$COMP_CWORD}"))
}

complete -F _%[1]s_complete %[1]s

# to persist bash completion add this to the end of your ~/.bashrc file by command:
#   %[2]s >>~/.bashrc
# or alternatively to enable it to only current session use command:
#   . <(%[2]s)
`

var (
	BashCompletionCommand = &Command{
		Name:   "_completion-func-bash",
		Action: PrintBashCompletionFunc,
		Hidden: true,
	}

	BashCompletionFlag      = &BoolFlag{F: F{Name: "_complete-bash", Hidden: true}}
	BashCompletionWordFlag  = &IntFlag{F: F{Name: "_complete-word", Hidden: true}}
	BashCompletionPointFlag = &IntFlag{F: F{Name: "_complete-point", Hidden: true}}
	BashCompletionLineFlag  = &StringFlag{F: F{Name: "_complete-line", Hidden: true}}
)

func DefaultCommandCompletion(c *Command) error {
	last := c.Args().Last()

	fmt.Printf("echo \"\"")

	for _, s := range c.Commands {
		if s.Name[0] == '_' {
			continue
		}
		if strings.HasPrefix(s.Name, last) {
			fmt.Printf(" %v", s.Name)
		}
		for _, a := range s.Aliases {
			if strings.HasPrefix(a, last) {
				fmt.Printf(" %v", a)
			}
		}
	}

	for c := c; c != nil; c = c.parent {
		for _, f := range c.Flags {
			b := f.Base()
			if b.Name[0] == '_' {
				continue
			}
			pref := "-"
			if len(b.Name) > 1 {
				pref = "--"
			}
			if strings.HasPrefix(pref+b.Name, last) {
				fmt.Printf(" %v", pref+b.Name)
			}
			for _, a := range b.Aliases {
				if len(a) > 1 {
					pref = "--"
				} else {
					pref = "-"
				}
				if strings.HasPrefix(pref+a, last) {
					fmt.Printf(" %v", pref+a)
				}
			}
		}
	}

	if len(c.Commands) == 0 {
		fmt.Printf(" && compgen -o default")
	}

	fmt.Printf("\n")

	return nil
}

func DefaultFlagCompletion(f Flag) func(string) error {
	return func(l string) error {
		fmt.Printf("compgen -o default %v", l)
		return nil
	}
}

func IsCompletionSet() bool {
	return BashCompletionFlag.Value
}

func ifCompletion() (_ bool, _ int) {
	if !IsCompletionSet() {
		return false, 0
	}

	//	log.Printf("word: %v", BashCompletionWordFlag.Value)
	//	log.Printf("line: %q", BashCompletionLineFlag.Value)
	if BashCompletionWordFlag.Value == len(strings.Fields(BashCompletionLineFlag.Value)) {
		return true, 0
	}

	return true, 1
}

func PrintBashCompletionFunc(c *Command) error {
	arg0 := path.Base(os.Args[0])
	_, err := fmt.Fprintf(Writer, BashCompletionTemplate, arg0, strings.Join(os.Args[:2], " "))
	return err
}

func AddBashCompletion() {
	App.Commands = append(App.Commands, BashCompletionCommand)
	App.Flags = append(App.Flags,
		BashCompletionFlag,
		BashCompletionWordFlag,
		BashCompletionPointFlag,
		BashCompletionLineFlag,
	)
}
