package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kardianos/service"
)

type watchdog struct {
	exit        chan error // for Windows Service
	filewatcher fileWatcher
	healthcheck httpHealthcheck
}

func (e *watchdog) run() (err error) {

	// healthcheck
	// TODO: go-datadog sdk で metrics を直接投げる
	// MEMO: jobrunner で60s に一回実行。
	go func() {
		e.exit <- e.healthcheck.run()
	}()

	// filewatcher
	e.exit <- e.filewatcher.initialize()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	e.exit <- e.filewatcher.run()

	// monitor stopped
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case tm := <-ticker.C:
			// TODO: 消す
			logger.Infof("Still running at %v", tm)
		case <-e.exit:
			ticker.Stop()
			logger.Info("watchdog-symlinker Stop ...")
			return nil
		}
	}
}

func (e *watchdog) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	e.exit = make(chan error)

	go e.run()
	return nil
}

func (e *watchdog) Stop(s service.Service) error {
	close(e.exit)
	return nil
}
