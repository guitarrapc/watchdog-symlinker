package main

import (
	"context"
	"time"

	"github.com/kardianos/service"
)

type watchdog struct {
	exit        chan struct{}
	exitError   chan error
	filewatcher fileWatcher
	healthcheck healthcheck
}

func (w *watchdog) run() (err error) {

	// context : goroutine leak prevention
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: http health check は外す
	// TODO: go-datadog sdk で metrics を直接投げる
	// MEMO: jobrunner で60s に一回実行。
	// healthcheck
	go w.healthcheck.run(ctx, w.exitError)

	// filewatcher
	go w.filewatcher.run(ctx, w.exitError)

	// monitor stopped
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case tm := <-ticker.C:
			// TODO: 消す
			logger.Infof("Still running at %v", tm)
		case <-w.exit:
			logger.Info("watchdog-symlinker exit called ...")
			ticker.Stop()
			return nil
		case err := <-w.exitError:
			logger.Errorf("watchdog-symlinker exit called via error ...\n%s", err)
			ticker.Stop()
			return err
		}
	}
}

func (w *watchdog) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal ...")
	} else {
		logger.Info("Running under service manager ...")
	}
	w.exit = make(chan struct{})
	w.exitError = make(chan error)

	go w.run()
	return nil
}

func (w *watchdog) Stop(s service.Service) error {
	logger.Info("stopping watchdog-symlinker ...")
	close(w.exit)
	close(w.exitError)
	logger.Info("successfully stopped ...")
	return nil
}
