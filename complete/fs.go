package complete

import (
	"embed"
	"fmt"
	"io"
	"strings"
)

//go:embed complete.bash complete.zsh
var scripts embed.FS

func Template(shell string) ([]byte, error) {
	return scripts.ReadFile("complete." + shell)
}

func ExecTemplate(w io.Writer, shell string, args []string) (err error) {
	f, err := Template(shell)
	if err != nil {
		return fmt.Errorf("read script: %w", err)
	}

	_, err = fmt.Fprintf(w, string(f), args[0], strings.Join(args, " "))

	return
}
