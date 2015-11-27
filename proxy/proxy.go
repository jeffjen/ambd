package proxy

import (
	"github.com/jeffjen/go-libkv/libkv"

	disc "github.com/jeffjen/go-discovery"
	"github.com/jeffjen/go-proxy/proxy"

	log "github.com/Sirupsen/logrus"
	ctx "golang.org/x/net/context"

	"encoding/json"
	"errors"
)

var (
	ErrProxyExist = errors.New("proxy exist")

	Cancel      ctx.CancelFunc
	RootContext ctx.Context

	Store *libkv.Store
)

func init() {
	RootContext, Cancel = ctx.WithCancel(ctx.Background())

	Store = libkv.NewStore()
}

type ProxyFunc func(ctx.Context, *proxy.ConnOptions) error

type Info struct {
	Net  string `json:"net"`
	From string `json:"src"`

	// static assignment
	To []string `json:"dst,omitempty"`

	// read from discovery
	Service string `json:"srv,omitempty"`

	// Proxy handler
	handle ProxyFunc          `json:"-"`
	opts   *proxy.ConnOptions `json:"-"`

	// Abort proxy handler
	Cancel ctx.CancelFunc `json:"-"`
}

func (i *Info) Listen() {
	if i.Service != "" {
		discovery := &proxy.DiscOptions{
			Service:   i.Service,
			Endpoints: disc.Endpoints(),
		}
		i.handle = proxy.Srv
		i.opts = &proxy.ConnOptions{
			Net:       i.Net,
			From:      i.From,
			Discovery: discovery,
		}
	} else if len(i.To) != 0 {
		i.handle = proxy.To
		i.opts = &proxy.ConnOptions{
			Net:  i.Net,
			From: i.From,
			To:   i.To,
		}
	}

	// Attach context to proxy daemon
	order, abort := ctx.WithCancel(RootContext)

	// This proxy shall have its isolated abort feature
	i.Cancel = abort

	fields := log.Fields{"Net": i.Net, "From": i.From, "To": i.To, "Service": i.Service}
	go func() {
		log.WithFields(fields).Info("begin")
		err := i.handle(order, i.opts)
		log.WithFields(fields).Warning(err)
	}()
}

func parse(spec string) (*Info, error) {
	var i = new(Info)
	if err := json.Unmarshal([]byte(spec), i); err != nil {
		return nil, err
	}
	return i, nil
}

func Listen(meta *Info) error {
	if Store.Get(meta.From) != nil {
		return ErrProxyExist
	}
	meta.Listen()
	Store.Set(meta.From, meta)
	return nil
}

func Reload() {
	Store.IterateFunc(func(iden string, x interface{}) {
		meta := x.(*Info)
		meta.Cancel()
		meta.Listen()
	})
}

func RunProxyDaemon(targets []string) {
	for _, spec := range targets {
		meta, err := parse(spec)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Warning("RunProxyDaemon")
			continue
		}
		if err = Listen(meta); err != nil {
			if err != ErrProxyExist {
				log.WithFields(log.Fields{"err": err}).Warning("RunProxyDaemon")
			}
		}
	}
}
