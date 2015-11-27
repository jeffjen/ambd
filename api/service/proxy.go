package service

import (
	api "github.com/jeffjen/docker-ambassador/api"
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
		} else {
			log.WithFields(log.Fields{"err": err}).Warning("proxy failed")
		}
	}

	w.Write([]byte("done"))
}

func ProxyRemove(w http.ResponseWriter, r *http.Request, args []string) {
	if err := common("DELETE", r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var From string = args[0]

	if x := proxy.Store.Get(From); x != nil {
		meta := x.(*proxy.Info)
		meta.Cancel()
		proxy.Store.Del(From)
		w.Write([]byte("done"))
	} else {
		http.Error(w, "not found", 404)
	}

	return
}

func ProxyList(w http.ResponseWriter, r *http.Request) {
	if err := common("GET", r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var (
		listing = make([]*proxy.Info, 0)

		enc = json.NewEncoder(api.NewStreamWriter(w))
	)

	proxy.Store.IterateFunc(func(iden string, x interface{}) {
		meta := x.(*proxy.Info)
		listing = append(listing, meta)
	})

	enc.Encode(listing)
	return
}
