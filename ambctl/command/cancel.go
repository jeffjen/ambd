package command

import (
	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func NewCancelCmd() cli.Command {
	return cli.Command{
		Name:   "cancel",
		Usage:  "Cancel proxy endpoint by source",
		Action: cancel,
	}
}

func cancel(ctx *cli.Context) {
	var failed bool = false
	for _, src := range ctx.Args() {
		if err := CancelReq(src); err != nil {
			fmt.Fprintf(os.Stderr, "failed cancel - %s\n", src)
			failed = true
		} else {
			fmt.Printf("%s canceled\n", src)
		}
	}
	if failed {
		os.Exit(1)
	}
}
