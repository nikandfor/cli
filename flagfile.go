package cli

import (
	"os"
	"unicode"
	"unicode/utf8"

	"nikand.dev/go/cli/flag"
	"tlog.app/go/errors"
)

// FlagfileFlag replaces this flag occurrence with the given file content split on spaces.
// # comments are also supported.
var FlagfileFlag = &Flag{
	Name:        "flagfile,ff",
	Description: "load flags from file",
	Action:      flagfile,
}

var readFile = os.ReadFile

func skip(d []byte, i int, f func(r rune) bool) int {
	for w := 0; i < len(d); i += w {
		var r rune
		r, w = utf8.DecodeRune(d[i:])

		if !f(r) {
			return i
		}
	}

	return i
}

func untilNewline(r rune) bool { return r != '\n' }
func isArg(r rune) bool        { return !unicode.IsSpace(r) }

//	isArg := func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsPunct(r) }

func flagfile(f *Flag, arg string, args []string) (_ []string, err error) {
	_, val, args, err := flag.ParseArg(arg, args, true, false)
	if err != nil {
		return nil, err
	}

	d, err := readFile(val)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	var add []string

	for i := 0; i < len(d); i++ {
		i = skip(d, i, unicode.IsSpace)

		if d[i] == '#' {
			i = skip(d, i, untilNewline)
			continue
		}

		st := i
		i = skip(d, i, isArg)

		add = append(add, string(d[st:i]))
	}

	return append(add, args...), nil
}
