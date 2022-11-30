package proxy

import (
	"net/http/httputil"
	"net/url"
)

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(u), nil
}
