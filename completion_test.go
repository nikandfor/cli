package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompletionPref(t *testing.T) {
	act := false
	comp := false

	c := &Command{
		Name:     "command",
		Aliases:  []string{"com", "do"},
		Action:   func(c *Command) error { act = true; return nil },
		Complete: func(c *Command) error { comp = true; return nil },
		Flags: []Flag{
			F{Name: "flag"}.NewBool(false),
			F{Name: "val", Aliases: []string{"v"}}.NewString("str"),
			F{Name: "int", Aliases: []string{"i"}}.NewInt(0),
		},
	}
	AddCompletionToApp(c)

	err := c.Run([]string{"do", "--_comp-bash", "--_comp-line", "do --flag a", "--_comp-word=3", "--flag", "a"})
	assert.NoError(t, err)
	assert.Equal(t, Args{"a"}, c.Args())
	assert.Equal(t, true, c.Flags[0].VBool())
	assert.False(t, act)
	assert.True(t, comp)
}

func TestCompletionNew(t *testing.T) {
	act := false
	comp := false

	c := &Command{
		Name:     "command",
		Aliases:  []string{"com", "do"},
		Action:   func(c *Command) error { act = true; return nil },
		Complete: func(c *Command) error { comp = true; return nil },
		Flags: []Flag{
			F{Name: "flag"}.NewBool(false),
			F{Name: "val", Aliases: []string{"v"}}.NewString("str"),
			F{Name: "int", Aliases: []string{"i"}}.NewInt(0),
		},
	}
	AddCompletionToApp(c)

	err := c.Run([]string{"do", "--_comp-bash", "--_comp-line", "do --flag ", "--_comp-word=3", "--flag"})
	assert.NoError(t, err)
	assert.Equal(t, Args{""}, c.Args())
	assert.Equal(t, true, c.Flags[0].VBool())
	assert.False(t, act)
	assert.True(t, comp)
}
