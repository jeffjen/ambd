package cmd

import (
	proxy "github.com/jeffjen/ambd/proxy"
	web "github.com/jeffjen/ambd/web"
	disc "github.com/jeffjen/go-discovery"
	dcli "github.com/jeffjen/go-discovery/cli"

	log "github.com/Sirupsen/logrus"
	cli "github.com/codegangsta/cli"
)

func Ambassador(ctx *cli.Context) {
	var (
		addr = ctx.String("addr")

		proxycfg = ctx.String("proxycfg")

		proxyTargets = ctx.StringSlice("proxy")

		stop = make(chan struct{}, 1)
	)

	// setup register path for discovery
	disc.RegisterPath = ctx.String("prefix")

	if err := dcli.Before(ctx); err != nil {
		if err == dcli.ErrRequireAdvertise || err == dcli.ErrRequireDiscovery {
			log.WithFields(log.Fields{"err": err}).Warning("discovery feature disabled")
		} else {
			log.WithFields(log.Fields{"err": err}).Fatal("halt")
		}
	}

	if cfgkey := proxy.ConfigKey(); cfgkey != "" {
		proxy.ProxyConfigKey = cfgkey
	} else if proxycfg != "" {
		proxy.ProxyConfigKey = proxycfg
	}
	if proxyTargets != nil {
		proxy.Targets = proxyTargets
	}

	proxy.Follow()

	if addr != "" {
		log.WithFields(log.Fields{"addr": addr}).Info("API endpoint begin")
		web.RunAPIEndpoint(addr, stop)
	} else {
		log.Warning("API endpoint disabled")
	}

	<-stop // we should never reach here

	// TODO: we should gracefully shutdown proxied connections
}
