package main

import (
	"log"
	"os"
	"path"

	"github.com/kardianos/service"
	flag "github.com/spf13/pflag"
)

// global logger
var logger service.Logger

func main() {
	svcConfig := &service.Config{
		Name:        "watchdog-symlinker",
		DisplayName: "watchdog-symlinker",
		Description: "watch folder and create symlink to the latest file.",
	}

	w := &watchdog{}
	w.filewatcher = fileWatcher{}
	w.healthchecks = []healthcheck{&healthcheckhttp{}, &healthcheckstatsd{}}

	// flags
	command := flag.StringP("command", "c", "", "specify service command from install|uninstall|start|stop")
	flag.StringVarP(&w.filewatcher.pattern, "pattern", "p", "", "specify file name pattern to watch changes")
	flag.StringVarP(&w.filewatcher.watchFolder, "folder", "f", "", "specify path to the file watcher's target folder")
	flag.StringVarP(&w.filewatcher.symlinkName, "symlink", "s", "", "specify symlink name")
	flag.Parse()
	if w.filewatcher.watchFolder != "" && w.filewatcher.symlinkName != "" {
		w.filewatcher.dest = path.Join(w.filewatcher.watchFolder, w.filewatcher.symlinkName)
	}
	if *command == "" && (w.filewatcher.pattern == "" || w.filewatcher.watchFolder == "" || w.filewatcher.symlinkName == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// create service
	svcConfig.Arguments = []string{"-p", w.filewatcher.pattern, "-f", w.filewatcher.watchFolder, "-s", w.filewatcher.symlinkName}
	s, err := service.New(w, svcConfig)
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
	if *command != "" {
		err = service.Control(s, *command)
		if err != nil {
			logger.Errorf("Failed (%s): %s\n", *command, err)
			return
		}
		logger.Infof("Succeeded (%s)\n", *command)
		return
	}

	// Run in terminal
	s.Run()
}
