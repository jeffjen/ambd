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
	app.Commands = []cli.Command{
		NewListCmd(),
		NewCreateCmd(),
		NewCancelCmd(),
		NewInfoCmd(),
		NewConfigCmd(),
	}
	app.Run(os.Args)
}
