package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	option func(f *F)

	F struct {
		Name string
		Help string

		set bool
	}

	Bool struct {
		F
		Value bool
	}

	Int struct {
		F
		Value int
	}

	String struct {
		F
		Value string
	}

	Duration struct {
		F
		Value time.Duration
	}
)

func (f *F) setOpts(ops ...option) {
	for _, o := range ops {
		o(f)
	}
}

func (f *F) Base() *F    { return f }
func (f *F) IsSet() bool { return f.set }

func (f *F) check() error {

	return nil
}

func (f *Bool) Names() string     { return f.Name }
func (f *Int) Names() string      { return f.Name }
func (f *String) Names() string   { return f.Name }
func (f *Duration) Names() string { return f.Name }

func NewBool(n string, v bool, h string, opts ...option) *Bool {
	f := F{Name: n, Help: h}
	f.setOpts(opts...)
	return &Bool{F: f, Value: v}
}
func NewInt(n string, v int, h string, opts ...option) *Int {
	f := F{Name: n, Help: h}
	f.setOpts(opts...)
	return &Int{F: f, Value: v}
}
func NewString(n, v, h string, opts ...option) *String {
	f := F{Name: n, Help: h}
	f.setOpts(opts...)
	return &String{F: f, Value: v}
}
func NewDuration(n string, v time.Duration, h string, opts ...option) *Duration {
	f := F{Name: n, Help: h}
	f.setOpts(opts...)
	return &Duration{F: f, Value: v}
}

func (f *Bool) Parse(n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	}

	switch strings.ToLower(v) {
	case "", "true", "yes", "yeah", "1":
		f.Value = true
	case "false", "no", "nope", "0":
		f.Value = false
	default:
		return nil, fmt.Errorf("can't parse bool value: %v", v)
	}

	f.set = true

	return more, nil
}

func (f *Int) Parse(n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, fmt.Errorf("value expected")
	}

	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	f.Value = int(val)

	f.set = true

	return more, nil
}

func (f *String) Parse(n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, fmt.Errorf("value expected")
	}

	f.Value = v

	f.set = true

	return more, nil
}

func (f *Duration) Parse(n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, fmt.Errorf("value expected")
	}

	f.Value, err = time.ParseDuration(v)
	if err != nil {
		return nil, err
	}

	f.set = true

	return more, nil
}
