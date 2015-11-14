package api

import (
	"io"
	"net/http"
)

var (
	mux *http.ServeMux
	s   *http.Server
)

func GetServeMux() *http.ServeMux {
	return mux
}

func GetServer() *http.Server {
	return s
}

type StreamWriter struct {
	w http.ResponseWriter
}

func (s *StreamWriter) Write(data []byte) (n int, err error) {
	n, err = s.w.Write(data)
	if err == nil {
		s.w.(http.Flusher).Flush()
	}
	return
}

func NewStreamWriter(w http.ResponseWriter) io.Writer {
	if _, ok := w.(http.Flusher); ok {
		return &StreamWriter{w}
	} else {
		return w
	}
}
