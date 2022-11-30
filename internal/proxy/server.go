package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	bind  string
	mux   *Mux
	proxy *httputil.ReverseProxy
}

func NewServer(bind string, mux *Mux, proxy *httputil.ReverseProxy) *Server {
	return &Server{
		bind:  bind,
		mux:   mux,
		proxy: proxy,
	}
}

func (srv *Server) StartServer() {
	rt := srv.mux.SetRoutes(srv.proxy)

	if err := http.ListenAndServe(srv.bind, rt); err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
