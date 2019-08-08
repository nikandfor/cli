package flag

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	IntFlag struct {
		Name  string
		Value int
		Help  string
	}

	StringFlag struct {
		Name  string
		Value string
		Help  string
	}

	DurationFlag struct {
		Name  string
		Value time.Duration
		Help  string
	}

	BoolFlag struct {
		Name  string
		Value bool
		Help  string
	}
)

func (f *IntFlag) Names() string { return f.Name }

func (f *StringFlag) Names() string { return f.Name }

func (f *DurationFlag) Names() string { return f.Name }

func (f *BoolFlag) Names() string { return f.Name }

func (f *IntFlag) Parse(n, v string, more []string) (rest []string, err error) {
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

	return more, nil
}

func (f *StringFlag) Parse(n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, fmt.Errorf("value expected")
	}

	f.Value = v

	return more, nil
}

func (f *DurationFlag) Parse(n, v string, more []string) (rest []string, err error) {
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

	return more, nil
}

func (f *BoolFlag) Parse(n, v string, more []string) (rest []string, err error) {
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

	return more, nil
}
