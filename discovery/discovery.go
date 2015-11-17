package discovery

import (
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	ctx "golang.org/x/net/context"

	"path"
	"strings"
	"time"
)

var (
	RegisterPath = "/srv/monitor"

	Advertise string
	Discovery string
	Endpoints []string

	Hearbeat time.Duration
	TTL      time.Duration

	Cancel ctx.CancelFunc
)

func parse(endpoint string) []string {
	parts := strings.Split(strings.TrimPrefix(endpoint, "etcd://"), ",")
	for idx, p := range parts {
		parts[idx] = "http://" + p
	}
	return parts
}

func NewDiscovery() (client etcd.Client) {
	Endpoints = parse(Discovery)
	cfg := etcd.Config{
		Endpoints:               Endpoints,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	if cli, err := etcd.New(cfg); err != nil {
		log.Fatal(err)
	} else {
		client = cli
	}
	return
}

func upkeep(kAPI etcd.KeysAPI, iden string, opts *etcd.SetOptions) (err error) {
	_, err = kAPI.Set(ctx.Background(), iden, "", opts)
	return
}

func Register(heartbeat time.Duration, ttl time.Duration) {
	Hearbeat, TTL = heartbeat, ttl

	// begin book keeping "THIS" montior unit
	go func() {
		var work ctx.Context

		work, Cancel = ctx.WithCancel(ctx.Background())

		client := NewDiscovery()

		var (
			iden = path.Join(RegisterPath, Advertise)
			opts = etcd.SetOptions{TTL: ttl}
			kAPI = etcd.NewKeysAPI(client)
			f    = log.Fields{"heartbeat": heartbeat, "ttl": ttl}
			t    = time.NewTicker(heartbeat)
		)
		defer t.Stop()

		// Allow for implicit bootstrap on register
		if err := upkeep(kAPI, iden, &opts); err != nil {
			log.Error("1:", err)
		} else {
			log.WithFields(f).Info("uptime")
		}

		// only update key TTL
		opts.PrevExist = etcd.PrevExist

		// Tick... Toc...
		for {
			select {
			case <-t.C:
				if err := upkeep(kAPI, iden, &opts); err != nil {
					log.Error("2:", err)
				} else {
					log.WithFields(f).Info("uptime")
				}
			case <-work.Done():
				log.WithFields(f).Info("abort")
				return
			}
		}
	}()
}
