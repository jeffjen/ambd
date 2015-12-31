package main

import (
	dcli "github.com/jeffjen/go-discovery/cli"

	cli "github.com/codegangsta/cli"
)

var (
	Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cluster",
			Usage: "cluster to apply for discovery",
			Value: "debug",
		},
		cli.StringFlag{
			Name:  "addr",
			Usage: "API endpoint for admin",
		},
		cli.StringFlag{
			Name:  "proxycfg",
			Usage: "Proxy specification from config key",
		},
		cli.StringSliceFlag{
			Name:  "proxy",
			Usage: "Proxy specification on startup",
		},
		cli.BoolFlag{
			Name:  "proxy2discovery",
			Usage: "Direct ambassador to setup proxy for discovery",
		},
	}
)

func NewFlag() []cli.Flag {
	return append(Flags, dcli.Flags...)
}
