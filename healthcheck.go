package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gin-gonic/gin"
)

type (
	healthcheck interface {
		run(ctx context.Context, exit chan<- error) (err error)
	}
	healthcheckhttp   struct{}
	healthcheckstatsd struct{}
)

func (*healthcheckhttp) run(ctx context.Context, exit chan<- error) (err error) {
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "health")
	})

	select {
	case <-ctx.Done():
		logger.Info("cancel called in healthcheckhttp ...")
		return
	case exit <- routes.Run("127.0.0.1:8080"):
		return
	}
}

func (e *healthcheckstatsd) run(ctx context.Context, exit chan<- error) (err error) {
	c, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
		return err
	}

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	c.Namespace = "watchdog-symlinker."
	c.Tags = append(c.Tags, "watcher:watchdog-symlinker")
	for {
		select {
		case <-ctx.Done():
			logger.Info("cancel called in healthcheckstatsd ...")
			return
		case <-ticker.C:
			logger.Info("sending metrics to datadog")
			err = c.Incr("health", nil, 1)
			if err != nil {
				logger.Errorf("error while sending datadog metrics", err)
			}
		}
	}
}
