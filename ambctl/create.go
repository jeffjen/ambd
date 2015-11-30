package main

import (
	arg "github.com/jeffjen/docker-ambassador/ambctl/arg"

	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func NewCreateCmd() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create proxy endpoint by spec",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name", Usage: "Proxy endpoint alias"},
			cli.StringFlag{Name: "net", Usage: "Network type", Value: "tcp"},
			cli.StringSliceFlag{Name: "src", Usage: "Origin address to listen"},
			cli.StringSliceFlag{Name: "dst", Usage: "Target to proxy to"},
			cli.StringFlag{Name: "srv", Usage: "Service identity in discovery"},
		},
		Action: create,
	}
}

func issue(Name, Net string, From, To []string) {
	if len(From) == 1 {
		if err := CreateReq(arg.Info{Name: Name, Net: Net, From: From[0], To: To}); err != nil {
			fmt.Fprintln(os.Stderr, "failed request")
			os.Exit(1)
		} else {
			fmt.Println("done")
		}
	} else {
		if err := CreateReq(arg.Info{Name: Name, Net: Net, FromRange: From, To: To}); err != nil {
			fmt.Fprintln(os.Stderr, "failed request")
			os.Exit(1)
		} else {
			fmt.Println("done")
		}
	}
}

func issueSrv(Name, Net string, From []string, Srv string) {
	if len(From) == 1 {
		if err := CreateReq(arg.Info{Name: Name, Net: Net, From: From[0], Service: Srv}); err != nil {
			fmt.Fprintln(os.Stderr, "failed request")
			os.Exit(1)
		} else {
			fmt.Println("done")
		}
	} else {
		if err := CreateReq(arg.Info{Name: Name, Net: Net, FromRange: From, Service: Srv}); err != nil {
			fmt.Fprintln(os.Stderr, "failed request")
			os.Exit(1)
		} else {
			fmt.Println("done")
		}
	}
}

func create(ctx *cli.Context) {
	var (
		Name = ctx.String("name")
		Net  = ctx.String("net")

		FromRange []string
	)

	if Net == "" {
		fmt.Fprintln(os.Stderr, "missing required flag --net")
		os.Exit(1)
	}
	if from := ctx.StringSlice("src"); len(from) != 0 {
		FromRange = make([]string, len(from))
		for idx, one_from := range from {
			FromRange[idx] = one_from
		}
		fmt.Println(FromRange)
	}
	if len(FromRange) == 0 {
		fmt.Fprintln(os.Stderr, "missing required flag --src")
		os.Exit(1)
	}
	if dst := ctx.StringSlice("dst"); len(dst) != 0 {
		var To = make([]string, len(dst))
		for idx, one_dst := range dst {
			To[idx] = one_dst
		}
		issue(Name, Net, FromRange, To)
	} else if srv := ctx.String("srv"); srv != "" {
		issueSrv(Name, Net, FromRange, srv)
	}
}
