package main

import (
	cli "github.com/codegangsta/cli"
)

func init() {
	cli.CommandHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS]{{end}}

{{.Usage}}

{{if .Flags}}Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
`
}
