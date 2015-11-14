package cmd

import (
	api "github.com/jeffjen/docker-ambassador/api"
	disc "github.com/jeffjen/docker-ambassador/discovery"

	log "github.com/Sirupsen/logrus"
	cli "github.com/codegangsta/cli"
)

func init() {
	disc.RegisterPath = "/srv/ambassador"
}

func runAPIEndpoint(addr string, stop chan<- struct{}) {
	defer close(stop)

	server := api.GetServer()

	server.Addr = addr
	log.Error(server.ListenAndServe())
}

func Ambassador(ctx *cli.Context) {
	var (
		addr = ctx.String("addr")

		stop = make(chan struct{}, 1)
	)

	disc.Check(ctx)

	if addr != "" {
		log.WithFields(log.Fields{"addr": addr}).Info("API endpoint begin")
		runAPIEndpoint(addr, stop)
	} else {
		log.Warning("API endpoint disabled")
	}
}
