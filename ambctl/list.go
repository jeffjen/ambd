package main

import (
	cli "github.com/codegangsta/cli"

	"bytes"
	"encoding/json"
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

func list(c *cli.Context) {
	resp, fail := ListProxyReq(), false
	for ret := range resp {
		if ret.Err != nil {
			fmt.Fprintf(os.Stderr, "%s - failed list\n", ret.Host)
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
