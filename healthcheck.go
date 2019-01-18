package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpHealthcheck struct{}

func (e *httpHealthcheck) run(ctx context.Context, exit chan<- error) (err error) {
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "health")
	})

	select {
	case <-ctx.Done():
		logger.Info("cancel called in httpHealthchec ...")
		return
	case exit <- routes.Run("127.0.0.1:8080"):
		return
	}
}
