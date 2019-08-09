package cli

import (
	"os"
	"text/template"
)

var commandHelpTemplate = template.Must(template.New("command help").Parse(`{{ .Name }} - {{ .Description }}
{{- if .HelpText }}

{{ .HelpText }}
{{ end }}
{{ if .Commands -}}
Subcommands:
{{- range .Commands }}
    {{ .Name }} - {{ .Description }}
{{- end }}
{{- end }}
{{- if .Flags }}
Flags:
{{- range .Flags }}
    {{ .Base.Name }} - {{ .Description }} ({{ .Value }})
{{- end }}
{{- end }}
`))

func defaultHelp(c *Command) error {
	err := commandHelpTemplate.Execute(os.Stdout, c)
	if err != nil {
		return err
	}
	return nil
}
