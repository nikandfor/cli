package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandRunSimple(t *testing.T) {
	c := &Command{
		Name:        "long,l",
		Description: "test command",
		Action:      nil,
		HelpText: `Some long descriptive help message here.
Possible multiline.
    With paddings.`,
		Commands: []*Command{{
			Name:        "sub,s,alias",
			Description: "subcommand",
			Action:      func(*Command) error { return nil },
		}},
		Flags: []*Flag{
			NewFlag("flag,f,ff", false, "some flag"),
		},
	}

	err := c.run([]string{"base", "first", "second", "--flag", "-"})
	assert.NoError(t, err)
	assert.Equal(t, c.Args, Args{"first", "second", "-"})
}
