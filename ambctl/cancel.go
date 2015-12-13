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
	fail := false
	for _, name := range ctx.StringSlice("name") {
		resp := CancelReq(name)
		for ret := range resp {
			if ret.Err != nil {
				fmt.Fprintf(os.Stderr, "%s - failed cancel: %s\n", ret.Host, name)
				fail = true
			} else {
				fmt.Printf("%s canceled\n", name)
			}
		}
	}
	if fail {
		os.Exit(1)
	}
}
