package service

import (
	"github.com/BD777/ipproxypool/pkg/config"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(cfg.Mode)

	// r.Use(mw.InitMiddleware())

	{
		// r.GET("", )
	}

	// apiGroup := r.Group("/api")

	return r
}
