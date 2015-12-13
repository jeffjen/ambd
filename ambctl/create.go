package main

import (
	arg "github.com/jeffjen/ambd/ambctl/arg"

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

func createreq(info *arg.Info) {
	resp, fail := CreateReq(info), false
	for ret := range resp {
		if ret.Err != nil {
			fmt.Fprintln(os.Stderr, ret.Err)
			fail = true
		} else {
			fmt.Println("done")
		}
	}
	if fail {
		os.Exit(1)
	}
}

func issue(Name, Net string, From, To []string) {
	if len(From) == 1 {
		createreq(&arg.Info{Name: Name, Net: Net, From: From[0], To: To})
	} else {
		createreq(&arg.Info{Name: Name, Net: Net, FromRange: From, To: To})
	}
}

func issueSrv(Name, Net string, From []string, Srv string) {
	if len(From) == 1 {
		createreq(&arg.Info{Name: Name, Net: Net, From: From[0], Service: Srv})
	} else {
		createreq(&arg.Info{Name: Name, Net: Net, FromRange: From, Service: Srv})
	}
}

func create(c *cli.Context) {
	var (
		Name = c.String("name")
		Net  = c.String("net")

		FromRange []string
	)

	if Net == "" {
		fmt.Fprintln(os.Stderr, "missing required flag --net")
		os.Exit(1)
	}
	if from := c.StringSlice("src"); len(from) != 0 {
		FromRange = make([]string, len(from))
		for idx, one_from := range from {
			FromRange[idx] = one_from
		}
	}
	if len(FromRange) == 0 {
		fmt.Fprintln(os.Stderr, "missing required flag --src")
		os.Exit(1)
	}
	if dst := c.StringSlice("dst"); len(dst) != 0 {
		var To = make([]string, len(dst))
		for idx, one_dst := range dst {
			To[idx] = one_dst
		}
		issue(Name, Net, FromRange, To)
	} else if srv := c.String("srv"); srv != "" {
		issueSrv(Name, Net, FromRange, srv)
	}
}
