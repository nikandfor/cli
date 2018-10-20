package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type (
	Flag interface {
		Base() *F
		IsSet() bool
		String() string
		Int() int
		Bool() bool

		parse(n, v string, as Args) (Args, error)
		complete(string) error
	}

	F struct {
		Name    string
		Aliases []string
		Hidden  bool

		compl func(f Flag) func(string) error

		isSet bool
	}

	StringFlag struct {
		F
		Value string
	}

	BoolFlag struct {
		F
		Value bool
	}

	IntFlag struct {
		F
		Value int
	}
)

func (f F) StringFlag(v string) *StringFlag {
	return &StringFlag{
		F:     f,
		Value: v,
	}
}
func (f F) BoolFlag(v bool) *BoolFlag {
	return &BoolFlag{
		F:     f,
		Value: v,
	}
}
func (f F) IntFlag(v int) *IntFlag {
	return &IntFlag{
		F:     f,
		Value: v,
	}
}

func (f *StringFlag) String() string { return f.Value }

func (f *BoolFlag) Bool() bool     { return f.Value }
func (f *BoolFlag) String() string { return fmt.Sprintf("%v", f.Value) }

func (f *IntFlag) Int() int       { return f.Value }
func (f *IntFlag) String() string { return fmt.Sprintf("%v", f.Value) }

func (f *F) Base() *F { return f }
func (f *F) parse(n, rest string, args Args) (Args, error) {
	return args, nil
}
func (f *F) complete(l string) error {
	if f.compl == nil {
		return DefaultFlagCompletion(f)(l)
	}
	return f.compl(f)(l)
}

func (f *F) IsSet() bool    { return f.isSet }
func (f *F) String() string { return fmt.Sprintf("{%v %v}", f.Name, f.Aliases) }
func (f *F) Int() int       { panic("wrong flag type") }
func (f *F) Bool() bool     { panic("wrong flag type") }

func (f *StringFlag) parse(n, rest string, args Args) (args_ Args, err_ error) {
	//	log.Printf("parse %10v  [%q %q %q]", f.Name, n, rest, args)
	//	defer func() {
	//		log.Printf("parse %10v  <%q %q>", f.Name, rest_, args_)
	//	}()
	var err error
	switch {
	case rest != "":
		if rest[0] == '=' {
			f.Value = rest[1:]
		} else {
			f.Value = rest
		}
		args = args[1:]
	case len(args) >= 2:
		f.Value = args[1]
		args = args[2:]
	default:
		err = errors.New("expected arg")
	}
	return args, err
}

func (f *IntFlag) parse(n, rest string, args Args) (args_ Args, err_ error) {
	//	log.Printf("parse %10v  [%q %q %q]", f.Name, n, rest, args)
	//	defer func() {
	//		log.Printf("parse %10v  <%q %q>", f.Name, rest_, args_)
	//	}()
	var err error
	switch {
	case rest != "":
		var v string
		if rest[0] == '=' {
			v = rest[1:]
		} else {
			v = rest
		}
		f.Value, err = strconv.Atoi(v)
		args = args[1:]
	case len(args) >= 2:
		f.Value, err = strconv.Atoi(args[1])
		args = args[2:]
	default:
		err = errors.New("expected arg")
	}
	return args, err
}

func (f *BoolFlag) parse(n, rest string, args Args) (args_ Args, err_ error) {
	//	log.Printf("parse %10v  [%q %q %q]", f.Name, n, rest, args)
	//	defer func() {
	//		log.Printf("parse %10v  <%q %q>", f.Name, rest_, args_)
	//	}()
	var err error
	switch {
	case rest != "" && rest[0] == '=':
		switch strings.ToLower(rest[1:]) {
		case "t", "true", "y", "1":
			f.Value = true
		case "f", "false", "n", "0":
			f.Value = false
		default:
			err = errors.New("expected bool")
		}
		args = args[1:]
	default:
		f.Value = true
		args = args[1:]
	}
	return args, err
}
