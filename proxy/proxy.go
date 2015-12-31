package proxy

import (
	"github.com/jeffjen/go-libkv/libkv"

	disc "github.com/jeffjen/go-discovery"
	"github.com/jeffjen/go-proxy/proxy"

	log "github.com/Sirupsen/logrus"
	ctx "golang.org/x/net/context"

	"errors"
)

const (
	DiscoveryPath = "/docker/ambd/nodes"
)

var (
	ErrProxyExist = errors.New("proxy exist")

	ErrMissingName = errors.New("proxy name empty")

	Targets []string

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
	// describe what TLS role this instnace is running
	// Either *leave empty*, which is not engage in TLS
	// or *client*, which is to connect with provided certificate
	// or *server*, which is to listen with provided certificate
	ServerRole string `json:"tlsrole"`

	// certificate root path
	CA string `json:"tlscacert"`

	// certificate public private key pair
	Cert string `json:"tlscert"`
	Key  string `json:"tlskey"`

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
		// proxy type handler
		handle ProxyFunc

		// proxy connection option
		opts *proxy.ConnOptions

		// tls configuration info
		tlscfg proxy.TLSConfig

		logger = log.WithFields(log.Fields{
			"Name":    i.Name,
			"Net":     i.Net,
			"From":    i.From,
			"Range":   i.FromRange,
			"To":      i.To,
			"Service": i.Service,
			"Role":    i.ServerRole,
		})
	)

	// Setup TLS configuration
	switch {
	case i.ServerRole == "server":
		cfg, err := proxy.LoadCertificate(proxy.CertOptions{
			CA:      i.CA,
			TlsCert: i.Cert,
			TlsKey:  i.Key,
			Server:  true,
		})
		if err != nil {
			logger.Error(err)
			return
		}
		tlscfg = proxy.TLSConfig{Server: cfg}
		break
	case i.ServerRole == "client":
		cfg, err := proxy.LoadCertificate(proxy.CertOptions{
			CA:      i.CA,
			TlsCert: i.Cert,
			TlsKey:  i.Key,
		})
		if err != nil {
			logger.Error(err)
			return
		}
		tlscfg = proxy.TLSConfig{Client: cfg}
		break

	case i.ServerRole == "": // not processing TLS
		break

	default: // fall through as if not processing TLS
		break
	}

	// Setup destination and listener information
	switch {
	case i.Service != "":
		discovery := &proxy.DiscOptions{
			Service:   i.Service,
			Endpoints: disc.Endpoints(),
		}
		opts = &proxy.ConnOptions{
			Net:       i.Net,
			Discovery: discovery,
			TLSConfig: tlscfg,
		}
		if len(i.FromRange) != 0 {
			handle = proxy.ClusterSrv
			opts.FromRange = i.FromRange
		} else {
			handle = proxy.Srv
			opts.From = i.From
		}
		break
	case len(i.To) != 0:
		opts = &proxy.ConnOptions{
			Net:       i.Net,
			To:        i.To,
			TLSConfig: tlscfg,
		}
		if len(i.FromRange) != 0 {
			handle = proxy.ClusterTo
			opts.FromRange = i.FromRange
		} else {
			handle = proxy.To
			opts.From = i.From
		}
		break
	}

	// Attach context to proxy daemon
	order, abort := ctx.WithCancel(RootContext)

	// This proxy shall have its isolated abort feature
	i.Cancel = abort

	go func() {
		logger.Info("begin")
		err := handle(order, opts)
		logger.Warning(err)
	}()
}

func Listen(meta *Info) error {
	if Store.Get(meta.Name) != nil {
		return ErrProxyExist
	}
	meta.Listen()
	Store.Set(meta.Name, meta)
	return nil
}
