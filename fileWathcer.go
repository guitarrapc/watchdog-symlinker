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
	filePattern string
	useFileWalk bool
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
	var dirs []string
	pattern := regexp.MustCompile(e.directoryPattern)
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	found := false
loop:
	for {
		select {
		case <-t.C:
			// child only
			logger.Infof("walking directories in %s ...", basePath)
			dirs, err = directory.Dirwalk(basePath)
			if err != nil {
				logger.Error(err)
				logger.Info("retrying to find target directory check ...")
				break
			}
			// TODO: bad knowhow, should fix.
			if runtime.GOOS == "windows" && len(dirs) == 0 {
				logger.Infof("no child directories found in %s, add basePath to monitoring target ...", basePath)
				dirs = append(dirs, basePath)
			}

			// check each directory
			logger.Infof("matching directories with pattern %s ...", pattern.String())
			for _, directory := range dirs {
				dir := directory
				isMatch := pattern.MatchString(dir)
				logger.Infof("(%s) %s", strconv.FormatBool(isMatch), dir)
				if isMatch {
					d := path.Join(dir, e.symlinkName)
					logger.Infof("start checking %s ...", d)
					h := filewatch.Handler{
						Dest:         d,
						FilePattern:  e.option.filePattern,
						SymlinkName:  e.symlinkName,
						Directory:    dir,
						Logger:       logger,
						UseFileEvent: false,
					}
					if !e.option.useFileWalk && (runtime.GOOS == "windows" || runtime.GOOS == "linux") {
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
