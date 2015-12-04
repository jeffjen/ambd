package main

import (
	cli "github.com/codegangsta/cli"
)

func init() {
	cli.AppHelpTemplate = `Usage: {{.Name}} [COMMAND]
{{.Usage}}

Commands:
	{{range .Commands}}{{.Name}}{{ "\t " }}{{.Usage}}
	{{end}}
Run '{{.Name}} help COMMAND' for more information on a command.

`

	cli.CommandHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS]{{else if .ArgsUsage}}[CONFIG_KEY]{{end}}

{{.Usage}}{{if .Flags}}
Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}{{if .ArgsUsage}}

CONFIG - {{.ArgsUsage}}
{{end}}
`
}
