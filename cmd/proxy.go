package cmd

import (
	disc "github.com/jeffjen/docker-ambassador/discovery"
	"github.com/jeffjen/go-proxy/proxy"

	log "github.com/Sirupsen/logrus"
	ctx "golang.org/x/net/context"

	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"
)

var (
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

func Listen(iden string, meta *Info) {
	var (
		handle ProxyFunc
		opt    *proxy.ConnOptions
	)

	if meta.Service != "" {
		discovery := &proxy.DiscOptions{
			Service:   meta.Service,
			Endpoints: disc.Endpoints,
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

	ProxyStore[iden] = meta
}

func Genkey(seed string) (key string) {
	key = fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().String()+seed)))
	return
}

func runProxyDaemon(targets []string) {
	for _, spec := range targets {
		meta, err := parse(spec)
		if err != nil {
			log.Warning(err)
		} else {
			Listen(Genkey(""), meta)
		}
	}
}
