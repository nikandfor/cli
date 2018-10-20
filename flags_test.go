package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagsNew(t *testing.T) {
	fbase := F{Name: "name", Aliases: []string{"flag", "n", "f"}}

	var f Flag
	f = fbase.NewInt(4)
	assert.Equal(t, 4, f.VInt())

	f = fbase.NewBool(true)
	assert.Equal(t, true, f.VBool())

	f = fbase.NewString("str_v")
	assert.Equal(t, "str_v", f.VString())

	assert.Equal(t, "{name [flag n f]}", f.(*StringFlag).String())
}

func TestFlagsParseInt(t *testing.T) {
	fbase := F{Name: "name", Aliases: []string{"flag", "n", "f"}}

	var f FlagDev = fbase.NewInt(4)

	more, err := f.Parse("flag", "=3", false)
	assert.NoError(t, err)
	assert.False(t, more)
	assert.Equal(t, 3, f.VAny())

	more, err = f.Parse("flag", "", false)
	assert.NoError(t, err)
	assert.True(t, more)
	more, err = f.Parse("flag", "5", true)
	assert.NoError(t, err)
	assert.False(t, more)
	assert.Equal(t, 5, f.VAny())
}

func TestFlagsParseString(t *testing.T) {
	fbase := F{Name: "name", Aliases: []string{"flag", "n", "f"}}

	var f FlagDev = fbase.NewString("4")

	more, err := f.Parse("flag", "=3", false)
	assert.NoError(t, err)
	assert.False(t, more)
	assert.Equal(t, "3", f.VAny())

	more, err = f.Parse("flag", "", false)
	assert.NoError(t, err)
	assert.True(t, more)
	more, err = f.Parse("flag", "5", true)
	assert.NoError(t, err)
	assert.False(t, more)
	assert.Equal(t, "5", f.VAny())
}

func TestFlagsParseBool(t *testing.T) {
	fbase := F{Name: "name", Aliases: []string{"flag", "n", "f"}}

	var f FlagDev = fbase.NewBool(true)

	more, err := f.Parse("flag", "=n", false)
	assert.NoError(t, err)
	assert.False(t, more)
	assert.Equal(t, false, f.VAny())

	more, err = f.Parse("flag", "", false)
	assert.NoError(t, err)
	assert.False(t, more)
	assert.Equal(t, true, f.VAny())
}
