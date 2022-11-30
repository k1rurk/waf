package proxy

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"waf/internal/filter"
	"waf/internal/sqli"
)

type Handle struct {
	cache *filter.Cache
}

// ProxyRequestHandler handles the http request using proxy
func (c *Handle) ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		buf, _ := ioutil.ReadAll(r.Body)
		r2 := r.Clone(r.Context())
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		r2.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		err := r2.ParseForm()
		if err != nil {
			log.Printf("Cannot parse form %v", err)
			http.Error(w, "Parsing parameter error", http.StatusInternalServerError)
		}
		for _, arrVal := range r2.Form {
			for _, val := range arrVal {
				for _, f := range sqli.FArray {
					val = f(val)
				}
				log.Println(val)
				for _, fltr := range c.cache.Filter {
					if matched, _ := regexp.MatchString(fltr.Rule, val); matched {
						log.Printf("Rule %s captured vulnerability: %s", fltr.Id, fltr.Description)
						http.Error(w, "The request has been blocked", http.StatusForbidden)
						return
					}
				}
			}
		}

		proxy.ServeHTTP(w, r)
	}
}
