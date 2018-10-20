package app

import (
	"bytes"
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

	err := c.Run([]string{"do", "--_comp-bash", "--_comp-line", "do --flag a", "--_comp-word=2", "--flag", "a"})
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

	err := c.Run([]string{"do", "--_comp-bash", "--_comp-line", "do --flag ", "--_comp-word=2", "--flag"})
	assert.NoError(t, err)
	assert.Equal(t, Args{""}, c.Args())
	assert.Equal(t, true, c.Flags[0].VBool())
	assert.False(t, act)
	assert.True(t, comp)
}

func TestDefaultCommandCompletion(t *testing.T) {
	c := &Command{
		Name: "root",
		Commands: []*Command{
			{Name: "cmda", Aliases: []string{"ca"}},
			{Name: "cmdb", Aliases: []string{"cb"}},
			{Name: "bmd", Aliases: []string{"cd"}},
		},
		Flags: []Flag{
			&F{Name: "flag", Aliases: []string{"f"}},
			&F{Name: "fmag", Aliases: []string{"m"}},
		},
	}

	var buf bytes.Buffer
	Writer = &buf

	err := DefaultCommandComplete(c)
	assert.NoError(t, err)
	assert.Equal(t, `compgen -W "cmda cmdb bmd --flag --fmag"`, buf.String())

	buf.Reset()
	c.args = Args{"c"}
	err = DefaultCommandComplete(c)
	assert.NoError(t, err)
	assert.Equal(t, `compgen -W "cmda cmdb cd"`, buf.String())

	buf.Reset()
	c.args = Args{"-"}
	err = DefaultCommandComplete(c)
	assert.NoError(t, err)
	assert.Equal(t, `compgen -W "--flag --fmag"`, buf.String())

	buf.Reset()
	c.args = Args{"-m"}
	err = DefaultCommandComplete(c)
	assert.NoError(t, err)
	assert.Equal(t, `compgen -W "-m"`, buf.String())

	c.Commands = nil

	buf.Reset()
	c.args = Args{""}
	err = DefaultCommandComplete(c)
	assert.NoError(t, err)
	assert.Equal(t, `compgen -o default ""`, buf.String())

	buf.Reset()
	c.args = Args{"a"}
	err = DefaultCommandComplete(c)
	assert.NoError(t, err)
	assert.Equal(t, `compgen -o default "a"`, buf.String())
}

func TestCompleteFlag(t *testing.T) {
	act := false
	comp := false
	fcomp := false

	c := &Command{
		Name:     "command",
		Aliases:  []string{"com", "do"},
		Action:   func(c *Command) error { act = true; return nil },
		Complete: func(c *Command) error { comp = true; return nil },
		Flags: []Flag{
			F{Name: "flag"}.NewBool(false),
			F{Name: "int", Aliases: []string{"i"}}.NewInt(0),
			F{Name: "val", Aliases: []string{"v"},
				Complete: func(f Flag, _ *Command, last string) error {
					fcomp = true
					return nil
				},
			}.NewString("str"),
		},
	}
	AddCompletionToApp(c)

	err := c.Run([]string{"do", "--_comp-bash", "--_comp-line", "do --flag --val ", "--_comp-word=3", "--flag", "--val"})
	assert.NoError(t, err)
	assert.Equal(t, Args(nil), c.Args())
	assert.Equal(t, true, c.Flags[0].VBool())
	assert.False(t, act)
	assert.False(t, comp)
	assert.True(t, fcomp)
}
