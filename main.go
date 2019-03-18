package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
	flag "github.com/spf13/pflag"
)

// version info
var (
	Version  string
	Revision string
)

// global logger
var logger service.Logger

func main() {
	svcConfig := &service.Config{
		Name:        "watchdog-symlinker",
		DisplayName: "watchdog-symlinker",
		Description: "watch directory and create symlink to the latest file.",
	}

	// setup watchdog service
	filewatcher := &fileWatcher{option: fileWatcherOption{}}
	httphealthcheck := &healthcheckhttp{}
	statsdhealthcheck := &healthcheckstatsd{}

	// flags
	command := flag.StringP("command", "c", "", "specify service command. (available list : install|uninstall|start|stop)")
	flag.StringVarP(&filewatcher.option.filePattern, "file", "f", "", "specify file name pattern to watch changes. (regex string)")
	flag.StringVarP(&filewatcher.directoryPattern, "directory", "d", "", "specify full path to watch directory. (regex string)")
	flag.StringVarP(&filewatcher.symlinkName, "symlink", "s", "", "specify symlink name.")
	flag.BoolVar(&filewatcher.option.useFileWalk, "useFileWalk", false, "use walk directory instead of file event.")
	flag.BoolVar(&httphealthcheck.disable, "healthcheckHttpDisabled", false, "disable local http healthcheck.")
	flag.StringVar(&httphealthcheck.addr, "healthcheckHttpAddr", "127.0.0.1:12250", "specify http healthcheck waiting host:port.")
	flag.BoolVar(&statsdhealthcheck.disable, "healthcheckStatsdDisabled", false, "disable datadog statsd healthcheck.")
	flag.StringVar(&statsdhealthcheck.addr, "healthcheckStatsdAddr", "127.0.0.1:8125", "specify statsd healthcheck waiting host:port.")
	flag.Parse()
	if *command == "" && (filewatcher.option.filePattern == "" || filewatcher.directoryPattern == "" || filewatcher.symlinkName == "") {
		log.Println("version", Version)
		log.Println("revision", Revision)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// create service
	w := &watchdog{filewatcher: *filewatcher, healthchecks: []healthcheck{httphealthcheck, statsdhealthcheck}}
	svcConfig.Arguments = []string{
		"-f", filewatcher.option.filePattern,
		"-d", filewatcher.directoryPattern,
		"-s", filewatcher.symlinkName,
		"--healthcheckHttpAddr", httphealthcheck.addr,
		"--healthcheckStatsdAddr", statsdhealthcheck.addr,
	}
	if filewatcher.option.useFileWalk {
		svcConfig.Arguments = append(svcConfig.Arguments, "--useFileWalk")
	}
	if httphealthcheck.disable {
		svcConfig.Arguments = append(svcConfig.Arguments, "--healthcheckHttpDisabled")
	}
	if statsdhealthcheck.disable {
		svcConfig.Arguments = append(svcConfig.Arguments, "--healthcheckStatsdDisabled")
	}
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
