package command

import (
	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func NewListCmd() cli.Command {
	return cli.Command{
		Name:   "list",
		Usage:  "List active proxy endpoint",
		Action: list,
	}
}

func list(ctx *cli.Context) {
	if err := ListProxyReq(); err != nil {
		fmt.Fprintln(os.Stderr, "unable to reach server")
	}
}
