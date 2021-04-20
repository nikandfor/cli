package cli

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

var HelpFlag = &Flag{
	Name:        "help,h",
	Description: "show that message",
	After:       defaultHelp,
	Value:       boolptr(false),
}

var funcs = map[string]interface{}{
	"formatCmd": func(c *Command, full bool) string {
		var b strings.Builder

		pad := func(s, w int) {
			if s >= w {
				return
			}

			_, _ = b.WriteString("                                          "[:w-s])
		}

		_, _ = b.WriteString(c.Name)

		if c.Usage != "" {
			_, _ = b.WriteString(" ")
			_, _ = b.WriteString(c.Usage)
		}

		pad(b.Len(), 25)

		if c.Description != "" {
			_, _ = b.WriteString(" - ")
			_, _ = b.WriteString(c.Description)
		}

		if full && c.HelpText != "" {
			_, _ = b.WriteString("\n\n")
			_, _ = b.WriteString(c.HelpText)
		}

		return b.String()
	},
	"formatFlag": func(f *Flag) string {
		var b strings.Builder

		pad := func(s, w int) {
			if s >= w {
				return
			}

			_, _ = b.WriteString("                                          "[:w-s])
		}

		_, _ = b.WriteString(f.Name)

		var val string
		if f.Value != nil {
			r := reflect.ValueOf(f.Value)
			for r.Kind() == reflect.Ptr && !r.IsNil() {
				r = r.Elem()
			}

			if !r.IsZero() {
				val = fmt.Sprintf("=%v", r)
			}
		}

		if val != "" {
			_, _ = b.WriteString(val)
		}

		pad(b.Len(), 25)

		if f.Description != "" {
			_, _ = b.WriteString(" - ")
			_, _ = b.WriteString(f.Description)
		}

		return b.String()
	},
}

var commandHelpTemplate = template.Must(template.New("command help").Funcs(funcs).Parse(`{{ formatCmd . true }}
{{ if .Commands }}
Subcommands:
{{- range .Commands }}
    {{ if not .Hidden }}{{ formatCmd . false }}{{ end }}
{{- end }}
{{- end }}
{{- if or .Flags .Parent }}
Flags:
{{- block "flags" . }}
{{- range .Flags }}
    {{ if not .Hidden }}{{ formatFlag . }}{{ end }}
{{- end }}
{{- with .Parent }}{{ template "flags" . }}{{ end }}
{{- end }}
{{- end }}
`))

func defaultHelp(f *Flag, c *Command) error {
	err := commandHelpTemplate.Execute(stdout, c)
	if err != nil {
		return err
	}
	return ErrFlagExit
}

func boolptr(v bool) *bool { return &v }
