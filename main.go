package main

import (
	"fmt"
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

	var command string
	w := &watchdog{}
	w.filewatcher = fileWatcher{}
	w.healthchecks = []healthcheck{&healthcheckhttp{}, &healthcheckstatsd{}}

	// arguments
	switch len(os.Args) {
	case 2:
		// command line "service start/stop/uninstall" action called
		command = os.Args[1]
	case 4:
		// main invokation
		w.filewatcher.pattern = os.Args[1]
		w.filewatcher.watchFolder = os.Args[2]
		w.filewatcher.symlinkName = os.Args[3]
		w.filewatcher.dest = path.Join(w.filewatcher.watchFolder, w.filewatcher.symlinkName)
	case 5:
		// command line "service install" action called (also start/stop/uninstall can work.)
		command = os.Args[1]
		svcConfig.Arguments = []string{os.Args[2], os.Args[3], os.Args[4]}
	default:
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Printf("Usage: %s filenamePattern FolderToWatch symlinkName\nSample: %s ^.*.log$ %s current.log", os.Args[0], os.Args[0], dir)
		os.Exit(1)
	}

	// create service
	s, err := service.New(w, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// service action
	if command != "" {
		err = service.Control(s, os.Args[1])
		if err != nil {
			logger.Warning("Failed (%s): %s\n", os.Args[1], err)
			return
		}
		logger.Info("Succeeded (%s)\n", os.Args[1])
		return
	}

	// setup the logger
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal()
	}

	// Run in terminal
	s.Run()
}
