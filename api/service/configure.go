package service

import (
	disc "github.com/jeffjen/docker-ambassador/discovery"
	proxy "github.com/jeffjen/docker-ambassador/proxy"

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
		hbStr  = r.Form.Get("hb")
		ttlStr = r.Form.Get("ttl")

		heartbeat time.Duration
		ttl       time.Duration
	)

	heartbeat, err := time.ParseDuration(hbStr)
	if err != nil {
		http.Error(w, "invalid arg", 400)
		return
	}
	ttl, err = time.ParseDuration(ttlStr)
	if err != nil {
		http.Error(w, "invalid arg", 400)
		return
	}
	if discovery := r.Form.Get("discovery"); discovery != "" {
		// attach provided new discovery endpoint
		disc.Discovery = discovery
	}

	disc.Cancel() // abort previous session

	// start new session
	disc.Register(heartbeat, ttl)

	proxy.Reload() // restart all proxy session

	w.Write([]byte("done"))
}
