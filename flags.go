package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	option func(f *Flag)

	FlagAction func(f *Flag, c *Command) error

	Flag struct {
		Name        string
		Description string
		Before      FlagAction
		After       FlagAction

		Value FlagValue

		IsSet bool
	}

	FlagValue interface {
		Parse(f *Flag, name, val string, more []string) (rest []string, err error)
	}

	Bool struct {
		Value bool
	}

	Int struct {
		Value int
	}

	String struct {
		Value string
	}

	Duration struct {
		Value time.Duration
	}
)

var ErrFlagExit = errors.New("flag exit")

func (f *Flag) setOpts(ops ...option) {
	for _, o := range ops {
		o(f)
	}
}

func (f *Flag) check() error {

	return nil
}

func NewFlag(n string, v interface{}, d string, opts ...option) *Flag {
	var val FlagValue

	switch v := v.(type) {
	case FlagValue:
		val = v
	case bool:
		val = &Bool{v}
	case int:
		val = &Int{v}
	case string:
		val = &String{v}
	case time.Duration:
		val = &Duration{v}
	default:
		panic("unsupported flag value type")
	}

	f := &Flag{Name: n, Description: d, Value: val}
	f.setOpts(opts...)
	return f
}

func (fv *Bool) Parse(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	}

	switch strings.ToLower(v) {
	case "", "true", "yes", "yeah", "1":
		fv.Value = true
	case "false", "no", "nope", "0":
		fv.Value = false
	default:
		return nil, fmt.Errorf("can't parse bool value: %v", v)
	}

	return more, nil
}

func (fv *Int) Parse(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, fmt.Errorf("value expected")
	}

	vl, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	fv.Value = int(vl)

	return more, nil
}

func (fv *String) Parse(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, fmt.Errorf("value expected")
	}

	fv.Value = v

	return more, nil
}

func (fv *Duration) Parse(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, fmt.Errorf("value expected")
	}

	fv.Value, err = time.ParseDuration(v)
	if err != nil {
		return nil, err
	}

	return more, nil
}
