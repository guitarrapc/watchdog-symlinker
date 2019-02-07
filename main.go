package main

import (
	"log"
	"os"

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

	// setup watchdog service
	filewatcher := &fileWatcher{}
	httphealthcheck := &healthcheckhttp{}
	statsdhealthcheck := &healthcheckstatsd{}

	// flags
	command := flag.StringP("command", "c", "", "specify service command from install|uninstall|start|stop")
	flag.StringVarP(&filewatcher.pattern, "pattern", "p", "", "specify file name pattern to watch changes")
	flag.StringVarP(&filewatcher.watchFolderPattern, "folder", "f", "", "specify folder name pattern contains file name to watch")
	flag.StringVarP(&filewatcher.symlinkName, "symlink", "s", "", "specify symlink name")
	flag.BoolVar(&httphealthcheck.enable, "healthcheckHttpEnabled", true, "Use local http healthcheck or not.")
	flag.StringVar(&httphealthcheck.addr, "healthcheckHttpAddr", "127.0.0.1:12250", "specify http healthcheck's waiting host:port.")
	flag.BoolVar(&statsdhealthcheck.enable, "healthcheckStatsdEnabled", true, "Use datadog statsd healthcheck or not.")
	flag.StringVar(&statsdhealthcheck.addr, "healthcheckStatsdAddr", "127.0.0.1:8125", "specify statsd healthcheck's waiting host:port.")
	flag.Parse()
	if *command == "" && (filewatcher.pattern == "" || filewatcher.watchFolderPattern == "" || filewatcher.symlinkName == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// create service
	w := &watchdog{filewatcher: *filewatcher, healthchecks: []healthcheck{httphealthcheck, statsdhealthcheck}}
	svcConfig.Arguments = []string{"-p", w.filewatcher.pattern, "-f", w.filewatcher.watchFolderPattern, "-s", w.filewatcher.symlinkName}
	s, err := service.New(w, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	w.service = s

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

	// run in terminal
	err = s.Run()

	// notify Run error
	logger.Infof("Exiting service.")
	if err != nil {
		logger.Error(err)
	}
}
