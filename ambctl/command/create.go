package command

import (
	pxy "github.com/jeffjen/docker-ambassador/proxy"

	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func NewCreateCmd() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create proxy endpoint by spec",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "net", Usage: "Network type", Value: "tcp"},
			cli.StringFlag{Name: "src", Usage: "Origin address to listen"},
			cli.StringSliceFlag{Name: "dst", Usage: "Target to proxy to"},
			cli.StringFlag{Name: "srv", Usage: "Service identity in discovery"},
		},
		Action: create,
	}
}

func create(ctx *cli.Context) {
	var (
		Net  = ctx.String("net")
		From = ctx.String("src")
	)

	if Net == "" {
		fmt.Fprintln(os.Stderr, "missing required flag --net")
		os.Exit(1)
	}
	if From == "" {
		fmt.Fprintln(os.Stderr, "missing required flag --net")
		os.Exit(1)
	}
	if dst := ctx.StringSlice("dst"); len(dst) != 0 {
		var To = make([]string, len(dst))
		for idx, one_dst := range dst {
			To[idx] = one_dst
		}
		if err := CreateReq(pxy.Info{Net: Net, From: From, To: To}); err != nil {
			fmt.Fprintln(os.Stderr, "failed request")
			os.Exit(1)
		} else {
			fmt.Println("done")
		}
	}
}
