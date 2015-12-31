package service

import (
	proxy "github.com/jeffjen/ambd/proxy"
	disc "github.com/jeffjen/go-discovery"

	_ "github.com/Sirupsen/logrus"

	"net/http"
	"time"
)

func Follow(w http.ResponseWriter, r *http.Request) {
	if err := common("PUT", r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var (
		proxycfg = r.Form.Get("key")

		discovery = r.Form.Get("discovery")

		heartbeat time.Duration
		ttl       time.Duration
	)

	if proxycfg == "" {
		http.Error(w, "Bad Request", 400)
		return
	} else {
		proxy.ProxyConfigKey = proxycfg
	}

	if hbStr := r.Form.Get("hb"); hbStr == "" {
		heartbeat = disc.DefaultHeartbeat
	} else {
		if hb, err := time.ParseDuration(hbStr); err != nil {
			heartbeat = disc.DefaultHeartbeat
		} else {
			heartbeat = hb
		}
	}

	if ttlStr := r.Form.Get("ttl"); ttlStr == "" {
		ttl = disc.DefaultTTL
	} else {
		if t, err := time.ParseDuration(ttlStr); err != nil {
			ttl = disc.DefaultTTL
		} else {
			ttl = t
		}
	}

	if proxy.ConfigReset != nil {
		proxy.ConfigReset() // abort config
	}
	if disc.Cancel != nil {
		disc.Cancel() // abort previous session
	}

	// resume advertising node, if we were advertising
	if discovery != "null" && discovery != "" {
		disc.Discovery = discovery
	}
	if disc.Discovery != "" {
		disc.Register(heartbeat, ttl)
		proxy.Follow()
	}

	w.Write([]byte("done"))
}
