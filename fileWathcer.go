package main

import (
	"context"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/guitarrapc/watchdog-symlinker/directory"
	"github.com/guitarrapc/watchdog-symlinker/filewatch"
)

type fileWatcher struct {
	directoryPattern string
	symlinkName      string
	option           fileWatcherOption
}

type fileWatcherOption struct {
	filePattern  string
	useFileEvent bool
}

// runWatcher
// @summary: file watcher to replace symlink to latest
func (e *fileWatcher) run(ctx context.Context, exit chan<- struct{}, exitError chan<- error) {

	logger.Info("starting filewatcher ...")

	// extract base path
	logger.Infof("extract base path for %s ...", e.directoryPattern)
	basePath, err := directory.GetBasePath(e.directoryPattern)
	if err != nil {
		exitError <- err
		return
	}

	// loop until target directory found
	var directories []string
	pattern := regexp.MustCompile(e.directoryPattern)
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	found := false
loop:
	for {
		select {
		case <-t.C:
			logger.Infof("walking directories in %s ...", basePath)
			directories, err = directory.Dirwalk(basePath)
			if err != nil {
				logger.Error(err)
				logger.Info("retrying to find target directory check ...")
				break
			}
			directories = append(directories, basePath)

			// check each directory
			logger.Infof("matching directories with pattern %s ...", pattern.String())
			for _, directory := range directories {
				isMatch := pattern.MatchString(directory)
				logger.Infof("(%s) %s", strconv.FormatBool(isMatch), directory)
				if isMatch {
					d := path.Join(directory, e.symlinkName)
					logger.Infof("start checking %s ...", d)
					h := filewatch.Handler{Dest: d, FilePattern: e.option.filePattern, SymlinkName: e.symlinkName, Directory: directory, Logger: logger}
					if e.option.useFileEvent && runtime.GOOS == "windows" {
						go h.RunEvent(ctx, exit, exitError)
					} else {
						go h.Run(ctx, exit, exitError)
					}
					found = true
				}
			}

			if found {
				break loop
			}
		}
	}
}
