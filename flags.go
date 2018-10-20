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

		Parse(arg, val string, rep bool) (bool, error)
	}

	F struct {
		Name     string
		Aliases  []string
		Hidden   bool
		After    FlagAction
		Complete func(f Flag, last string) error
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

	FileFlag struct {
		StringFlag
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

func (f *IntFlag) Parse(name, val string, rep bool) (bool, error) {
	if val == "" {
		return true, nil
	}
	if !rep && val[0] == '=' {
		val = val[1:]
	}

	q, err := strconv.Atoi(val)
	if err != nil {
		return false, err
	}
	f.Value = q

	return false, nil
}
func (f *BoolFlag) Parse(name, val string, _ bool) (bool, error) {
	if val == "" {
		f.Value = true
		return false, nil
	}
	if val[0] == '=' {
		val = val[1:]
	}
	switch strings.ToLower(val) {
	case "t", "y", "1", "true", "yes":
		f.Value = true
	case "f", "n", "0", "false", "no":
		f.Value = false
	default:
		return false, errors.New("expected bool value")
	}
	return false, nil

}
func (f *StringFlag) Parse(name, val string, rep bool) (bool, error) {
	if val == "" {
		return true, nil
	}
	if !rep && val[0] == '=' {
		val = val[1:]
	}
	f.Value = val
	return false, nil
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
