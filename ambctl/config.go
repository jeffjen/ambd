package main

import (
	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func NewConfigCmd() cli.Command {
	return cli.Command{
		Name:  "config",
		Usage: "Command ambassador to use this config",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "dsc", Usage: "Discovery service endpoint"},
			cli.StringFlag{Name: "cluster", Usage: "Set cluster node it belongs to"},
			cli.BoolFlag{Name: "proxy2discovery", Usage: "Enable proxy to discovery"},
		},
		Action:    config,
		ArgsUsage: "Key for config value",
	}
}

func config(ctx *cli.Context) {
	var (
		proxycfg []string = ctx.Args()

		dsc     = ctx.String("dsc")
		cluster = ctx.String("cluster")
		enable  = ctx.Bool("proxy2discovery")
	)

	if len(proxycfg) == 0 {
		fmt.Fprintln(os.Stderr, "Must have exactly one argument")
		os.Exit(1)
	}
	if dsc == "" {
		dsc = "null"
	}

	resp, fail := ConfigReq(proxycfg[0], dsc, cluster, enable), false
	for ret := range resp {
		if ret.Err != nil {
			fmt.Fprintf(os.Stderr, "%s - failed config\n", ret.Host)
			fail = true
		} else {
			fmt.Println("done")
		}
	}
	if fail {
		os.Exit(1)
	}
}
