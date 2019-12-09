package cli

import (
	"os"
	"text/template"
)

var HelpFlag = &Flag{
	Name:        "help,h",
	Description: "show that message",
	After:       defaultHelp,
	Value:       false,
}

var commandHelpTemplate = template.Must(template.New("command help").Parse(`{{ .Name }} {{ .Usage }} - {{ .Description }}
{{- if .HelpText }}

{{ .HelpText -}}
{{ end }}
{{ if .Commands }}
Subcommands:
{{ range .Commands }}    {{ .Name }} {{ .Usage }}{{ if .Description }} - {{ .Description }}{{ end }}
{{ end }}
{{- end }}
{{- if .Flags }}
Flags:
{{- range .Flags }}
    {{ .Name }}{{ with .Value }}={{ . }}{{ end }}				- {{ .Description }}
{{- end }}
{{- end }}
`))

func defaultHelp(f *Flag, c *Command) error {
	err := commandHelpTemplate.Execute(os.Stdout, c)
	if err != nil {
		return err
	}
	return ErrFlagExit
}
