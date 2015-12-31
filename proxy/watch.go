package proxy

import (
	disc "github.com/jeffjen/go-discovery"
	"github.com/jeffjen/go-libkv/libkv"
	"github.com/jeffjen/go-proxy/proxy"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	ctx "golang.org/x/net/context"

	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
)

var (
	ProxyConfigKey string = "__nobody__"

	ConfigReset ctx.CancelFunc

	retry = &proxy.Backoff{}
)

var (
	EnableDiscoveryProxy bool

	DiscoveryProxyInfo = &Info{
		Name: "discovery",
		Net:  "tcp4",
		From: ":2379",
	}
)

func splitDiscovery() (dst []string) {
	return strings.Split(strings.TrimPrefix(disc.Discovery, "etcd://"), ",")
}

func DiscoveryURI() {
	if buf, err := ioutil.ReadFile(".discovery"); err == nil {
		disc.Discovery = string(buf)
	}
	go func() {
		for _ = range time.Tick(2 * time.Minute) {
			ioutil.WriteFile(".discovery", []byte(disc.Discovery), 0644)
		}
	}()
}

func ConfigKey() string {
	var cfgkey string
	if buf, err := ioutil.ReadFile(".proxycfg"); err == nil {
		cfgkey = string(buf)
	}
	go func() {
		for _ = range time.Tick(2 * time.Minute) {
			ioutil.WriteFile(".proxycfg", []byte(ProxyConfigKey), 0644)
		}
	}()
	return cfgkey
}

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
		log.WithFields(log.Fields{"err": err, "value": value}).Warning("config")
		targets = nil
	}
	return
}

func reload(pxycfg []*Info) {
	it, mod := Store.IterateW()
	for elem := range it {
		mod <- &libkv.Value{R: true}
		elem.X.(*Info).Cancel()
	}
	log.WithFields(log.Fields{"count": len(Targets)}).Debug("reload from args")
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
	log.WithFields(log.Fields{"count": len(pxycfg)}).Debug("reload from cfgkey")
	for _, meta := range pxycfg {
		if err := Listen(meta); err != nil {
			if err != ErrProxyExist {
				log.WithFields(log.Fields{"err": err}).Debug("reload")
			}
		}
	}
	if EnableDiscoveryProxy {
		log.Debug("reload discovery proxy")
		DiscoveryProxyInfo.To = splitDiscovery()
		if err := Listen(DiscoveryProxyInfo); err != nil {
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
			reload(o)
		}
	}()
	return order
}

func watcheWorker(c ctx.Context, watcher etcd.Watcher) <-chan []*Info {
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
				log.WithFields(log.Fields{"key": evt.Node.Key}).Warning("cfgkey")
				switch evt.Action {
				default:
					v <- nil
					break
				case "set":
					v <- get(evt.Node.Value)
					break
				case "del":
				case "expire":
					v <- make([]*Info, 0)
					break
				}
			}
		}
	}()
	return v
}

func followBootStrap() {
	kAPI, err := disc.NewKeysAPI(etcd.Config{Endpoints: disc.Endpoints()})
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warning("bootstrap")
		return
	}
	resp, err := kAPI.Get(RootContext, ProxyConfigKey, nil)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warning("bootstrap")
		reload(make([]*Info, 0))
	} else if resp.Node.Dir {
		log.WithFields(log.Fields{"key": resp.Node.Key}).Warning("not a valid node")
		reload(make([]*Info, 0))
	} else {
		log.WithFields(log.Fields{"key": resp.Node.Key, "val": resp.Node.Value}).Debug("cfgkey")
		if pxycfg := get(resp.Node.Value); pxycfg != nil {
			reload(pxycfg)
		} else {
			reload(make([]*Info, 0))
		}
	}
}

func Follow() {
	followBootStrap() // bootstrap proxy config

	var c ctx.Context

	c, ConfigReset = ctx.WithCancel(RootContext)
	go func() {
		watcher, err := disc.NewWatcher(&disc.WatcherOptions{
			Config: etcd.Config{Endpoints: disc.Endpoints()},
			Key:    ProxyConfigKey,
		})
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warning("config")
			return
		}
		order := reloadWorker()
		defer close(order)
		for yay := true; yay; {
			v := watcheWorker(c, watcher)
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
