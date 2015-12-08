package service

import (
	proxy "github.com/jeffjen/docker-ambassador/proxy"
	disc "github.com/jeffjen/go-discovery"

	_ "github.com/Sirupsen/logrus"

	"net/http"
	"time"
)

func Configure(w http.ResponseWriter, r *http.Request) {
	if err := common("PUT", r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var (
		heartbeat time.Duration
		ttl       time.Duration
	)

	if hbStr := r.Form.Get("hb"); hbStr != "" {
		hb, err := time.ParseDuration(hbStr)
		if err != nil {
			http.Error(w, "invalid arg", 400)
			return
		}
		heartbeat = hb
	}
	if ttlStr := r.Form.Get("ttl"); ttlStr != "" {
		t, err := time.ParseDuration(ttlStr)
		if err != nil {
			http.Error(w, "invalid arg", 400)
			return
		}
		ttl = t
	}
	if discovery := r.Form.Get("discovery"); discovery != "" {
		// attach provided new discovery endpoint
		disc.Discovery = discovery
	}

	disc.Cancel() // abort previous session

	// start new session
	disc.Register(heartbeat, ttl)

	if proxy.ConfigReset != nil {
		// abort config
		proxy.ConfigReset()
	}

	// reload config setting
	proxy.Follow()

	w.Write([]byte("done"))
}

func Follow(w http.ResponseWriter, r *http.Request) {
	if err := common("PUT", r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var (
		proxycfg = r.Form.Get("key")
	)

	if proxycfg == "" {
		http.Error(w, "Bad Request", 400)
		return
	}

	// Reassin proxy config key
	proxy.ProxyConfigKey = proxycfg

	if proxy.ConfigReset != nil {
		// abort config
		proxy.ConfigReset()
	}

	// reload config setting
	proxy.Follow()

	w.Write([]byte("done"))
}
