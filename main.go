package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/kardianos/service"
)

var logger service.Logger

func main() {
	svcConfig := &service.Config{
		Name:        "watchdog-symlinker",
		DisplayName: "watchdog-symlinker",
		Description: "watch folder and create symlink to the latest file.",
	}

	// arguments
	var command string
	watchdog := &watchdog{}
	switch len(os.Args) {
	case 2:
		// command line service start/stop/uninstall call
		command = os.Args[1]
	case 4:
		// service/command line invokation
		watchdog.pattern = os.Args[1]
		watchdog.watchFolder = os.Args[2]
		watchdog.symlinkName = os.Args[3]
		watchdog.dest = path.Join(watchdog.watchFolder, watchdog.symlinkName)
	case 5:
		// command line service install called (also start/stop/uninstall can work.)
		command = os.Args[1]
		svcConfig.Arguments = []string{os.Args[2], os.Args[3], os.Args[4]}
	default:
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		logger.Infof(`Usage: %s filenamePattern FolderToWatch symlinkName\nSample: %s ^.*.log$ %s current.log\n`, os.Args[0], os.Args[0], dir)
		os.Exit(1)
	}

	// create service
	s, err := service.New(watchdog, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// setup the logger
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal()
	}

	// service action
	if command != "" {
		err = service.Control(s, os.Args[1])
		if err != nil {
			logger.Warning("Failed (%s) : %s\n", os.Args[1], err)
			return
		}
		logger.Info("Succeeded (%s)\n", os.Args[1])
		return
	}

	// Run in terminal
	s.Run()
}
