package cli

import (
	"testing"

	"github.com/nikandfor/assert"
)

func TestFlagFile(t *testing.T) {
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
				NewFlag("flag", "second", ""),
			},
		}},
		Flags: []*Flag{
			NewFlag("flag", "first", ""),
			FlagfileFlag,
		},
	}

	err := c.run([]string{"first", "second", "--flag", "before", "--flagfile=qwe"}, nil)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "first", c.Flag("flag").Value)
	assert.Equal(t, "ffval", c.Commands[0].Flag("flag").Value)
}

func TestFlagFile2(t *testing.T) {
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
				NewFlag("flag", "second", ""),
			},
		}},
		Flags: []*Flag{
			NewFlag("flag", "first", ""),
			FlagfileFlag,
		},
	}

	err := c.run([]string{"first", "second", "--flag", "before", "--flagfile=qwe", "--flag=after"}, nil)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "first", c.Flag("flag").Value)
	assert.Equal(t, "after", c.Commands[0].Flag("flag").Value)
}
