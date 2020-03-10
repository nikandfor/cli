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

		Value interface{}

		Hidden    bool
		Mandatory bool
		IsSet     bool
	}

	FlagValue interface {
		Parse(f *Flag, name, val string, more []string) (rest []string, err error)
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
	var val interface{}

	switch v := v.(type) {
	case FlagValue, bool, int, string, time.Duration, []string,
		*bool, *int, *string, *time.Duration:
		val = v
	default:
		panic("unsupported flag value type")
	}

	f := &Flag{Name: n, Description: d, Value: val}
	f.setOpts(opts...)
	return f
}

func parseBool(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	}

	switch strings.ToLower(v) {
	case "", "true", "yes", "yeah", "y", "1":
		f.Value = true
	case "false", "no", "nope", "n", "0":
		f.Value = false
	default:
		return nil, fmt.Errorf("can't parse bool value: %v", v)
	}

	return more, nil
}

func parseInt(f *Flag, n, v string, more []string) (rest []string, err error) {
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

	f.Value = int(vl)

	return more, nil
}

func parseString(f *Flag, n, v string, more []string) (rest []string, err error) {
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

func parseDuration(f *Flag, n, v string, more []string) (rest []string, err error) {
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

func parseStringSlice(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, fmt.Errorf("value expected")
	}

	f.Value = append(f.Value.([]string), v)

	return more, nil
}
