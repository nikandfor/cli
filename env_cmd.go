package cli

import (
	"fmt"
	"strings"
)

var EnvCmd = &Command{
	Name:   "env",
	Action: envAction,
}

func envAction(c *Command) (err error) {
	var pref []string

	for q := c; q != nil; q = q.Parent {
		if q.EnvPrefix != "" {
			pref = append(pref, q.EnvPrefix)
		}
	}

	if len(pref) == 0 {
		return nil
	}

	for _, e := range c.OSEnv {
		for _, p := range pref {
			if strings.HasPrefix(e, p) {
				fmt.Fprintf(c.Stdout, "%v\n", e)
				break
			}
		}
	}

	return nil
}
