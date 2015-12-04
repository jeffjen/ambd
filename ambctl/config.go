package main

import (
	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func NewConfigCmd() cli.Command {
	return cli.Command{
		Name:      "config",
		Usage:     "Command ambassador to use this config",
		Action:    config,
		ArgsUsage: "Key for config value",
	}
}

func config(ctx *cli.Context) {
	var proxycfg []string = ctx.Args()

	if len(proxycfg) == 0 {
		fmt.Fprintln(os.Stderr, "Must have exactly one argument")
		os.Exit(1)
	}

	if err := ConfigReq(proxycfg[0]); err != nil {
		fmt.Fprintln(os.Stderr, "failed request")
		os.Exit(1)
	} else {
		fmt.Println("done")
	}
}
