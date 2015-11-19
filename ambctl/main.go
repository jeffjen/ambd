package main

import (
	"github.com/jeffjen/docker-ambassador/ambctl/command"

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
	app.Commands = []cli.Command{
		command.NewListCmd(),
		command.NewCreateCmd(),
		command.NewCancelCmd(),
		command.NewInfoCmd(),
	}
	app.Run(os.Args)
}
