package service

import (
	"fmt"
	"net/http"

	"github.com/BD777/ipproxypool/pkg/config"
)

func NewHTTPServer(cfg *config.Config) *http.Server {
	router := NewRouter(cfg)

	return &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}
}
