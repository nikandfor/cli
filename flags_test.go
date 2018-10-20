package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagsSmoke(t *testing.T) {
	args := Args{"--bool", "-b=false", "--string", "str_val", "--int", "5", "--s=val2", "--i=3"}
	a := args
	var err error

	fb := F{Name: "bool", Aliases: []string{"b"}}.BoolFlag(false)

	a, err = fb.parse(a.First(), "", a)
	assert.NoError(t, err)
	assert.Equal(t, args[1:], a)
	assert.Equal(t, true, fb.Value)

	a, err = fb.parse(a.First(), "=false", a)
	assert.NoError(t, err)
	assert.Equal(t, args[2:], a)
	assert.Equal(t, false, fb.Value)

	fs := F{Name: "string", Aliases: []string{"s"}}.StringFlag("")

	a, err = fs.parse(a.First(), "", a)
	assert.NoError(t, err)
	assert.Equal(t, args[4:], a)
	assert.Equal(t, "str_val", fs.Value)

	fi := F{Name: "int", Aliases: []string{"i"}}.IntFlag(0)

	a, err = fi.parse(a.First(), "", a)
	assert.NoError(t, err)
	assert.Equal(t, args[6:], a)
	assert.Equal(t, 5, fi.Value)

	// --name=val
	a, err = fs.parse(a.First(), "=val2", a)
	assert.NoError(t, err)
	assert.Equal(t, args[7:], a)
	assert.Equal(t, "val2", fs.Value)

	a, err = fi.parse(a.First(), "=3", a)
	assert.NoError(t, err)
	assert.Equal(t, args[8:], a)
	assert.Equal(t, 3, fi.Value)
}
