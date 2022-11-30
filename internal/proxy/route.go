package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"waf/internal/filter"
)

type Mux struct {
	mux     *http.ServeMux
	handler *Handle
}

func NewMux(cache *filter.Cache) *Mux {
	return &Mux{
		handler: &Handle{cache: cache},
		mux:     http.NewServeMux(),
	}
}

func (m *Mux) SetRoutes(proxy *httputil.ReverseProxy) *Logger {
	m.mux.HandleFunc("/", m.handler.ProxyRequestHandler(proxy))
	return NewLogger(m.mux)
}

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr + " " + r.Method + " " + r.URL.String() + " " + r.Proto + " " + r.UserAgent())
	l.handler.ServeHTTP(w, r)
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}
