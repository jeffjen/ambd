package api

import (
	srv "github.com/jeffjen/docker-ambassador/api/service"
	d "github.com/jeffjen/go-discovery/info"

	log "github.com/Sirupsen/logrus"

	"net/http"
)

func init() {
	mux = http.NewServeMux()
	s = &http.Server{Handler: mux}

	vmux := &VarServeMux{}
	vmux.HandleFunc(`/proxy/(\w+)`, srv.ProxyRemove)

	mux.HandleFunc("/info", d.Info)
	mux.HandleFunc("/config", srv.Configure)
	mux.HandleFunc("/proxy", srv.ProxyHelper)
	mux.Handle("/proxy/", vmux)
}

func RunAPIEndpoint(addr string, stop chan<- struct{}) {
	defer close(stop)

	server := GetServer()

	server.Addr = addr
	log.Error(server.ListenAndServe())
}
