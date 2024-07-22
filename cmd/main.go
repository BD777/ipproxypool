package main

import (
	"flag"
	"log"

	"github.com/BD777/ipproxypool/pkg/config"
	"github.com/BD777/ipproxypool/pkg/service"
)

func main() {
	cfg := config.DefaultConfig

	// parse flags: mode, host, port
	flag.StringVar(&cfg.Mode, "mode", cfg.Mode, "debug, release")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "")
	flag.Parse()

	log.Printf("service runs with config: %+v", cfg)

	s := service.NewHTTPServer(cfg)
	s.ListenAndServe()
}
