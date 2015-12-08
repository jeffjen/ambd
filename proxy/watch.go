package proxy

import (
	disc "github.com/jeffjen/go-discovery"
	"github.com/jeffjen/go-libkv/libkv"
	"github.com/jeffjen/go-proxy/proxy"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	ctx "golang.org/x/net/context"

	"encoding/json"
)

var (
	ProxyConfigKey string

	ConfigReset ctx.CancelFunc

	retry = &proxy.Backoff{}
)

func parse(spec string) (*Info, error) {
	var i = new(Info)
	if err := json.Unmarshal([]byte(spec), i); err != nil {
		return nil, err
	}
	if i.Name == "" {
		return nil, ErrMissingName
	}
	return i, nil
}

func get(value string) (targets []*Info) {
	targets = make([]*Info, 0)
	if err := json.Unmarshal([]byte(value), &targets); err != nil {
		log.WithFields(log.Fields{"err": err}).Warning("bad proxy spec")
		targets = nil
	}
	return
}

func doReload(pxycfg []*Info) {
	it, mod := Store.IterateW()
	for elem := range it {
		mod <- &libkv.Value{R: true}
		elem.X.(*Info).Cancel()
	}
	for _, spec := range Targets {
		meta, err := parse(spec)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warning("reload")
			continue
		}
		if err = Listen(meta); err != nil {
			if err != ErrProxyExist {
				log.WithFields(log.Fields{"err": err}).Debug("reload")
			}
		}
	}
	for _, meta := range pxycfg {
		if err := Listen(meta); err != nil {
			if err != ErrProxyExist {
				log.WithFields(log.Fields{"err": err}).Debug("reload")
			}
		}
	}
}

func reloadWorker() chan<- []*Info {
	order := make(chan []*Info)
	go func() {
		for o := range order {
			doReload(o)
		}
	}()
	return order
}

func doWatch(c ctx.Context, watcher etcd.Watcher) <-chan []*Info {
	v := make(chan []*Info)
	go func() {
		evt, err := watcher.Next(c)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Debug("config")
			retry.Delay()
			v <- nil
		} else {
			retry.Reset()
			if evt.Node.Dir {
				log.WithFields(log.Fields{"key": evt.Node.Key}).Warning("not a valid node")
				v <- nil
			} else {
				// FIXME: check that is is not a del or expire
				v <- get(evt.Node.Value)
			}
		}
	}()
	return v
}

func followBootStrap() {
	cfg := etcd.Config{Endpoints: disc.Endpoints()}
	kAPI, err := proxy.NewKeysAPI(cfg)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warning("bootstrap")
		return
	}
	resp, err := kAPI.Get(RootContext, ProxyConfigKey, nil)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warning("bootstrap")
		doReload(make([]*Info, 0))
	} else if resp.Node.Dir {
		log.WithFields(log.Fields{"key": resp.Node.Key}).Warning("not a valid node")
		doReload(make([]*Info, 0))
	} else {
		doReload(get(resp.Node.Value))
	}
}

func Follow() {
	followBootStrap() // bootstrap proxy config

	var c ctx.Context

	c, ConfigReset = ctx.WithCancel(RootContext)
	go func() {
		cfg := etcd.Config{Endpoints: disc.Endpoints()}
		watcher, err := proxy.NewWatcher(cfg, ProxyConfigKey, 0)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warning("config")
			return
		}
		order := reloadWorker()
		defer close(order)
		for yay := true; yay; {
			v := doWatch(c, watcher)
			select {
			case <-c.Done():
				yay = false
			case proxyTargets, ok := <-v:
				if ok && proxyTargets != nil {
					order <- proxyTargets
				}
				yay = ok
			}
		}
	}()
}
