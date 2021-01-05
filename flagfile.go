package cli

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var fopen func(string) (io.ReadCloser, error) = func(n string) (io.ReadCloser, error) { return os.Open(n) }

var FlagfileFlag = &Flag{
	Name:        "flagfile,ff",
	Description: "load flags from file",
	After:       flagfile,
	Value:       StringSlicePtr(nil),
}

func flagfile(f *Flag, c *Command) error {
	args, err := func() (args Args, err error) {
		var ifexists bool
		fnames := f.Value.(*[]string)
		last := len(*fnames) - 1
		fname := (*fnames)[last]
		if strings.HasSuffix(fname, "?") {
			fname = strings.TrimSuffix(fname, "?")
			(*fnames)[last] = fname
			ifexists = true
		}

		f, err := fopen(fname)
		if os.IsNotExist(err) && ifexists {
			*fnames = (*fnames)[:last]
			return nil, nil
		}
		if err != nil {
			return
		}
		defer func() {
			if e := f.Close(); err == nil {
				err = e
			}
		}()

		r := bufio.NewScanner(f)
		r.Split(bufio.ScanWords)

		for r.Scan() {
			args = append(args, r.Text())
		}

		if err = r.Err(); err != nil {
			return
		}

		return
	}()
	if err != nil {
		return err
	}

	for args.Len() != 0 {
		args, err = c.parseFlag(args[0], args)
		if err != nil {
			return err
		}
	}

	return nil
}

func StringPtr(s string) *string { return &s }

func StringSlicePtr(s []string) *[]string { return &s }
