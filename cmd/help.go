package cmd

import (
	cli "github.com/codegangsta/cli"
)

func init() {
	cli.AppHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS]{{end}} DISCOVERY_URI

{{.Usage}}

Version: {{.Version}}

DISCOVERY_URI:
	- etcd://10.0.1.10:2379,10.0.1.11:2379

Example proxy specification:
	- {"net": "tcp", "src": ":6379", "dst": ["redis-01:6379", "redis-02:6379"]}
	- {"net": "tcp", "src": ":3370", "srv": "/srv/mysql/debug"}

{{if .Flags}}Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
{{if .Commands}}Commands:
	{{range .Commands}}{{.Name}}{{ "\t " }}{{.Usage}}
	{{end}}
Run '{{.Name}} help COMMAND' for more information on a command.{{end}}

`
}
