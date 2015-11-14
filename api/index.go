package api

import (
	srv "github.com/jeffjen/docker-ambassador/api/service"

	"net/http"
)

func init() {
	mux = http.NewServeMux()
	s = &http.Server{Handler: mux}

	mux.HandleFunc("/info", srv.Info)
	mux.HandleFunc("/config", srv.Configure)
}
