package cli

import (
	"bufio"
	"os"
	"strings"
)

var EnvfileFlag *Flag

func init() {
	EnvfileFlag = &Flag{
		Name:        "envfile",
		Description: "load env variables from file",
		After:       EnvfileFlagAction,
		Value:       StringPtr(""),
	}
}

func EnvfileFlagAction(ff *Flag, c *Command) error {
	//	tlog.Printf("EnvfileFlagAction: %v\n", *ff.Value.(*string))

	f, err := fopen(*ff.Value.(*string))
	if err != nil {
		return err
	}
	defer func() {
		if e := f.Close(); err == nil {
			err = e
		}
	}()

	r := bufio.NewScanner(f)
	r.Split(bufio.ScanLines)

	//vars:
	for r.Scan() {
		e := r.Text()
		e = strings.TrimSpace(e)
		if strings.HasPrefix(e, "#") {
			continue
		}

		p := strings.Index(e, "=")
		if p == -1 {
			os.Setenv(e, "")
		} else {
			os.Setenv(e[:p], e[p+1:])
		}
	}

	if err = r.Err(); err != nil {
		return err
	}

	if err = c.parseEnv(true); err != nil {
		return err
	}

	return nil
}

func varname(s string) string {
	p := strings.Index(s, "=")
	if p != -1 {
		s = s[:p]
	}

	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "_", "-")

	return s
}
