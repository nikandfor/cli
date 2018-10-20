package app

import (
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

func DefaultCommandComplete(c *Command) error {
	return nil
}
