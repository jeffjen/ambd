package web

import (
	api "github.com/jeffjen/docker-ambassador/api"
	srv "github.com/jeffjen/docker-ambassador/api/service"
	d "github.com/jeffjen/go-discovery/info"

	log "github.com/Sirupsen/logrus"
)

func init() {
	vmux := &api.VarServeMux{}
	vmux.HandleFunc(`/proxy/(.+)`, srv.ProxyRemove)

	api.GetServeMux().HandleFunc("/info", d.Info)
	api.GetServeMux().HandleFunc("/config", srv.Configure)
	api.GetServeMux().HandleFunc("/proxy", srv.ProxyHelper)
	api.GetServeMux().HandleFunc("/proxy/list", srv.ProxyList)
	api.GetServeMux().HandleFunc("/proxy/app-config", srv.Follow)
	api.GetServeMux().Handle("/proxy/", vmux)
}

func RunAPIEndpoint(addr string, stop chan<- struct{}) {
	defer close(stop)

	server := api.GetServer()

	server.Addr = addr
	log.Error(server.ListenAndServe())
}
