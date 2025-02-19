package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"

	"nikand.dev/go/cli"
	"nikand.dev/go/cli/flag"
)

var (
	version = "dev"
	commit  = "dev"
	date    string
)

func main() {
	app := &cli.Command{
		Name:        "greeter",
		Usage:       "",
		Description: "greeter greets you",

		Help: `You may add some greetings and then run the program to greet you
with one of the greetings saved.

Greetings are golang 'text/template's with Command argument. Example greetings:
	"Hello, {{ .Flag "name" }}!!"
`,

		Before: before,
		Action: greet,

		EnvPrefix: "GREETER_",

		Flags: []*cli.Flag{
			cli.NewFlag("file", "greetings.txt", "file with greetings to use"),
			cli.NewFlag("name", "", "name to greet", flag.Required),

			cli.NewFlag("greeting,g", -1, "greeting to use to greet", flag.Local),

			cli.EnvfileFlag,
			cli.FlagfileFlag,
			cli.HelpFlag,
		},
		Commands: []*cli.Command{
			{
				Name:        "list,ls",
				Description: "list saved greetings",
				Action:      list,
			},
			{
				Name:        "add,new",
				Description: "add greeting",
				Args:        cli.Args{}, // if you expect arguments, you must specify it
				Action:      add,
				Flags: []*cli.Flag{
					cli.NewFlag("unique,uniq", false, "do not add if already exists"),
				},
			},
			{
				Name:        "delete,del,remove,rm",
				Description: "delete saved greeting",
				Action:      del,
			},
			{
				Name:        "exec",
				Description: "execute greeting template. useful for testing before adding",
				Args:        cli.Args{},
				Action:      exec,
			},
			cli.Version(version, commit, date),
			cli.CompleteCmd,
		},
	}

	cli.RunAndExit(app, os.Args, os.Environ())
}

func before(c *cli.Command) error {
	rand.Seed(time.Now().UnixNano())

	return nil
}

func greet(c *cli.Command) error {
	data, err := os.ReadFile(c.String("file"))
	if os.IsNotExist(err) {
		return errors.New("no greetings added")
	}
	if err != nil {
		return wrap(err, "read greetings file")
	}

	r := bytes.NewReader(data)

	d := json.NewDecoder(r)

	var raw json.RawMessage

	var choice int

	if f := c.Flag("greeting"); f.IsSet {
		choice = f.Int()
	} else {
		n := 0
		for ; d.More(); n++ {
			err = d.Decode(&raw)
			if err != nil {
				return wrap(err, "decode greeting")
			}
		}

		choice = rand.Intn(n)

		r.Reset(data)
	}

	for n := 0; d.More() && n <= choice; n++ {
		err = d.Decode(&raw)
		if err != nil {
			return wrap(err, "decode greeting")
		}
	}

	var text string

	err = json.Unmarshal(raw, &text)
	if err != nil {
		return wrap(err, "decode greeting")
	}

	t, err := template.New("greeting").Parse(text)
	if err != nil {
		return wrap(err, "parse template")
	}

	err = t.Execute(c.Stdout, c)
	if err != nil {
		return wrap(err, "execute template")
	}

	_, err = c.Stdout.Write([]byte{'\n'})
	if err != nil {
		return wrap(err, "write newline")
	}

	return nil
}

func add(c *cli.Command) (err error) {
	text := strings.Join(c.Args, " ")

	_, err = template.New("greeting").Parse(text)
	if err != nil {
		return wrap(err, "parse template")
	}

	f, err := os.OpenFile(c.String("file"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return wrap(err, "open greetings file")
	}

	defer func() {
		e := f.Close()
		if err == nil {
			err = wrap(e, "close greetings file")
		}
	}()

	e := json.NewEncoder(f)

	err = e.Encode(text)
	if err != nil {
		return wrap(err, "encode greeting")
	}

	return nil
}

func list(c *cli.Command) error {
	data, err := os.ReadFile(c.String("file"))
	if err != nil {
		return wrap(err, "read greetings file")
	}

	r := bytes.NewReader(data)
	d := json.NewDecoder(r)

	var raw json.RawMessage

	for n := 0; d.More(); n++ {
		err = d.Decode(&raw)
		if err != nil {
			return wrap(err, "decode greeting")
		}

		fmt.Fprintf(c.Stdout, "%d: %s\n", n, raw)
	}

	return nil
}

func del(c *cli.Command) error {
	data, err := os.ReadFile(c.String("file"))
	if err != nil {
		return wrap(err, "read greetings file")
	}

	r := bytes.NewReader(data)
	d := json.NewDecoder(r)

	var buf bytes.Buffer
	e := json.NewEncoder(&buf)

	var raw json.RawMessage

	skip := c.Int("greeting")

	for n := 0; d.More(); n++ {
		err = d.Decode(&raw)
		if err != nil {
			return wrap(err, "decode greeting")
		}

		if n == skip {
			continue
		}

		err = e.Encode(raw)
		if err != nil {
			return wrap(err, "encode greeting")
		}
	}

	err = os.WriteFile(c.String("file"), buf.Bytes(), 0644)
	if err != nil {
		return wrap(err, "write greetings file")
	}

	return nil
}

func exec(c *cli.Command) (err error) {
	text := strings.Join(c.Args, " ")

	t, err := template.New("greeting").Parse(text)
	if err != nil {
		return wrap(err, "parse template")
	}

	err = t.Execute(c.Stdout, c)
	if err != nil {
		return wrap(err, "execute template")
	}

	_, err = c.Stdout.Write([]byte{'\n'})
	if err != nil {
		return wrap(err, "write newline")
	}

	return nil
}

func wrap(err error, msg string) error {
	return fmt.Errorf("%v: %w", msg, err)
}
