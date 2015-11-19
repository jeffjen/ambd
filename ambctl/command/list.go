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
	fmt.Fprintln(os.Stderr, "Not yet implemented")
	os.Exit(1)
}
