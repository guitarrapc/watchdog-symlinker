package main

import (
	"context"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gin-gonic/gin"
)

type (
	healthcheck interface {
		run(ctx context.Context, exitError chan<- error) (err error)
	}
	healthcheckhttp struct {
		enable bool
		addr   string
	}
	healthcheckstatsd struct {
		enable bool
		addr   string
	}
)

func (h *healthcheckhttp) run(ctx context.Context, exitError chan<- error) (err error) {
	// validate
	if !h.enable {
		logger.Info("healthcheckhttp is disabled ...")
		return nil
	}

	logger.Infof("starting healthcheckhttp on %s ...", h.addr)

	// execute
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "health")
	})

	logger.Info("successfully start healthcheckhttp ... ")

	err = routes.Run(h.addr)
	select {
	case <-ctx.Done():
		logger.Info("cancel called in healthcheckhttp ...")
		return
	case exitError <- err:
		return
	}
}

func (h *healthcheckstatsd) run(ctx context.Context, exitError chan<- error) (err error) {
	// validate
	if !h.enable {
		logger.Info("healthcheckstatsd is disabled ...")
		return nil
	}

	logger.Infof("starting healthcheckstatsd on %s ...", h.addr)

	// execute
	c, err := statsd.New(h.addr)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("successfully start healthcheckstatsd ... ")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	c.Namespace = "watchdog_symlinker."
	c.Tags = append(c.Tags, "watcher:watchdog_symlinker")
	metricName := "health"
	for {
		select {
		case <-ctx.Done():
			logger.Info("cancel called in healthcheckstatsd ...")
			return
		case <-ticker.C:
			logger.Infof("sending metrics to datadog: %s%s", c.Namespace, metricName)
			err = c.Incr(metricName, nil, 1)
			if err != nil {
				logger.Errorf("error while sending datadog metrics, keep runing healtchcheck. %s", err)
			}
		}
	}
}
