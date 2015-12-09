package main

import (
	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func NewInfoCmd() cli.Command {
	return cli.Command{
		Name:   "info",
		Usage:  "Get ambd info",
		Action: info,
	}
}

func info(ctx *cli.Context) {
	if err := InfoReq(); err != nil {
		fmt.Fprintln(os.Stderr, "unable to reach server")
	}
}
