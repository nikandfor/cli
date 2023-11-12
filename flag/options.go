package flag

import (
	"os"

	"tlog.app/go/errors"
)

func Default(v interface{}) Option {
	return func(f *Flag) {
		f.Value = v
	}
}

func Hidden(f *Flag) {
	f.Hidden = true
}

func Required(f *Flag) {
	f.Required = true
}

func Local(f *Flag) {
	f.Local = true
}

// AtFile replaces flag value of the form @file with the file contents.
func AtFile(f *Flag) {
	orig := f.Action

	f.Action = func(f *Flag, arg string, args []string) ([]string, error) {
		key, val, args, err := ParseArg(arg, args, true, false)
		if err != nil {
			return args, err
		}

		if len(val) != 0 && val[0] == '@' {
			data, err := readFile(val[1:])
			if err != nil {
				return nil, errors.Wrap(err, "read file")
			}

			arg = key + "=" + string(data)
		}

		return orig(f, arg, args)
	}
}

var readFile = os.ReadFile
