package api

import (
	"net/http"
	"regexp"
)

type VarHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, args []string)
}

type VarHandlerFunc func(http.ResponseWriter, *http.Request, []string)

func (n VarHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, args []string) {
	n(w, r, args)
}

type route struct {
	pattern *regexp.Regexp
	handler VarHandler
}

type VarServeMux struct {
	routes []*route
}

func (v *VarServeMux) Handle(pattern string, handler VarHandler) {
	v.routes = append(v.routes, &route{
		regexp.MustCompile(pattern),
		handler.(VarHandler),
	})
}

func (v *VarServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request, []string)) {
	v.routes = append(v.routes, &route{
		regexp.MustCompile(pattern),
		VarHandlerFunc(handler),
	})
}

func (v *VarServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range v.routes {
		matched := route.pattern.FindStringSubmatch(r.URL.Path)
		if len(matched) > 0 {
			route.handler.ServeHTTP(w, r, matched[1:])
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
