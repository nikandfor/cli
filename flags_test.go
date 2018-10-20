package app

import (
	"fmt"
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

	args, err := f.Parse("flag", "=3", []string{"--flag=3", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"rest"}, args)
	assert.Equal(t, 3, f.VAny())

	args, err = f.Parse("flag", "", []string{"--flag", "5", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"rest"}, args)
	assert.Equal(t, 5, f.VAny())

	args, err = f.Parse("f", "=3", []string{"-f=3", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"rest"}, args)
	assert.Equal(t, 3, f.VAny())

	args, err = f.Parse("f", "", []string{"-f", "5", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"rest"}, args)
	assert.Equal(t, 5, f.VAny())
}

func TestFlagsParseString(t *testing.T) {
	fbase := F{Name: "name", Aliases: []string{"flag", "n", "f"}}

	var f FlagDev = fbase.NewString("4")

	for _, name := range []string{"f", "name"} {
		for _, eq := range []bool{false, true} {
			for _, val := range []string{"3", "5", "-", "--"} {
				v, vsep := "", ""
				rest := []string{"--" + name}
				if eq {
					v = "=" + val
					rest[0] += v
				} else {
					vsep = val
					rest = append(rest, val)
				}
				rest = append(rest, "rest")
				t.Run(fmt.Sprintf("--%s%s_%s", name, v, vsep), func(t *testing.T) {
					args, err := f.Parse(name, v, rest)
					assert.NoError(t, err)
					assert.Equal(t, []string{"rest"}, args)
					assert.Equal(t, val, f.VAny())
				})
			}
		}
	}
}

func TestFlagsParseBool(t *testing.T) {
	fbase := F{Name: "name", Aliases: []string{"flag", "n", "f"}}

	var f FlagDev = fbase.NewBool(true)

	args, err := f.Parse("flag", "=n", []string{"--flag=n", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"rest"}, args)
	assert.Equal(t, false, f.VAny())

	args, err = f.Parse("name", "", []string{"--name", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"rest"}, args)
	assert.Equal(t, true, f.VAny())

	args, err = f.Parse("f", "=0", []string{"-f=0", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"rest"}, args)
	assert.Equal(t, false, f.VAny())

	args, err = f.Parse("f", "", []string{"-f", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"rest"}, args)
	assert.Equal(t, true, f.VAny())

	f.(*BoolFlag).Value = false

	args, err = f.Parse("f", "abc", []string{"-fabc", "rest"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"-abc", "rest"}, args)
	assert.Equal(t, true, f.VAny())
}
