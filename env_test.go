package cli

import (
	"testing"

	"github.com/nikandfor/assert"
	"nikand.dev/go/cli/flag"
)

func TestEnvfile(t *testing.T) {
	readFile = func(n string) ([]byte, error) {
		assert.Equal(t, ".env", n)

		return []byte(`PREF_F1=1
		PREF_F2 2
		# PREF_F3=a
		PREF_F4 abc def
		NOT_PREF_F3=3`), nil
	}

	var ok bool

	c := &Command{
		Name:   "testcmd",
		Args:   Args{},
		Action: func(c *Command) error { ok = true; return nil },
		Flags: []*Flag{
			flag.New("f1", "", ""),
			flag.New("f2", 1, ""),
			flag.New("f3", "def", ""),
			flag.New("f4", "", ""),
			flag.New("f5", "", ""),
			flag.New("f6", 0, ""),
			EnvfileFlag,
		},
		EnvPrefix: "PREF_",
		ParseEnv: func(c *Command, env []string) ([]string, error) {
			rest, err := DefaultParseEnv(c, env)
			assert.NoError(t, err)
			if err != nil {
				return nil, err
			}

			return rest, nil
		},
	}

	err := Run(c, []string{"cmd", "--envfile=.env", "a", "-f6=9"}, []string{"PREF_F5=4"})

	assert.NoError(t, err)
	assert.Equal(t, Args{"a"}, c.Args)
	assert.True(t, ok, "command not called")

	assert.Equal(t, "1", c.Flag("f1").Value)
	assert.Equal(t, 2, c.Flag("f2").Value)
	assert.Equal(t, "def", c.Flag("f3").Value)
	assert.Equal(t, "abc def", c.Flag("f4").Value)
	assert.Equal(t, "4", c.Flag("f5").Value)
	assert.Equal(t, 9, c.Flag("f6").Value)

	assert.Equal(t, []string{"NOT_PREF_F3=3"}, c.Env)
}
