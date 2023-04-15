package cli

import "time"

func (c *Command) Bool(n string) bool {
	return c.mustflag(n).Bool()
}

func (c *Command) Duration(n string) time.Duration {
	return c.mustflag(n).Duration()
}

func (c *Command) Float64(n string) float64 {
	return c.mustflag(n).Float64()
}

func (c *Command) Float32(n string) float32 {
	return c.mustflag(n).Float32()
}

func (c *Command) Int(n string) int {
	return c.mustflag(n).Int()
}

func (c *Command) Int64(n string) int64 {
	return c.mustflag(n).Int64()
}

func (c *Command) Uint(n string) uint {
	return c.mustflag(n).Uint()
}

func (c *Command) Uint64(n string) uint64 {
	return c.mustflag(n).Uint64()
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
