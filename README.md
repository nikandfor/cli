[![Documentation](https://pkg.go.dev/badge/nikand.dev/go/cli)](https://pkg.go.dev/nikand.dev/go/cli?tab=doc)
[![Go workflow](https://github.com/nikandfor/cli/actions/workflows/go.yml/badge.svg)](https://github.com/nikandfor/cli/actions/workflows/go.yml)
[![CircleCI](https://circleci.com/gh/nikandfor/cli.svg?style=svg)](https://circleci.com/gh/nikandfor/cli)
[![codecov](https://codecov.io/gh/nikandfor/cli/branch/master/graph/badge.svg)](https://codecov.io/gh/nikandfor/cli)
[![Go Report Card](https://goreportcard.com/badge/nikand.dev/go/cli)](https://goreportcard.com/report/nikand.dev/go/cli)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/nikandfor/cli?sort=semver)

# cli
This is a lightweight yet extensible library for creating convenient command-line applications in Go.
It follows the general principle of being dead simple, efficient, and highly customizable to fit your needs.

Many cool features are not supported by the library, but they can be easily implemented if needed.
For example GNU grouping of oneletter options are not supported, but if you want, just set your custom `Command.ParseFlag` handler.

## Usage

### Hello World

```go
func main() {
    app := &cli.Command{
        Name: "hello",
        Action: hello, // this is how the library connects to your business logic
    }

    cli.RunAndExit(app, os.Args, nil)
}

func hello(c *cli.Command) error {
    fmt.Println("Hello world")

    return nil
}
```

### Flags

Flags are much similar to commands, they have Action property too which does all the job.
The default action is to parse the value and assign it to `(&Flag{}).Value`.

If flag is mentioned multiple times Action will be called for each occurance.
Flag value can be taken from env variables, which uses the same approach with the same Action.
But it's called before that command arguments are parsed.

```go
func main() {
    app := &cli.Command{
        Name: "hello",
        Action: hello,
        EnvPrefix: "HELLO_",
        Flags: []*cli.Flag{
            // much less boilerplate than &cli.StringFlag{Name: ...}
            // this was one of the main reasons I started this lib
            cli.NewFlag("name", "world", "name to say hello to"),

            // supported flag value types are string, some ints, time.Duration
            cli.NewFlag("full-flag-name,flag,f", 3, "alias names are added with comma"),

            // action can be passed instead of value
            cli.NewFlag("json,j", jsonFlagAction, "json encoded value"),

            // stdlib flag.Value is also supported

            cli.FlagfileFlag, // configs without configs
            cli.EnvFlag,
            cli.HelpFlag, // must be added explicitly
        },
    }

    cli.RunAndExit(app, os.Args, os.Environ()) // env as well as args are passed here
}

func hello(c *cli.Command) error {
    fmt.Println("Hello", c.String("name"))

    return nil
}

// f is a Flag we are working on now.
// arg is an arg that triggered the flag.
// args is the rest of arguments.
// Command can be received as f.CurrentCommand.(*cli.Command).
// Conciseness is sacrificed here in favor of extracting flags into its own package.
// Rest of unparsed args must be returned back.
func jsonFlagAction(f *cli.Flag, arg string, args []string) ([]string, error) {
    key, val, args, err := flag.ParseArg(arg, args, true, false)
                           // last two booleans:
                           // eat next arg from args if value is not present in arg
                           // allow no value at all. Like bool flags.
    if err != nil {
        return args, err
    }

    // key
    // We may slightly change parsing behaviour depending on a flag name was used.

    // default value if any now in f.Value

    err = json.Unmarshal([]byte(val), &f.Value)
    if err != nil {
        return nil, errors.Wrap(err, "parse value as json")
    }

    // You can actually return more args than if was.
    // --flagfile is working just like that.
    return args, nil
}

func hello(c *cli.Command) error {
    fmt.Println("Hello object", c.Flag("flag").Value)

    return nil
}
```

### Arguments

Command do not accept arguments by default. This saved me multiple times from doing something I wasn't going to ask for.
For example
```
command --flag= value           # flag="", args=["value"]
VAR="a b" command --flag=$VAR   # flag=a, args=["b"]
command -f --another value      # if f is a string flag it would be f="--another", args=["value"]
```
You'll get error in all the previous cases if you wasn't expecting arguments.

So if you need arguments, do it explicit.

```go
func main() {
    app := &cli.Command{
        Name: "hello",
        Args: cli.Args{}, // non-nil value must be here
        Action: hello,
    }

    cli.RunAndExit(app, os.Args, nil)
}

func hello(c *cli.Command) error {
    fmt.Println("Hello ", c.Args.Get(0)) // Get is, by the way, a lazy c.Args[0] if you don't want to check for len(c.Args)

    return nil
}
```

### Subcommands

```go
func main() {
    app := &cli.Command{
        Name:        "todo",
        Description: "todo list. add, list and remove jobs to do",
        // Command may not have action on their own. Help would be printed by default.
        Flags: []*cli.Flag{
            cli.NewFlag("storage", "todo.json", "file to store tasks"), // global flag
            cli.HelpFlag,
        },
        Commands: []*cli.Command{{
            Name: "add",
            Action: add,
            Args: cli.Args{}, // we'll take job description from args
        }, {
            Name: "remove,rm",
            Action: remove,
            Flags: []*cli.Flag{
                cli.NewFlag("job", -1, "job id to remove"), // local flag
            },
        }},
    }

    cli.RunAndExit(app, os.Args, nil)
}
```

### Flag values from the environment

```go
func main() {
    app := &cli.Command{
        Name: "hello",
        Action: hello,
        EnvPrefix: "HELLO_",
        Flags: []*cli.Flag{
            cli.NewFlag("flag,f", "", "flag description"),
        },
        Commands: []*cli.Command{{
            Name: "subcommand",
            // EnvPrefix could be overwritten here for nested flags
            // The parent is used if empty
            Action: subcommand,
            Flags: []*cli.Flag{
                cli.NewFlag("another,a", 1, "another description"),
            },
        }},
    }

    cli.RunAndExit(app, os.Args, os.Environ()) // ENV VARIABLE COME FROM HERE, not from os package. Explicitness.
}
```
Usage
```
HELLO_FLAG=value hello
HELLO_FLAG=v2 HELLO_ANOTHER=4 hello subcommand
```

#### Order of precedence

* --flag=first
* ENV_FLAG=second
* cli.NewFlag("flag", "the_last", "help")
