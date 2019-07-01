package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type (
	FlagAction func(f Flag, c *Command) error

	Flag interface {
		Base() *F

		Type() string
		IsSet() bool

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
		Name           string
		Aliases        []string
		Hidden         bool
		Before         FlagAction
		After          FlagAction
		Description    string
		Completion     func(f Flag, c *Command, last string) error
		CompletionHelp string

		isSet bool
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

	LevelFlag struct {
		IntFlag
	}

	StringSliceFlag struct {
		F
		Value []string
	}

	EnumFlag struct {
		F
		Value   string
		Options []string
	}
)

func (f F) NewInt(v int) *IntFlag { return &IntFlag{F: f, Value: v} }
func (f F) NewBool(v bool) *BoolFlag {
	return &BoolFlag{F: f, Value: v}
}
func (f F) NewString(v string) *StringFlag {
	return &StringFlag{F: f, Value: v}
}
func (f F) NewFile(v string) *FileFlag {
	return &FileFlag{StringFlag{F: f, Value: v}}
}
func (f F) NewLevel(v int) *LevelFlag                  { return &LevelFlag{IntFlag{F: f, Value: v}} }
func (f F) NewStringSlice(v []string) *StringSliceFlag { return &StringSliceFlag{F: f, Value: v} }
func (f F) NewEnum(v string, opts []string) *EnumFlag  { return &EnumFlag{F: f, Value: v, Options: opts} }

func (f F) IsSet() bool { return f.isSet }

func (f *F) Base() *F               { return f }
func (f *IntFlag) Base() *F         { return &f.F }
func (f *BoolFlag) Base() *F        { return &f.F }
func (f *StringFlag) Base() *F      { return &f.F }
func (f *StringSliceFlag) Base() *F { return &f.F }
func (f *EnumFlag) Base() *F        { return &f.F }

func (f *F) Type() string               { return "" }
func (f *IntFlag) Type() string         { return "int" }
func (f *BoolFlag) Type() string        { return "bool" }
func (f *StringFlag) Type() string      { return "string" }
func (f *FileFlag) Type() string        { return "file" }
func (f *LevelFlag) Type() string       { return "level" }
func (f *StringSliceFlag) Type() string { return "string slice" }
func (f *EnumFlag) Type() string        { return "enum" }

func (f *F) Parse(name, val string, rep bool) (bool, error) { return false, nil }
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

	f.isSet = true

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
	f.isSet = true
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
	f.isSet = true
	return false, nil
}
func (f *LevelFlag) Parse(name, val string, rep bool) (bool, error) {
	if val == "" {
		f.Value = 1
		return false, nil
	}
	return f.IntFlag.Parse(name, val, rep)
}
func (f *StringSliceFlag) Parse(name, val string, rep bool) (bool, error) {
	if val == "" && !rep {
		return true, nil
	}
	if !rep && val[0] == '=' {
		val = val[1:]
	}
	f.Value = append(f.Value, val)
	f.isSet = true
	return false, nil
}
func (f *EnumFlag) Parse(name, val string, rep bool) (bool, error) {
	if val == "" {
		return true, nil
	}
	if !rep && val[0] == '=' {
		val = val[1:]
	}

	for _, o := range f.Options {
		if val == o {
			f.Value = val
			f.isSet = true
			return false, nil
		}
	}

	return false, fmt.Errorf("got %v exptected one of %v", val, f.Options)
}

func (f F) VInt() int          { panic("wrong type") }
func (f F) VBool() bool        { panic("wrong type") }
func (f F) VString() string    { panic("wrong type") }
func (f F) Value() interface{} { panic("wrong type") }

func (f F) String() string { return fmt.Sprintf("{%v %v}", f.Name, f.Aliases) }

func (f *IntFlag) VInt() int          { return f.Value }
func (f *BoolFlag) VBool() bool       { return f.Value }
func (f *StringFlag) VString() string { return f.Value }
func (f *EnumFlag) VString() string   { return f.Value }

func (f *F) VAny() interface{}               { return nil }
func (f *IntFlag) VAny() interface{}         { return f.Value }
func (f *BoolFlag) VAny() interface{}        { return f.Value }
func (f *StringFlag) VAny() interface{}      { return f.Value }
func (f *StringSliceFlag) VAny() interface{} { return f.Value }
func (f *EnumFlag) VAny() interface{}        { return f.Options[0] }
