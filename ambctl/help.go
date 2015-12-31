package main

import (
	cli "github.com/codegangsta/cli"
)

func init() {
	cli.AppHelpTemplate = `Usage: {{.Name}} [COMMAND]

{{.Usage}}{{if .Flags}}
Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
Commands:
	{{range .Commands}}{{.Name}}{{ "\t " }}{{.Usage}}
	{{end}}
Run '{{.Name}} help COMMAND' for more information on a command.

`

	cli.CommandHelpTemplate = `Usage: {{.Name}}{{if .ArgsUsage}} [CONFIG_KEY]{{end}}{{if .Flags}} [OPTIONS]{{end}}

{{.Usage}}{{if .ArgsUsage}}

CONFIG - {{.ArgsUsage}}
{{end}}{{if .Flags}}
Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}{{if .Subcommands}}
Commands:
	{{range .Subcommands}}{{.Name}}{{ "\t " }}{{.Usage}}
	{{end}}{{end}}
`
}
