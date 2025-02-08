package cli

import (
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"

	"nikand.dev/go/cli/flag"
)

// FlagfileFlag replaces this flag occurrence with the given file content split on spaces.
// # comments are also supported.
var FlagfileFlag = &Flag{
	Name:        "flagfile,ff",
	Description: "load flags from file",
	Action:      flagfile,
}

var readFile = os.ReadFile

func flagfile(f *Flag, arg string, args []string) (_ []string, err error) {
	_, val, args, err := flag.ParseArg(arg, args, true, false)
	if err != nil {
		return nil, err
	}

	d, err := readFile(val)
	if err != nil {
		return nil, wrap(err, "read file")
	}

	var add []string
	var buf []byte

	for i := 0; i < len(d); i++ {
		i = skip(d, i, unicode.IsSpace)
		if i == len(d) {
			break
		}

		if d[i] == '#' {
			i = skip(d, i, untilNewline)
			continue
		}

		buf, i, err = decodeArg(d, i, buf[:0])
		if err != nil {
			return nil, err
		}

		add = append(add, string(buf))
	}

	return append(add, args...), nil
}

func decodeArg(d []byte, i int, buf []byte) ([]byte, int, error) {
	done := i
	var esc, single, double bool

	flush := func(w int) {
		buf = append(buf, d[done:i]...)
		done = i + w
	}

loop:
	for w := 0; i < len(d); i += w {
		var r rune
		r, w = utf8.DecodeRune(d[i:])

		switch {
		case unicode.IsSpace(r) && !double && !single && !esc:
			break loop
		case esc:
			esc = !esc
		case r == '\\' && !single:
			flush(w)
			esc = true
		case r == '"' && !single:
			flush(w)
			double = !double
		case r == '\'' && !double:
			flush(w)
			single = !single
		}
	}
	if esc || double || single {
		return buf, i, fmt.Errorf("bad string")
	}

	flush(0)

	return buf, i, nil
}

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
