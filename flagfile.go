package cli

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/nikandfor/errors"
)

var FlagfileFlag = &Flag{
	Name:        "flagfile,ff",
	Description: "load flags from file",
	Action:      flagfile,
}

var readFile func(string) ([]byte, error) = os.ReadFile

func flagfile(c *Command, f *Flag, arg string, args []string) (_ []string, err error) {
	args, err = ParseFlagString(c, f, arg, args)
	if err != nil {
		return nil, err
	}

	data, err := readFile(f.Value.(string))
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	r := bufio.NewScanner(bytes.NewReader(data))
	r.Split(bufio.ScanLines)

	var add []string

	for r.Scan() {
		e := r.Text()
		e = strings.TrimSpace(e)
		if strings.HasPrefix(e, "#") {
			continue
		}

		add = append(add, e)
	}

	if err = r.Err(); err != nil {
		return nil, errors.Wrap(err, "scan file")
	}

	return append(add, args...), nil
}

func StringPtr(s string) *string { return &s }

func StringSlicePtr(s []string) *[]string { return &s }
