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

	ErrMissingName = errors.New("proxy name empty")

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
	Name string `json:"name"`

	Net       string   `json:"net"`
	From      string   `json:"src"`
	FromRange []string `json:"range"`

	// static assignment
	To []string `json:"dst,omitempty"`

	// read from discovery
	Service string `json:"srv,omitempty"`

	// Abort proxy handler
	Cancel ctx.CancelFunc `json:"-"`
}

func (i *Info) Listen() {
	var (
		handle ProxyFunc
		opts   *proxy.ConnOptions
	)

	if i.Service != "" {
		discovery := &proxy.DiscOptions{
			Service:   i.Service,
			Endpoints: disc.Endpoints(),
		}
		opts = &proxy.ConnOptions{
			Net:       i.Net,
			Discovery: discovery,
		}
		if len(i.FromRange) != 0 {
			handle = proxy.ClusterSrv
			opts.FromRange = i.FromRange
		} else {
			handle = proxy.Srv
			opts.From = i.From
		}
	} else if len(i.To) != 0 {
		opts = &proxy.ConnOptions{
			Net: i.Net,
			To:  i.To,
		}
		if len(i.FromRange) != 0 {
			handle = proxy.ClusterTo
			opts.FromRange = i.FromRange
		} else {
			handle = proxy.To
			opts.From = i.From
		}
	}

	// Attach context to proxy daemon
	order, abort := ctx.WithCancel(RootContext)

	// This proxy shall have its isolated abort feature
	i.Cancel = abort

	fields := log.Fields{"Name": i.Name, "Net": i.Net, "From": i.From, "Range": i.FromRange, "To": i.To, "Service": i.Service}
	go func() {
		log.WithFields(fields).Info("begin")
		err := handle(order, opts)
		log.WithFields(fields).Warning(err)
	}()
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

func Listen(meta *Info) error {
	if Store.Get(meta.Name) != nil {
		return ErrProxyExist
	}
	meta.Listen()
	Store.Set(meta.Name, meta)
	return nil
}

func Reload() {
	for it := range Store.IterateR() {
		meta := it.X.(*Info)
		meta.Cancel()
		meta.Listen()
	}
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
