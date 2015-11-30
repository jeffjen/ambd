package main

import (
	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func NewCancelCmd() cli.Command {
	return cli.Command{
		Name:  "cancel",
		Usage: "Cancel proxy endpoint by source",
		Flags: []cli.Flag{
			cli.StringSliceFlag{Name: "name", Usage: "Proxy endpoint alias"},
		},
		Action: cancel,
	}
}

func cancel(ctx *cli.Context) {
	var failed bool = false

	for _, name := range ctx.StringSlice("name") {
		if err := CancelReq(name); err != nil {
			fmt.Fprintf(os.Stderr, "failed cancel - %s\n", name)
			failed = true
		} else {
			fmt.Printf("%s canceled\n", name)
		}
	}
	if failed {
		os.Exit(1)
	}
}
