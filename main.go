package main

import (
	proxy "github.com/jeffjen/ambd/proxy"
	web "github.com/jeffjen/ambd/web"
	disc "github.com/jeffjen/go-discovery"
	dcli "github.com/jeffjen/go-discovery/cli"

	log "github.com/Sirupsen/logrus"
	cli "github.com/codegangsta/cli"

	"os"
	"path"
)

func main() {
	app := cli.NewApp()
	app.Name = "ambd"
	app.Usage = "Facilitate dynamic Ambassador pattern"
	app.Authors = []cli.Author{
		cli.Author{"Yi-Hung Jen", "yihungjen@gmail.com"},
	}
	app.Flags = NewFlag()
	app.Action = Ambassador
	app.Run(os.Args)
}

func Ambassador(ctx *cli.Context) {
	var (
		addr = ctx.String("addr")

		proxycfg = ctx.String("proxycfg")

		proxyTargets = ctx.StringSlice("proxy")

		stop = make(chan struct{}, 1)
	)

	// setup register path for discovery
	disc.RegisterPath = path.Join(ctx.String("cluster"), proxy.DiscoveryPath)

	if err := dcli.Before(ctx); err != nil { // We don't want to setup discovery
		if err == dcli.ErrRequireDiscovery {
			log.WithFields(log.Fields{"err": err}).Warning("discovery feature disabled")
		} else {
			log.WithFields(log.Fields{"err": err}).Fatal("halt")
		}
	} else { // We had successfully setup discovery
		if cfgkey := proxy.ConfigKey(); cfgkey != "" {
			proxy.ProxyConfigKey = cfgkey
		} else if proxycfg != "" {
			proxy.ProxyConfigKey = proxycfg
		}
		if proxyTargets != nil {
			proxy.Targets = proxyTargets
		}
		proxy.EnableDiscoveryProxy = ctx.Bool("proxy2discovery")
		proxy.Follow()
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
