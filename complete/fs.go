package complete

import (
	"embed"
	"fmt"
	"io"
	"strings"

	"github.com/nikandfor/errors"
)

//go:embed complete.bash complete.zsh
var scripts embed.FS

func Template(shell string) ([]byte, error) {
	return scripts.ReadFile("complete." + shell)
}

func ExecTemplate(w io.Writer, shell string, args []string) (err error) {
	f, err := Template(shell)
	if err != nil {
		return errors.Wrap(err, "read script")
	}

	_, err = fmt.Fprintf(w, string(f), args[0], strings.Join(args, " "))

	return
}
