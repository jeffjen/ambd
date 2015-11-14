package cmd

import (
	cli "github.com/codegangsta/cli"
)

func init() {
	cli.AppHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS]{{end}} DISCOVERY_URI

{{.Usage}}

Version: {{.Version}}

DISCOVERY_URI:
	EXAMPLE URI - etcd://10.0.1.10:2379,10.0.1.11:2379

{{if .Flags}}Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
{{if .Commands}}Commands:
	{{range .Commands}}{{.Name}}{{ "\t " }}{{.Usage}}
	{{end}}
Run '{{.Name}} help COMMAND' for more information on a command.{{end}}
`

	cli.CommandHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS]{{end}} DISCOVERY_URI

{{.Usage}}

DISCOVERY_URI:
	EXAMPLE URI - etcd://10.0.1.10:2379,10.0.1.11:2379

{{if .Flags}}Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
`
}
