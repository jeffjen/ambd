package main

import (
	"github.com/jeffjen/ambd/cmd"

	cli "github.com/codegangsta/cli"

	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ambd"
	app.Usage = "Facilitate dynamic Ambassador pattern"
	app.Authors = []cli.Author{
		cli.Author{"Yi-Hung Jen", "yihungjen@gmail.com"},
	}
	app.Flags = cmd.NewFlag()
	app.Action = cmd.Ambassador
	app.Run(os.Args)
}
