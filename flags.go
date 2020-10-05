package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nikandfor/errors"
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

var (
	ErrFlagExit          = errors.New("flag exit")
	ErrFlagValueExpected = errors.New("value expected")
)

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
	case FlagValue:
		val = v
	case bool:
		val = &v
	case int:
		val = &v
	case string:
		val = &v
	case time.Duration:
		val = &v
	case []string:
		val = &v
	case *bool, *int, *string, *time.Duration, *[]string:
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

	q := false
	switch strings.ToLower(v) {
	case "", "true", "yes", "yeah", "y", "1":
		q = true
	case "false", "no", "nope", "n", "0":
	default:
		return nil, fmt.Errorf("can't parse bool value: %v", v)
	}

	*f.Value.(*bool) = q

	return more, nil
}

func parseInt(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, newValExp(n)
	}

	vl, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	*f.Value.(*int) = int(vl)

	return more, nil
}

func parseString(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, newValExp(n)
	}

	*f.Value.(*string) = v

	return more, nil
}

func parseDuration(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, newValExp(n)
	}

	q, err := time.ParseDuration(v)
	if err != nil {
		return nil, err
	}

	*f.Value.(*time.Duration) = q

	return more, nil
}

func parseStringSlice(f *Flag, n, v string, more []string) (rest []string, err error) {
	if len(v) != 0 && v[0] == '=' {
		v = v[1:]
	} else if v == "" && len(more) != 0 {
		v = more[0]
		more = more[1:]
	} else {
		return nil, newValExp(n)
	}

	list := f.Value.(*[]string)
	*list = append(*list, v)

	return more, nil
}

func newValExp(n string) error {
	return errors.Wrap(ErrFlagValueExpected, "flag %v", n)
}
