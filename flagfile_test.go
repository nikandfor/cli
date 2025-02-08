package cli

import (
	"testing"

	"github.com/nikandfor/assert"
	"nikand.dev/go/cli/flag"
)

func TestFlagfile(t *testing.T) {
	readFile = func(n string) ([]byte, error) {
		assert.Equal(t, "file.flagfile", n)

		return []byte(`first second
		third --fourth`), nil
	}

	f := flag.New("ff", "", "")
	args, err := flagfile(f, "--ff=file.flagfile", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"first", "second", "third", "--fourth", "a", "b"}, args)
}

func TestFlagfile2(t *testing.T) {
	readFile = func(n string) ([]byte, error) {
		assert.Equal(t, "file.flagfile", n)

		return []byte(`-a 9
		-b 10 # 11
		-c 12`), nil
	}

	f := flag.New("ff", "", "")
	args, err := flagfile(f, "--ff", []string{"file.flagfile", "a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"-a", "9", "-b", "10", "-c", "12", "a", "b"}, args)
}

func TestFlagfileDecodeArg(t *testing.T) {
	var buf []byte
	var err error

	for _, tc := range []struct {
		s string
		r string
	}{
		{`a\ b\ c`, `a b c`},
		{`"s t r"`, `s t r`},
		{`"s \" t ' r"`, `s " t ' r`},
		{`'q w e'`, `q w e`},
		{`'q " w'`, `q " w`},
		{`ab' 1 'cd" 2 "ed`, `ab 1 cd 2 ed`},
	} {
		buf, _, err = decodeArg([]byte(tc.s), 0, buf[:0])
		assert.NoError(t, err)
		assert.Equal(t, tc.r, string(buf))
	}
}

func TestFlagfileQuotes(t *testing.T) {
	readFile = func(n string) ([]byte, error) {
		assert.Equal(t, "file.flagfile", n)

		return []byte(`-a "s t r" -a1="q \" w ' e"
		-b 'q w e' -b1='q " w'
		-c ab' 1 'cd" 2 "ef
		-d "#not end"
		# end`), nil
	}

	f := flag.New("ff", "", "")
	args, err := flagfile(f, "--ff", []string{"file.flagfile", "a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, []string{`-a`, `s t r`, `-a1=q " w ' e`, `-b`, `q w e`, `-b1=q " w`, `-c`, `ab 1 cd 2 ef`, `-d`, `#not end`, "a", "b"}, args)
}

func TestRunFlagfile(t *testing.T) {
	readFile = func(n string) ([]byte, error) {
		assert.Equal(t, "qwe", n)

		return []byte(`--flag ffval`), nil
	}

	var ok bool

	c := &Command{
		Name: "first",
		Commands: []*Command{{
			Name:   "second",
			Action: func(c *Command) error { ok = true; return nil },
			Flags: []*Flag{
				flag.New("flag", "second", ""),
			},
		}},
		Flags: []*Flag{
			flag.New("flag", "first", ""),
			FlagfileFlag,
		},
	}

	err := Run(c, []string{"first", "second", "--flag", "before", "--flagfile=qwe"}, nil)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "first", c.Flag("flag").Value)
	assert.Equal(t, "ffval", c.Commands[0].Flag("flag").Value)
}

func TestRunFlagfile2(t *testing.T) {
	readFile = func(n string) ([]byte, error) {
		return []byte(`--flag ffval`), nil
	}

	var ok bool

	c := &Command{
		Name: "first",
		Commands: []*Command{{
			Name:   "second",
			Action: func(c *Command) error { ok = true; return nil },
			Flags: []*Flag{
				flag.New("flag", "second", ""),
			},
		}},
		Flags: []*Flag{
			flag.New("flag", "first", ""),
			FlagfileFlag,
		},
	}

	err := Run(c, []string{"first", "second", "--flag", "before", "--flagfile=qwe", "--flag=after"}, nil)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "first", c.Flag("flag").Value)
	assert.Equal(t, "after", c.Commands[0].Flag("flag").Value)
}
