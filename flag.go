package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	Flag struct {
		Name        string
		Group       string
		Usage       string
		Description string
		Help        string

		Action FlagAction

		Hidden    bool
		Mandatory bool

		Value interface{}

		IsSet bool
	}

	// Setter is subset of stdlib flag.Value interface
	Setter interface {
		Set(v string) error
	}

	FlagAction func(c *Command, f *Flag, arg string, args []string) ([]string, error)

	FlagOption = func(f *Flag)
)

var (
	ErrFlagExit           = errors.New("flag exit")
	ErrFlagMandatory      = errors.New("flag is mandatory")
	ErrFlagValueRequired  = errors.New("flag value required")
	ErrFlagActionRequired = errors.New("flag action required")
	ErrNoSuchFlag         = errors.New("no such flag")
)

func NewFlag(name string, val interface{}, help string, opts ...FlagOption) (f *Flag) {
	f = &Flag{
		Name:        name,
		Description: help,

		Value: val,
	}

	switch val := val.(type) {
	case FlagAction:
		f.Value = nil
		f.Action = val
	case func(c *Command, f *Flag, arg string, args []string) ([]string, error):
		f.Value = nil
		f.Action = FlagAction(val)
	case bool:
		f.Action = ParseFlagBool
	case time.Duration:
		f.Action = ParseFlagDuration
	case float64:
		f.Action = ParseFlagFloat64
	case float32:
		f.Action = ParseFlagFloat32
	case int:
		f.Action = ParseFlagInt
	case uint:
		f.Action = ParseFlagUint
	case int64:
		f.Action = ParseFlagInt64
	case uint64:
		f.Action = ParseFlagUint64
	case string:
		f.Action = ParseFlagString
	case Setter:
		f.Action = ParseFlagValue(val, true, false)
	default:
		panic(fmt.Sprintf("unsupported value type: %T", val))
	}

	return f
}

func (f *Flag) MainName() string {
	return MainName(f.Name)
}

func DefaultParseFlag(c *Command, arg string, args []string) (nextArgs []string, err error) {
	st := 0
	for st < len(arg) && arg[st] == '-' {
		st++
	}

	end := st
	for end < len(arg) && arg[end] != '=' && arg[end] != ' ' {
		end++
	}

	f := c.Flag(arg[st:end])
	if f == nil {
		return nil, ErrNoSuchFlag
	}

	if f.Action == nil {
		return nil, ErrFlagActionRequired
	}

	return f.Action(c, f, arg, args)
}

// typed flag parsers

func ParseFlagBool(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	val, args, err := ParseFlagVal(arg, args, false, true)
	if err != nil {
		return nil, err
	}

	val = strings.ToLower(val)

	switch val {
	case "true", "t", "yes", "y", "":
		f.Value = true
	case "false", "f", "no", "n":
		f.Value = false
	default:
		return nil, errors.New("not a bool value")
	}

	f.IsSet = true

	return args, nil
}

func ParseFlagDuration(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFlagFunc(func(val string) (_ interface{}, err error) {
		v, err := time.ParseDuration(val)
		if err != nil {
			return nil, err
		}

		return v, nil
	}, true, false)

	return act(c, f, arg, args)
}

func ParseFlagFloat64(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFlagFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}

		return v, nil
	}, true, false)

	return act(c, f, arg, args)
}

func ParseFlagFloat32(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFlagFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, err
		}

		return float32(v), nil
	}, true, false)

	return act(c, f, arg, args)
}

func ParseFlagInt(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFlagFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}

		return int(v), nil
	}, true, false)

	return act(c, f, arg, args)
}

func ParseFlagInt64(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFlagFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}

		return v, nil
	}, true, false)

	return act(c, f, arg, args)
}

func ParseFlagString(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFlagFunc(func(val string) (_ interface{}, err error) {
		return val, nil
	}, true, false)

	return act(c, f, arg, args)
}

func ParseFlagUint(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFlagFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}

		return uint(v), nil
	}, true, false)

	return act(c, f, arg, args)
}

func ParseFlagUint64(c *Command, f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFlagFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}

		return v, nil
	}, true, false)

	return act(c, f, arg, args)
}

//

func ParseFlagValue(v Setter, eatnext, optional bool) FlagAction {
	return ParseFlagFunc(func(val string) (_ interface{}, err error) {
		err = v.Set(val)
		if err != nil {
			return
		}

		return v, nil
	}, eatnext, optional)
}

func ParseFlagFunc(p func(string) (interface{}, error), eatnext, optional bool) FlagAction {
	return func(c *Command, f *Flag, arg string, args []string) (_ []string, err error) {
		val, args, err := ParseFlagVal(arg, args, eatnext, optional)
		if err != nil {
			return args, err
		}

		v, err := p(val)
		if err != nil {
			return nil, err
		}

		f.Value = v
		f.IsSet = true

		return args, nil
	}
}

//

func ParseFlagVal(arg string, args []string, eatnext, optional bool) (val string, nextargs []string, err error) {
	_, val, _, nextargs, err = ParseFlagArg(arg, args, eatnext, optional)
	return
}

func ParseFlagArg(arg string, args []string, eatnext, optional bool) (k, val string, dashes int, _ []string, err error) {
	for dashes < len(arg) && arg[dashes] == '-' {
		dashes++
	}

	end := dashes
	for end < len(arg) && arg[end] != '=' && arg[end] != ' ' {
		end++
	}

	k = arg[dashes:end]

	switch {
	case end < len(arg):
		vst := end
		if vst < len(arg) && (arg[vst] == '=' || arg[vst] == ' ') {
			vst++
		}

		val = arg[vst:]
	case eatnext && len(args) != 0:
		val = args[0]
		args = args[1:]
	case optional:
		//
	default:
		err = ErrFlagValueRequired
		return
	}

	return k, val, dashes, args, nil
}

func (f *Flag) check() error {
	if f.Mandatory && !f.IsSet {
		return ErrFlagMandatory
	}

	return nil
}
