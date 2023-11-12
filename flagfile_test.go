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
	args, err := flagfile(f, "--ff", []string{"file.flagfile", "a", "b"})
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

func TestRunFlagfile(t *testing.T) {
	readFile = func(n string) ([]byte, error) {
		assert.Equal(t, "qwe", n)

		return []byte(`--flag ffval`), nil
	}

	var ok bool

	c := Command{
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

	err := c.run([]string{"first", "second", "--flag", "before", "--flagfile=qwe"}, nil)
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

	c := Command{
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

	err := c.run([]string{"first", "second", "--flag", "before", "--flagfile=qwe", "--flag=after"}, nil)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "first", c.Flag("flag").Value)
	assert.Equal(t, "after", c.Commands[0].Flag("flag").Value)
}
