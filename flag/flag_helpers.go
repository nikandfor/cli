package flag

import "time"

func (f *Flag) Bool() bool {
	return f.Value.(bool)
}

func (f *Flag) Duration() time.Duration {
	return f.Value.(time.Duration)
}

func (f *Flag) Float64() float64 {
	return f.Value.(float64)
}

func (f *Flag) Float32() float32 {
	return f.Value.(float32)
}

func (f *Flag) Int() int {
	return f.Value.(int)
}

func (f *Flag) Int64() int64 {
	return f.Value.(int64)
}

func (f *Flag) Uint() uint {
	return f.Value.(uint)
}

func (f *Flag) Uint64() uint64 {
	return f.Value.(uint64)
}

func (f *Flag) String() string {
	return f.Value.(string)
}
