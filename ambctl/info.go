package main

import (
	cli "github.com/codegangsta/cli"

	"bytes"
	"encoding/json"
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

func info(c *cli.Context) {
	resp, fail := InfoReq(), false
	for ret := range resp {
		if ret.Err != nil {
			fmt.Fprintf(os.Stderr, "%s - failed info\n", ret.Host)
			fail = true
		} else {
			var out = new(bytes.Buffer)
			json.Indent(out, ret.Data, "", "    ")
			out.WriteTo(os.Stdout)
		}
	}
	if fail {
		os.Exit(1)
	}
}
