package cli

import (
	"text/template"
)

var HelpFlag = &Flag{
	Name:        "help,h",
	Description: "show that message",
	After:       defaultHelp,
	Value:       false,
}

var commandHelpTemplate = template.Must(template.New("command help").Parse(`{{ .Name }} {{ .Usage }} - {{ .Description }}
{{- with .HelpText }}

{{ . -}}
{{ end }}
{{ if .Commands }}
Subcommands:
{{- range .Commands }}
	{{ .Name }} {{ .Usage }}{{ if .Description }} - {{ .Description }}{{ end }}
{{ end }}
{{- end }}
{{- if or .Flags .Parent }}
Flags:
{{- block "flags" . }}
{{- range .Flags }}
    {{ if not .Hidden }}{{ .Name }}{{ with .Value }}={{ . }}{{ end }}				- {{ .Description }}{{ end }}
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
