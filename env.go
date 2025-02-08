package cli

import (
	"bufio"
	"bytes"
	"errors"
	"strings"

	"nikand.dev/go/cli/flag"
)

var EnvfileFlag = &Flag{
	Name:        "envfile",
	Description: "load env variables from file",
	Action:      envfile,
}

func (c *Command) Getenv(key string) (val string) {
	val, _ = c.LookupEnv(key)
	return
}

func (c *Command) LookupEnv(key string) (string, bool) {
	for _, e := range c.Env {
		p := strings.IndexAny(e, "= ")
		if p == -1 {
			p = len(e)
		}

		if key == e[:p] {
			if p < len(e) {
				return e[p+1:], true
			} else {
				return "", true
			}
		}
	}

	return "", false
}

func envfile(f *Flag, arg string, args []string) (_ []string, err error) {
	c := f.CurrentCommand.(*Command)

	_, val, args, err := flag.ParseArg(arg, args, true, false)
	if err != nil {
		return nil, err
	}

	data, err := readFile(val)
	if err != nil {
		return nil, wrap(err, "read file")
	}

	r := bufio.NewScanner(bytes.NewReader(data))
	r.Split(bufio.ScanLines)

	var env []string

	for r.Scan() {
		e := r.Text()

		e = strings.TrimSpace(e)
		if strings.HasPrefix(e, "#") {
			continue
		}

		e = strings.TrimPrefix(e, "export ")
		e = strings.TrimSpace(e)

		env = append(env, e)
	}

	if err = r.Err(); err != nil {
		return nil, wrap(err, "scan file")
	}

	env, err = c.parseEnv(env)
	if err != nil {
		return nil, wrap(err, "parse env")
	}

	c.Env = append(c.Env, env...)

	return args, nil
}

func DefaultParseEnv(c *Command, env []string) (rest []string, err error) {
	prefix := GetEnvPrefix(c)
	if prefix == "" {
		return env, nil
	}

	for i := 0; i < len(env); i++ {
		if !strings.HasPrefix(env[i], prefix) {
			rest = append(rest, env[i])

			continue
		}

		e := strings.TrimPrefix(env[i], prefix)

		p := strings.Index(e, "=")
		if p == -1 {
			e = varname(e)
		} else {
			e = varname(e[:p]) + e[p:]
		}

		_, err = c.parseFlag(e, nil)
		if errors.Is(err, ErrNoSuchFlag) {
			rest = append(rest, env[i])

			continue
		}
		if err != nil {
			return nil, err
		}
	}

	return rest, nil
}

func varname(s string) string {
	s = strings.ToLower(s)

	s = strings.ReplaceAll(s, "_", "-")

	return s
}
