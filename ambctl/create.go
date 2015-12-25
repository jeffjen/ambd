package main

import (
	arg "github.com/jeffjen/ambd/ambctl/arg"

	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
	"path"
)

func NewCreateCmd() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create proxy endpoint by spec",
		Flags: []cli.Flag{
			// Proxy endpoint configuration
			cli.StringFlag{Name: "name", Usage: "Proxy endpoint alias"},
			cli.StringFlag{Name: "net", Usage: "Network type", Value: "tcp"},
			cli.StringSliceFlag{Name: "src", Usage: "Origin address to listen"},
			cli.StringSliceFlag{Name: "dst", Usage: "Target to proxy to"},
			cli.StringFlag{Name: "srv", Usage: "Service identity in discovery"},

			// Role of this service
			cli.StringFlag{Name: "tlsrole", Usage: "Specify ", EnvVar: "PROXY_CERTPATH"},

			// Certificate filepath with required keys
			cli.StringFlag{Name: "tlscertpath", Usage: "Specify certificate path", EnvVar: "PROXY_CERTPATH"},

			// CA, Cert, and Key filepath
			cli.StringFlag{Name: "tlscacert", Usage: "Trust certs signed only by this CA", EnvVar: "PROXY_CA"},
			cli.StringFlag{Name: "tlskey", Usage: "Path to TLS key file", EnvVar: "PROXY_KEY"},
			cli.StringFlag{Name: "tlscert", Usage: "Path to TLS certificate file", EnvVar: "PROXY_CERT"},
		},
		Action: create,
	}
}

func createreq(info *arg.Info) {
	resp, fail := CreateReq(info), false
	for ret := range resp {
		if ret.Err != nil {
			fmt.Fprintf(os.Stderr, "%s - failed create\n", ret.Host)
			fail = true
		} else {
			fmt.Println("done")
		}
	}
	if fail {
		os.Exit(1)
	}
}

func create(c *cli.Context) {
	var (
		// service role
		tlsrole = c.String("tlsrole")

		// certificate path with all necessary info
		tlscertpath = c.String("tlscertpath")

		// individual certificate file path
		tlscacert = c.String("tlscacert")
		tlskey    = c.String("tlskey")
		tlscert   = c.String("tlscert")

		Name = c.String("name")
		Net  = c.String("net")

		FromRange []string
	)

	if tlscertpath != "" {
		tlscacert = path.Join(tlscertpath, "ca.pem")
		tlskey = path.Join(tlscertpath, "key.pem")
		tlscert = path.Join(tlscertpath, "cert.pem")
	}

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
		if len(FromRange) == 1 {
			createreq(&arg.Info{
				ServerRole: tlsrole,
				CA:         tlscacert,
				Cert:       tlscert,
				Key:        tlskey,
				Name:       Name,
				Net:        Net,
				From:       FromRange[0],
				To:         To,
			})
		} else {
			createreq(&arg.Info{
				ServerRole: tlsrole,
				CA:         tlscacert,
				Cert:       tlscert,
				Key:        tlskey,
				Name:       Name,
				Net:        Net,
				FromRange:  FromRange,
				To:         To,
			})
		}
	} else if srv := c.String("srv"); srv != "" {
		if len(FromRange) == 1 {
			createreq(&arg.Info{
				ServerRole: tlsrole,
				CA:         tlscacert,
				Cert:       tlscert,
				Key:        tlskey,
				Name:       Name,
				Net:        Net,
				From:       FromRange[0],
				Service:    srv,
			})
		} else {
			createreq(&arg.Info{
				ServerRole: tlsrole,
				CA:         tlscacert,
				Cert:       tlscert,
				Key:        tlskey,
				Name:       Name,
				Net:        Net,
				FromRange:  FromRange,
				Service:    srv,
			})
		}
	}
}
