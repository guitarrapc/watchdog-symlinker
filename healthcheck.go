package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpHealthcheck struct{}

func (e *httpHealthcheck) run() (err error) {
	// TODO: health check は外す
	// health check
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "health")
	})
	// localhost 以外は拒否
	err = routes.Run("127.0.0.1:8080")
	return
}
