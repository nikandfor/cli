package flag

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

		Action Action  // flag parser
		Check  Visitor // called after all parsing but before command action for all flags
		//	Complete Visitor

		Hidden   bool // not shown in a help by default
		Required bool // must be set from args or env var
		Local    bool // do not inherited by child

		IsSet bool

		Value interface{}

		CurrentCommand interface{}
	}

	Action  func(f *Flag, arg string, args []string) ([]string, error)
	Visitor func(f *Flag) error
	Option  = func(f *Flag)

	// Setter is subset of stdlib flag.Value interface
	Setter interface {
		Set(v string) error
	}
)

var (
	ErrRequired      = errors.New("flag is required")
	ErrValueRequired = errors.New("flag value is required")
)

func New(name string, val interface{}, help string, opts ...Option) (f *Flag) {
	f = &Flag{
		Name:        name,
		Description: help,

		Value: val,
	}

	switch val := val.(type) {
	case Action:
		f.Value = nil
		f.Action = val
	case func(f *Flag, arg string, args []string) ([]string, error):
		f.Value = nil
		f.Action = Action(val)
	case bool:
		f.Action = ParseBool
	case time.Duration:
		f.Action = ParseDuration
	case float64:
		f.Action = ParseFloat64
	case float32:
		f.Action = ParseFloat32
	case int:
		f.Action = ParseInt
	case uint:
		f.Action = ParseUint
	case int64:
		f.Action = ParseInt64
	case uint64:
		f.Action = ParseUint64
	case string:
		f.Action = ParseString
	case []string:
		f.Action = ParseStringSlice
	case Setter:
		f.Action = ParseSetter(val, true, false)
	default:
		panic(fmt.Sprintf("unsupported value type: %T", val))
	}

	for _, o := range opts {
		o(f)
	}

	return f
}

func (f *Flag) MainName() string {
	p := strings.IndexByte(f.Name, ',')
	if p == -1 {
		return f.Name
	}

	return f.Name[:p]
}

func CheckFlag(f *Flag) error {
	if f.Check != nil {
		return f.Check(f)
	}

	if f.Required && !f.IsSet {
		return ErrRequired
	}

	return nil
}

// typed flag parsers

func ParseBool(f *Flag, arg string, args []string) ([]string, error) {
	_, val, args, err := ParseArg(arg, args, false, true)
	if err != nil {
		return nil, err
	}

	val = strings.ToLower(val)

	switch val {
	case "true", "t", "yes", "y", "", "1":
		f.Value = true
	case "false", "f", "no", "n", "0":
		f.Value = false
	default:
		return nil, errors.New("not a bool value")
	}

	f.IsSet = true

	return args, nil
}

func ParseDuration(f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFunc(func(val string) (_ interface{}, err error) {
		v, err := time.ParseDuration(val)
		if err != nil {
			return nil, err
		}

		return v, nil
	}, true, false)

	return act(f, arg, args)
}

func ParseFloat64(f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}

		return v, nil
	}, true, false)

	return act(f, arg, args)
}

func ParseFloat32(f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, err
		}

		return float32(v), nil
	}, true, false)

	return act(f, arg, args)
}

func ParseInt(f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}

		return int(v), nil
	}, true, false)

	return act(f, arg, args)
}

func ParseInt64(f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}

		return v, nil
	}, true, false)

	return act(f, arg, args)
}

func ParseString(f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFunc(func(val string) (_ interface{}, err error) {
		return val, nil
	}, true, false)

	return act(f, arg, args)
}

func ParseStringSlice(f *Flag, arg string, args []string) ([]string, error) {
	_, val, args, err := ParseArg(arg, args, true, false)
	if err != nil {
		return args, err
	}

	vals := strings.Split(val, ",")

	if !f.IsSet {
		f.IsSet = true
		f.Value = nil
	}

	if f.Value == nil {
		f.Value = vals
		return args, nil
	}

	have := f.Value.([]string)

	f.Value = append(have, vals...)

	return args, nil
}

func ParseUint(f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}

		return uint(v), nil
	}, true, false)

	return act(f, arg, args)
}

func ParseUint64(f *Flag, arg string, args []string) ([]string, error) {
	act := ParseFunc(func(val string) (_ interface{}, err error) {
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}

		return v, nil
	}, true, false)

	return act(f, arg, args)
}

//

func ParseSetter(v Setter, eatnext, optional bool) Action {
	return ParseFunc(func(val string) (_ interface{}, err error) {
		err = v.Set(val)
		if err != nil {
			return
		}

		return v, nil
	}, eatnext, optional)
}

func ParseFunc(p func(string) (interface{}, error), eatnext, optional bool) Action {
	return func(f *Flag, arg string, args []string) (_ []string, err error) {
		_, val, args, err := ParseArg(arg, args, eatnext, optional)
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

func ParseArg(arg string, args []string, eatnext, optional bool) (key, val string, nextargs []string, err error) {
	dashes := 0
	for dashes < 2 && dashes < len(arg) && arg[dashes] == '-' {
		dashes++
	}

	end := dashes
	for end < len(arg) && arg[end] != '=' && arg[end] != ' ' {
		end++
	}

	key = arg[dashes:end]

	switch {
	case end < len(arg):
		vst := end + 1
		val = arg[vst:]
	case eatnext && len(args) != 0:
		val = args[0]
		args = args[1:]
	case optional:
		// no value
	default:
		err = ErrValueRequired
		return
	}

	return key, val, args, nil
}
