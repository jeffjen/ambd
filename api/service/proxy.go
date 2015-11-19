package service

import (
	proxy "github.com/jeffjen/docker-ambassador/proxy"

	log "github.com/Sirupsen/logrus"

	"encoding/json"
	"net/http"
)

func ProxyHelper(w http.ResponseWriter, r *http.Request) {
	if err := common("POST", r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var (
		meta = new(proxy.Info)

		dec = json.NewDecoder(r.Body)
	)

	if err := dec.Decode(meta); err != nil {
		log.WithFields(log.Fields{"err": err}).Warning("bad proxy spec")
		http.Error(w, "bad proxy spec", 400)
		return
	}

	if err := proxy.Listen(meta); err != nil {
		if err != proxy.ErrProxyExist {
			log.WithFields(log.Fields{"err": err}).Warning("proxy failed")
			http.Error(w, "internal server error", 500)
			return
		}
	}

	w.Write([]byte("done"))
}
