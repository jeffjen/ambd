package proxy

import (
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

	ProxyStore map[string]*Info
)

func init() {
	RootContext, Cancel = ctx.WithCancel(ctx.Background())

	ProxyStore = make(map[string]*Info)
}

type Info struct {
	Net  string `json:"net"`
	From string `json:"src"`

	// static assignment
	To []string `json:"dst,omitempty"`

	// read from discovery
	Service string `json:"srv,omitempty"`

	Cancel ctx.CancelFunc `json:-`
}

func parse(spec string) (*Info, error) {
	var i Info
	if err := json.Unmarshal([]byte(spec), &i); err != nil {
		return nil, err
	}
	return &i, nil
}

type ProxyFunc func(ctx.Context, *proxy.ConnOptions) error

func Listen(meta *Info) error {
	var (
		handle ProxyFunc
		opt    *proxy.ConnOptions
	)

	if _, ok := ProxyStore[meta.From]; ok {
		return ErrProxyExist
	}

	if meta.Service != "" {
		discovery := &proxy.DiscOptions{
			Service:   meta.Service,
			Endpoints: disc.Endpoints(),
		}
		handle, opt = proxy.Srv, &proxy.ConnOptions{
			Net:       meta.Net,
			From:      meta.From,
			Discovery: discovery,
		}
	} else if len(meta.To) != 0 {
		handle, opt = proxy.To, &proxy.ConnOptions{
			Net:  meta.Net,
			From: meta.From,
			To:   meta.To,
		}
	}

	// Attach context to proxy daemon
	order, abort := ctx.WithCancel(RootContext)

	// This proxy shall have its isolated abort feature
	meta.Cancel = abort

	fields := log.Fields{"Net": meta.Net, "From": meta.From, "To": meta.To, "Service": meta.Service}
	go func() {
		log.WithFields(fields).Info("begin")
		err := handle(order, opt)
		log.WithFields(fields).Warning(err)
	}()

	ProxyStore[meta.From] = meta

	return nil
}

func Reload() {
	for from, meta := range ProxyStore {
		delete(ProxyStore, from)
		meta.Cancel()
		Listen(meta)
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
