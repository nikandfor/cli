package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	act := false

	c := &Command{
		Name:    "command",
		Aliases: []string{"com", "do"},
		Action:  func(c *Command) error { act = true; return nil },
		Flags: []Flag{
			F{Name: "flag"}.NewBool(false),
			F{Name: "val", Aliases: []string{"v"}}.NewString("str"),
			F{Name: "int", Aliases: []string{"i"}}.NewInt(0),
		},
	}

	err := c.Run([]string{"do", "--flag", "-v", "string", "--int=3", "arg"})
	assert.NoError(t, err)
	assert.Equal(t, Args{"arg"}, c.Args())
	assert.Equal(t, true, c.Flags[0].VBool())
	assert.Equal(t, "string", c.Flags[1].VString())
	assert.Equal(t, 3, c.Flags[2].VInt())
	assert.True(t, act)
}

func TestCommandSub(t *testing.T) {
	fok := false
	fnotok := false

	actok := func(c *Command) error { fok = true; return nil }
	actnotok := func(c *Command) error { fnotok = false; return nil }

	sub := &Command{
		Name:    "sub",
		Aliases: []string{"s"},
		Action:  actok,
	}
	c := &Command{
		Name:     "command",
		Aliases:  []string{"com", "do"},
		Action:   actnotok,
		Commands: []*Command{sub},
		Flags: []Flag{
			F{Name: "flag"}.NewBool(false),
			F{Name: "val"}.NewString("str"),
		},
	}

	err := c.Run([]string{"do", "s", "--flag", "--val", "string", "arg"})
	assert.NoError(t, err)
	assert.Equal(t, Args{"arg"}, sub.Args(), "root args: %v", c.Args())
	assert.Equal(t, true, c.Flags[0].VBool())
	assert.Equal(t, "string", c.Flags[1].VString())
	assert.True(t, fok, "action")
	assert.False(t, fnotok, "root action")
}

func TestCommandParse(t *testing.T) {
	act := false

	c := &Command{
		Name:    "command",
		Aliases: []string{"com", "do"},
		Action:  func(c *Command) error { act = true; return nil },
		Flags: []Flag{
			F{Name: "flag"}.NewBool(false),
			F{Name: "val"}.NewString("str"),
		},
	}

	err := c.Run([]string{"do", "--flag", "--val", "--", "--", "arg"})
	assert.NoError(t, err)
	assert.Equal(t, Args{"arg"}, c.Args())
	assert.Equal(t, true, c.Flags[0].VBool())
	assert.Equal(t, "--", c.Flags[1].VString())
	assert.True(t, act)
}
