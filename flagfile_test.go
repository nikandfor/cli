package cli

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagFile(t *testing.T) {
	var ff bytes.Buffer
	fopen = func(n string) (io.ReadCloser, error) {
		return ioutil.NopCloser(&ff), nil
	}
	ff.Write([]byte(`--flag ffval`))

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
			FlagFileFlag,
		},
	}

	err := c.run([]string{"first", "second", "--flag", "before", "--flagfile=qwe"})
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "first", c.Flag("flag").Value.(string))
	assert.Equal(t, "ffval", c.Commands[0].Flag("flag").Value.(string))
}

func TestFlagFile2(t *testing.T) {
	var ff bytes.Buffer
	fopen = func(n string) (io.ReadCloser, error) {
		return ioutil.NopCloser(&ff), nil
	}
	ff.Write([]byte(`--flag ffval`))

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
			FlagFileFlag,
		},
	}

	err := c.run([]string{"first", "second", "--flag", "before", "--flagfile=qwe", "--flag=after"})
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "first", c.Flag("flag").Value.(string))
	assert.Equal(t, "after", c.Commands[0].Flag("flag").Value.(string))
}
