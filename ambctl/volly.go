package main

import (
	cli "github.com/codegangsta/cli"
	etcd "github.com/coreos/etcd/client"
	ctx "golang.org/x/net/context"

	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

func NewVollyCmd() cli.Command {
	return cli.Command{
		Name:   "volly",
		Usage:  "Send command to cluster nodes",
		Before: nodes,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "cluster", Usage: "Cluster identifier"},
			cli.StringFlag{Name: "dsc", Usage: "Discovery service endpoints", EnvVar: "ETCDCTL_ENDPOINT"},
		},
		Subcommands: []cli.Command{
			NewListCmd(),
			NewCreateCmd(),
			NewCancelCmd(),
			NewInfoCmd(),
			NewConfigCmd(),
		},
	}
}

func nodes(c *cli.Context) error {
	Endpoint = make([]string, 0)

	var (
		cluster = path.Join(c.String("cluster"), "/docker/ambd/nodes")
	)

	cfg := etcd.Config{
		Endpoints: strings.Split(c.String("dsc"), ","),
	}
	cli, err := etcd.New(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	kAPI := etcd.NewKeysAPI(cli)

	wk, abort := ctx.WithTimeout(ctx.Background(), 100*time.Millisecond)
	defer abort()

	resp, err := kAPI.Get(wk, cluster, &etcd.GetOptions{
		Recursive: true,
	})
	if err != nil {
		if nrr, ok := err.(etcd.Error); ok && nrr.Code == 100 {
			fmt.Fprintln(os.Stderr, "cluster not found")
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

	for _, node := range resp.Node.Nodes {
		if !node.Dir {
			host := path.Base(node.Key)
			Endpoint = append(Endpoint, fmt.Sprintf("http://%s", host))
		}
	}

	return nil
}

func volly(c *cli.Context) {
	// NOOP
}
