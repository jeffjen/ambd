package cmd

import (
	proxy "github.com/jeffjen/docker-ambassador/proxy"
	web "github.com/jeffjen/docker-ambassador/web"
	dcli "github.com/jeffjen/go-discovery/cli"

	log "github.com/Sirupsen/logrus"
	cli "github.com/codegangsta/cli"

	"os"
)

func Ambassador(ctx *cli.Context) {
	var (
		addr = ctx.String("addr")

		proxycfg = ctx.String("proxycfg")

		proxyTargets = ctx.StringSlice("proxy")

		stop = make(chan struct{}, 1)
	)

	if err := dcli.Before(ctx); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if proxycfg != "" {
		proxy.ProxyConfigKey = proxycfg
		proxy.Follow()
	}

	if proxyTargets != nil {
		proxy.Targets = proxyTargets
		proxy.RunProxyDaemon()
	}

	if addr != "" {
		log.WithFields(log.Fields{"addr": addr}).Info("API endpoint begin")
		web.RunAPIEndpoint(addr, stop)
	} else {
		log.Warning("API endpoint disabled")
	}

	<-stop // we should never reach here

	// TODO: we should gracefully shutdown proxied connections
}
