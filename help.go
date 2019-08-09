package cli

import (
	"os"
	"text/template"
)

var HelpFlag = &Bool{F: F{
	Name:        "help,h",
	Description: "show that message",
	After:       defaultHelp,
}}

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
    {{ .Base.Name }} - {{ .Description }}{{ if .Value }} ({{ .Value }}){{ end }}
{{- end }}
{{- end }}
`))

func defaultHelp(f Flag, c *Command) error {
	err := commandHelpTemplate.Execute(os.Stdout, c)
	if err != nil {
		return err
	}
	return ErrFlagExit
}
