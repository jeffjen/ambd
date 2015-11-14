package main

import (
	"github.com/jeffjen/docker-ambassador/cmd"

	cli "github.com/codegangsta/cli"

	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "docker-ambassador"
	app.Usage = "Facilitate service consumer discovery and master election"
	app.Authors = []cli.Author{
		cli.Author{"Yi-Hung Jen", "yihungjen@gmail.com"},
	}
	app.Flags = cmd.Flags
	app.Commands = cmd.Commands
	app.Action = cmd.Ambassador
	app.Run(os.Args)
}
