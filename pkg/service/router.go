package service

import (
	"net/http"

	"github.com/BD777/ipproxypool/pkg/config"
	"github.com/BD777/ipproxypool/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(cfg.Mode)

	{
		r.GET("/", ListIPProxy)
		r.GET("/count", GetIPProxyCount)
		r.POST("/delete", DeleteIPProxy)
	}

	return r
}

func ListIPProxy(c *gin.Context) {
	req := &models.ListIPProxyRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		logrus.Errorf("ShouldBindQuery error: %v", err)
		c.JSON(http.StatusBadRequest, "invalid args")
		return
	}

	resp := models.ListIPProxy(req)
	if resp == nil {
		resp = []*models.IPProxy{}
	}
	c.JSON(http.StatusOK, resp)
}

func GetIPProxyCount(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetIPProxyCount())
}

func DeleteIPProxy(c *gin.Context) {
	req := &models.DeleteIPProxyRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		logrus.Errorf("ShouldBindQuery error: %v", err)
		c.JSON(http.StatusBadRequest, "invalid args")
		return
	}

	err := models.DeleteIPProxy(req)
	if err != nil {
		logrus.Errorf("delete ip error: %v with req %#v", err, req)
		c.JSON(http.StatusBadRequest, "delete fail")
		return
	}

	c.JSON(http.StatusOK, map[string]struct{}{})
}
