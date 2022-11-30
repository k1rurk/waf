package main

import (
	"log"
	"waf/internal/config"
	"waf/internal/filter"
	"waf/internal/proxy"
)

func main() {
	cfg := config.ReadConfigFile()
	cache := filter.ReadFilterFile(cfg.Filename)

	prx, err := proxy.NewProxy(cfg.Remote)
	if err != nil {
		log.Fatalln(err)
	}
	mux := proxy.NewMux(cache)

	server := proxy.NewServer(cfg.Bind, mux, prx)
	log.Printf("Listening on %s, forwarding to %s", cfg.Bind, cfg.Remote)
	server.StartServer()
}
