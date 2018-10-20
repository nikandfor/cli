package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type (
	FlagAction func(f Flag) error

	Flag interface {
		Base() *F

		VString() string
		VInt() int
		VBool() bool
		VAny() interface{}
	}

	FlagDev interface {
		Flag

		Parse(arg, val string, args []string) ([]string, error)
	}

	F struct {
		Name    string
		Aliases []string
		Hidden  bool
		After   FlagAction
	}

	IntFlag struct {
		F
		Value int
	}

	BoolFlag struct {
		F
		Value bool
	}

	StringFlag struct {
		F
		Value string
	}
)

func (f F) NewInt(v int) *IntFlag { return &IntFlag{F: f, Value: v} }
func (f F) NewBool(v bool) *BoolFlag {
	return &BoolFlag{F: f, Value: v}
}
func (f F) NewString(v string) *StringFlag {
	return &StringFlag{F: f, Value: v}
}

func (f *F) Base() *F          { return f }
func (f *IntFlag) Base() *F    { return &f.F }
func (f *BoolFlag) Base() *F   { return &f.F }
func (f *StringFlag) Base() *F { return &f.F }

func (f *IntFlag) Parse(name, val string, args []string) ([]string, error) {
	if val == "" {
		if len(args) < 2 {
			return nil, errors.New("argument expected")
		}
		val = args[1]
		args = args[1:]
		if val == "" {
			return nil, errors.New("argument int expected")
		}
	} else {
		val = val[1:]
	}

	q, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}
	f.Value = q
	return args[1:], nil
}
func (f *BoolFlag) Parse(name, val string, args []string) ([]string, error) {
	switch {
	case val == "":
		f.Value = true
		return args[1:], nil
	case val[0] == '=':
		val = val[1:] // remove '='
		switch strings.ToLower(val) {
		case "t", "y", "1", "true", "yes":
			f.Value = true
		case "f", "n", "0", "false", "no":
			f.Value = false
		default:
			return nil, errors.New("expected bool value")
		}
		return args[1:], nil
	default:
		f.Value = true
		args[0] = "-" + val
		return args, nil
	}
}
func (f *StringFlag) Parse(name, val string, args []string) (args_ []string, _ error) {
	if val == "" {
		if len(args) < 2 {
			return nil, errors.New("argument expected")
		}
		val = args[1]
		args = args[1:]
	} else {
		val = val[1:] // remove '='
	}

	f.Value = val
	return args[1:], nil
}

func (f F) VInt() int          { panic("wrong type") }
func (f F) VBool() bool        { panic("wrong type") }
func (f F) VString() string    { panic("wrong type") }
func (f F) Value() interface{} { panic("wrong type") }

func (f F) String() string { return fmt.Sprintf("{%v %v}", f.Name, f.Aliases) }

func (f *IntFlag) VInt() int          { return f.Value }
func (f *BoolFlag) VBool() bool       { return f.Value }
func (f *StringFlag) VString() string { return f.Value }

func (f *F) VAny() interface{}          { return nil }
func (f *IntFlag) VAny() interface{}    { return f.Value }
func (f *BoolFlag) VAny() interface{}   { return f.Value }
func (f *StringFlag) VAny() interface{} { return f.Value }
