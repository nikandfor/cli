package app

import (
	"fmt"
)

func ExampleApp() {
	App.Name = "exampleapp"
	App.Description = "This is the simples command example"
	App.Action = func(c *Command) error {
		fmt.Printf("flags are: file %q, int %v\n", c.Flag("file").VAny(), c.Flag("int").VAny())
		fmt.Printf("arguments: %q\n", c.Args())
		return nil
	}
	App.Commands = nil
	App.Flags = []Flag{
		F{Name: "file", Description: "file to store something"}.NewFile("default.txt"),
		F{Name: "int", Description: "number of pens you have now"}.NewInt(0),
		HelpFlag,
	}

	err := App.Run([]string{"exampleapp", "--int=3", "--file", "another.txt", "some", "args"}) // os.Args normally
	if err != nil {
		panic(err)
	}
}

func ExampleCompletion() {
	App.Name = "exampleapp"
	App.Description = "This is the simples command example"
	App.Action = func(c *Command) error {
		fmt.Printf("flags are: file %q, int %v\n", c.Flag("file").VAny(), c.Flag("int").VAny())
		fmt.Printf("arguments: %q\n", c.Args())
		return nil
	}
	App.Commands = nil
	App.Flags = []Flag{
		F{Name: "file", Description: "file to store something"}.NewFile("default.txt"),
		F{Name: "int", Description: "number of pens you have now"}.NewInt(0),
		HelpFlag,
	}

	EnableCompletion()

	// First you need to tell you bash environment to set our custom completion function to application command
	// It's done by this command
	// . <(exampleapp _comp-bash)

	err := App.Run([]string{"exampleapp",
		"--_comp-bash", // this flag is set by bash if completion script was set up
		// some other flags are set also but they are not important for us now
		"--int=3", "--f"})
	if err != nil {
		panic(err)
	}

	// So here you type command:
	// exampleapp --int=3 --f<Tab><Tab>
	// DefaultCommandCompletion function executed and offers you --file flag
	// Since that the only option with such prefix it will be substituted into command line automatically.
	// And you'll get this in you command line
	// exampleapp --int=3 --file // <- here after white space cursor will be waiting for you

	// Output: COMPREPLY=($(compgen -W "%s")); compopt -o nosort`, strings.Join(names, " "))

	// Function produced full command that defined COMPREPLY valiable explicit
	// So you can create any script you want with no limitations
	// Generated command will be evaluated info COMPREPLY variable set to only one option.
	// Bash uses that variable to change you command.
	// Last part disables options sorting so they appears in order they where defined and commands and flags are not mixed.
}
