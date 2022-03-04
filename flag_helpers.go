package cli

import "time"

func (f *Flag) Bool() bool {
	return f.Value.(bool)
}

func (f *Flag) Duration() time.Duration {
	return f.Value.(time.Duration)
}

func (f *Flag) Int() int {
	return f.Value.(int)
}

func (f *Flag) String() string {
	return f.Value.(string)
}

func (c *Command) Bool(n string) bool {
	return c.mustflag(n).Bool()
}

func (c *Command) Duration(n string) time.Duration {
	return c.mustflag(n).Duration()
}

func (c *Command) Int(n string) int {
	return c.mustflag(n).Int()
}

func (c *Command) String(n string) string {
	return c.mustflag(n).String()
}

func (c *Command) mustflag(n string) (f *Flag) {
	f = c.Flag(n)
	if f == nil {
		panic("no such flag: " + n)
	}

	return
}
