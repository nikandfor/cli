package cli

import (
	"bytes"
	"testing"

	"github.com/nikandfor/assert"
)

func TestCommandRunSimple(t *testing.T) {
	var ok bool
	var buf bytes.Buffer

	c := &Command{
		Name:        "long,l",
		Description: "test command",
		Args:        Args{},
		Action:      func(*Command) error { ok = true; return nil },
		Help: `Some long descriptive help message here.
Possible multiline.
    With paddings.`,
		Commands: []*Command{},
		Flags: []*Flag{
			NewFlag("flag,f,ff", false, "some flag"),
			NewFlag("flag2", "str", "some flag"),
		},

		Stderr: &buf,
	}
	assert.NotNil(t, c.Args) // require

	err := c.run([]string{"base", "first", "second", "--flag", "-"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, Args{"first", "second", "-"}, c.Args)

	assert.Equal(t, true, c.Flag("flag").Value)

	assert.True(t, ok, "expected command not called")

	assert.Equal(t, ``, buf.String())

	//

	err = c.run([]string{"base", "first", "second", "--flag2"}, nil)
	assert.ErrorIs(t, err, ErrFlagValueRequired)
}

func TestCommandRunSub(t *testing.T) {
	var ok bool
	var buf bytes.Buffer

	c := &Command{
		Name:        "long,l",
		Description: "test command",
		Action: func(*Command) error {
			assert.Fail(t, "command called")
			return nil
		},
		Help: `Some long descriptive help message here.
Possible multiline.
    With paddings.`,
		Commands: []*Command{{
			Name:        "sub,s,alias",
			Description: "subcommand",
			Args:        Args{},
			Action:      func(*Command) error { ok = true; return nil },
			Flags: []*Flag{
				NewFlag("subflag", 3, "some sub flag"),
			},
		}},
		Flags: []*Flag{
			NewFlag("flag,f,ff", "empty", "some flag"),
			HelpFlag,
		},
		Stderr: &buf,
	}

	sub := c.Command("sub")

	assert.NotNil(t, sub.Args) // require

	err := c.run([]string{"base", "sub", "first", "second", "--flag=value", "-", "--subflag", "4"}, nil)

	assert.NoError(t, err)

	assert.Equal(t, sub.Args, Args{"first", "second", "-"})

	assert.Equal(t, "value", sub.Flag("flag").Value)
	assert.Equal(t, 4, sub.Flag("subflag").Value)

	assert.True(t, ok, "expected command not called")

	assert.Equal(t, ``, buf.String())
}

func TestCommandRunSub2(t *testing.T) {
	var buf bytes.Buffer

	c := &Command{
		Name:        "long,l",
		Description: "test command",
		Action: func(*Command) error {
			assert.Fail(t, "command called")
			return nil
		},
		Help: `Some long descriptive help message here.
Possible multiline.
    With paddings.`,
		Commands: []*Command{{
			Name:        "sub,s,alias",
			Description: "subcommand",
			Args:        Args{},
			Action: func(*Command) error {
				assert.Fail(t, "command called")
				return nil
			},
			Flags: []*Flag{
				NewFlag("subflag", 3, "some sub flag"),
			},
		}},
		Flags: []*Flag{
			NewFlag("flag,f,ff", "empty", "some flag"),
			HelpFlag,
		},
		Stderr: &buf,
	}

	sub := c.Command("sub")

	assert.NotNil(t, sub.Args) // require

	err := c.run([]string{"base", "sub", "first", "second", "--flag=value", "-", "--subflag", "4", "--nonexisted"}, nil)

	assert.ErrorIs(t, err, ErrNoSuchFlag)

	assert.Equal(t, c.Command("sub").Args, Args{"first", "second", "-"})

	assert.Equal(t, "value", sub.Flag("flag").Value)
	assert.Equal(t, 4, sub.Flag("subflag").Value)

	assert.Equal(t, ``, buf.String())
}

func TestDoubleDash(t *testing.T) {
	var ok bool
	var buf bytes.Buffer

	c := &Command{
		Name:        "long,l",
		Description: "test command",
		Args:        Args{},
		Action:      func(*Command) error { ok = true; return nil },
		Help: `Some long descriptive help message here.
Possible multiline.
    With paddings.`,
		//	Commands: []*Command{},
		Flags: []*Flag{
			NewFlag("flag,f,ff", 0, "some flag"),
		},

		Stderr: &buf,
	}
	assert.NotNil(t, c.Args) // require

	err := c.run([]string{"base", "first", "--", "a", "--flag"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, Args{"first", "a", "--flag"}, c.Args)

	assert.True(t, ok, "expected command not called")

	assert.Equal(t, ``, buf.String())
}
