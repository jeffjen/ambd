package cmd

import (
	api "github.com/jeffjen/docker-ambassador/api"
	disc "github.com/jeffjen/docker-ambassador/discovery"
	proxy "github.com/jeffjen/docker-ambassador/proxy"

	log "github.com/Sirupsen/logrus"
	cli "github.com/codegangsta/cli"
)

func init() {
	disc.RegisterPath = "/srv/ambassador"
}

func Ambassador(ctx *cli.Context) {
	var (
		addr = ctx.String("addr")

		proxyTargets = ctx.StringSlice("proxy")

		stop = make(chan struct{}, 1)
	)

	disc.Check(ctx)

	if proxyTargets != nil {
		proxy.RunProxyDaemon(proxyTargets)
	} else {
		log.Info("No proxy at startup")
	}

	if addr != "" {
		log.WithFields(log.Fields{"addr": addr}).Info("API endpoint begin")
		api.RunAPIEndpoint(addr, stop)
	} else {
		log.Warning("API endpoint disabled")
	}

	<-stop // we should never reach here

	// TODO: we should gracefully shutdown proxied connections
}
