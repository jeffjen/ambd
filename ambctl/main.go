package main

import (
	cli "github.com/codegangsta/cli"

	"fmt"
	"os"
)

func endpoint(c *cli.Context) error {
	if host := c.String("host"); host != "" {
		Endpoint = append(Endpoint, fmt.Sprintf("http://%s", host))
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "ambctl"
	app.Usage = "Admin tool for ambd"
	app.Authors = []cli.Author{
		cli.Author{"Yi-Hung Jen", "yihungjen@gmail.com"},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "host", Usage: "ambd host", Value: "localhost:29091"},
	}
	app.Commands = []cli.Command{
		NewVollyCmd(),
		NewListCmd(),
		NewCreateCmd(),
		NewCancelCmd(),
		NewInfoCmd(),
		NewConfigCmd(),
	}
	app.Before = endpoint
	app.Run(os.Args)
}
