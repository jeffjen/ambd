package cmd

import (
	dcli "github.com/jeffjen/go-discovery/cli"

	cli "github.com/codegangsta/cli"
)

var (
	Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Usage: "API endpoint for admin",
		},
		cli.StringSliceFlag{
			Name:  "proxy",
			Usage: "Proxy specification on startup",
		},
	}
)

func NewFlag() []cli.Flag {
	return append(Flags, dcli.Flags...)
}
