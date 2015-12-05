package main

import (
	cli "github.com/codegangsta/cli"

	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ambctl"
	app.Usage = "Admin tool for docker-ambassador"
	app.Authors = []cli.Author{
		cli.Author{"Yi-Hung Jen", "yihungjen@gmail.com"},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "host", Usage: "docker-ambassador host", Value: "localhost:29091"},
	}
	app.Commands = []cli.Command{
		NewListCmd(),
		NewCreateCmd(),
		NewCancelCmd(),
		NewInfoCmd(),
		NewConfigCmd(),
	}
	app.Before = endpoint
	app.Run(os.Args)
}
