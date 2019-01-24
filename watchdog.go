package main

import (
	"context"
	"time"

	"github.com/kardianos/service"
)

type watchdog struct {
	exit         chan struct{}
	exitError    chan error
	filewatcher  fileWatcher
	healthchecks []healthcheck
}

func (w *watchdog) run() (err error) {

	// context : goroutine leak prevention
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// background action1. healthcheck
	for _, healthcheck := range w.healthchecks {
		go healthcheck.run(ctx, w.exitError)
	}

	// background action2. filewatcher
	go w.filewatcher.run(ctx, w.exit, w.exitError)

	// monitor every 1sec.
	// will be exit via standard or error
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
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
