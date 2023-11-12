package complete

import (
	"fmt"
	"path/filepath"
	"strconv"
)

type (
	Getenver interface {
		Getenv(string) string
	}

	LookupEnver interface {
		LookupEnv(string) (string, bool)
	}
)

func Shell(env LookupEnver) (sh string, ok bool) {
	_, ok = env.LookupEnv("BASH")
	if ok {
		//	return filepath.Base(sh), true
		return "bash", true
	}

	_, ok = env.LookupEnv("ZSH_NAME")
	if ok {
		return "zsh", true
	}

	_, ok = env.LookupEnv("SHELL")
	if ok {
		return filepath.Base(sh), true
	}

	return "", false
}

func Current(env LookupEnver) (cur string) {
	cur, _ = env.LookupEnv("CLI_COMP_CUR")
	return
}

func Prev(env LookupEnver) (prev string) {
	prev, _ = env.LookupEnv("CLI_COMP_PREV")
	return
}

func Args(env LookupEnver) (args []string, cur int) {
	x, _ := env.LookupEnv("CLI_COMP_WORDS_INDEX")

	l, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		return
	}

	cur = int(l)

	x, _ = env.LookupEnv("CLI_COMP_WORDS_LENGTH")

	l, err = strconv.ParseInt(x, 10, 32)
	if err != nil {
		return
	}

	args = make([]string, l)

	for i := 0; i < int(l); i++ {
		x, _ = env.LookupEnv(fmt.Sprintf("CLI_COMP_WORDS_%d", i))

		args[i] = x
	}

	return
}
