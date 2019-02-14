package main

import (
	"context"
	"os"

	"github.com/kardianos/service"
)

type watchdog struct {
	exit         chan struct{}
	exitError    chan error
	err          error
	service      service.Service
	filewatcher  fileWatcher
	healthchecks []healthcheck
}

func (w *watchdog) run() (err error) {

	// stop service when exiting, because it never return to main
	defer func() {
		if service.Interactive() {
			w.Stop(w.service)
		} else {
			w.service.Stop()
		}
	}()

	// context : goroutine leak prevention
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// background action1. healthcheck
	for _, healthcheck := range w.healthchecks {
		go healthcheck.run(ctx, w.exitError)
	}

	// background action2. filewatcher
	go w.filewatcher.run(ctx, w.exit, w.exitError)

	// monitor exit.
	select {
	case <-w.exit:
		logger.Info("watchdog-symlinker exit called ...")
		return nil
	case err := <-w.exitError:
		logger.Errorf("watchdog-symlinker exit called via error ... %s", err)
		// pass service exit reason
		w.err = err
		return err
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

	if service.Interactive() {
		// exit process
		if w.err != nil {
			logger.Info("exiting watchdog-symlinker with exitcode 1 ...")
			os.Exit(1)
		}
		logger.Info("exiting watchdog-symlinker with exitcode 0 ...")
		os.Exit(0)
	}
	return w.err
}
