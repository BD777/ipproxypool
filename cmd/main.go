package main

import (
	"flag"

	"github.com/BD777/ipproxypool/pkg/config"
	"github.com/BD777/ipproxypool/pkg/crawler"
	"github.com/BD777/ipproxypool/pkg/service"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.DefaultConfig

	// parse flags: mode, host, port
	flag.StringVar(&cfg.Mode, "mode", cfg.Mode, "debug, release")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "")
	flag.Int64Var(&cfg.CountThreshold, "count", cfg.CountThreshold, "")
	flag.Parse()

	logrus.Infof("service runs with config: %+v", cfg)

	go crawler.Run(cfg)

	s := service.NewHTTPServer(cfg)
	s.ListenAndServe()
}
