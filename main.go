package main

import (
	"github.com/jeffjen/docker-ambassador/cmd"
	disc "github.com/jeffjen/go-discovery"

	cli "github.com/codegangsta/cli"

	"os"
)

func init() {
	disc.RegisterPath = "/srv/ambassador"
}

func main() {
	app := cli.NewApp()
	app.Name = "docker-ambassador"
	app.Usage = "Facilitate dynamic Ambassador pattern"
	app.Authors = []cli.Author{
		cli.Author{"Yi-Hung Jen", "yihungjen@gmail.com"},
	}
	app.Flags = cmd.NewFlag()
	app.Commands = cmd.Commands
	app.Action = cmd.Ambassador
	app.Run(os.Args)
}
