package main

import (
	"context"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/guitarrapc/watchdog-symlinker/directory"
	"github.com/guitarrapc/watchdog-symlinker/symlink"
	"github.com/radovskyb/watcher"
)

type fileWatcher struct {
	directoryPattern string
	symlinkName      string
	option           fileWatcherOption
}

type fileWatcherOption struct {
	filePattern string
}

type fileWatchHandler struct {
	filePattern string
	directory   string
	symlinkName string
	dest        string
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
					h := fileWatchHandler{dest: d, filePattern: e.option.filePattern, symlinkName: e.symlinkName, directory: directory}
					go h.run(ctx, exit, exitError)
					found = true
				}
			}

			if found {
				break loop
			}
		}
	}
}

func (e *fileWatchHandler) run(ctx context.Context, exit chan<- struct{}, exitError chan<- error) {

	defer logger.Info("exit fileWatcher mainhandler ...")

	// initialize existing symlink
	err := initSymlink(e.directory, e.filePattern, e.dest)
	if err != nil {
		exitError <- err
		return
	}

	// watcher
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create)
	defer w.Close()

	// generate directory
	if !directory.IsExists(e.directory) {
		logger.Info("target directory not found, generating %s...", e.directory)
		os.MkdirAll(e.directory, os.ModePerm)
	}

	// add watch directory
	if err = w.Add(e.directory); err != nil {
		logger.Error(err)
		exitError <- err
		return
	}

	// list watch directory contents
	logger.Info("List watching files ...")
	var fileList []string
	for path, f := range w.WatchedFiles() {
		if !f.IsDir() && f.Name() != e.symlinkName {
			fileList = append(fileList, path)
		}
	}
	logger.Info(strings.Join(fileList, "\n"))

	// monitor handler
	r := regexp.MustCompile(e.filePattern)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("cancel called in filewatcher ...")
				return
			case event := <-w.Event:
				// replace symlink to generated file = latest
				if event.Name() != e.symlinkName {
					logger.Info(event)
					source := path.Join(e.directory, event.Name())
					logger.Infof("Create/Replace symlink new: %s, old: %s ...", source, e.dest)
					err = symlink.Replace(source, e.dest)
					if err != nil {
						exitError <- err
					}
				}
			case err := <-w.Error:
				logger.Error(err)
				logger.Info("Restarting new filewatcher")
				go e.run(ctx, exit, exitError)
				return
			case <-w.Closed:
				logger.Info("file watcher ended because of watcher closed ...")
				var complete struct{}
				exit <- complete
				return
			}
		}
	}()

	go w.Wait()

	logger.Infof("successfully start filewatcher %s ...", e.directory)
	if err := w.Start(time.Second * 1); err != nil {
		logger.Error(err)
		exitError <- err
	}
}

func initSymlink(directoryPath string, pattern string, dest string) (err error) {
	// check directory exists
	if !directory.ContainsFile(directoryPath) {
		logger.Infof("%s not contains files, skip initialize symlink ...", directoryPath)
		return nil
	}

	// remove exisiting symlink (because re-link to latest log file, existing is waste)
	if symlink.Exists(dest) {
		logger.Infof("Removing current Symlink: %s", dest)
		symlink.Delete(dest)
	} else {
		logger.Infof("Symlink %s not found ...", dest)
	}

	// list files
	logger.Infof("Checking latest file ...")
	latest, err := directory.GetLatestFile(directoryPath, pattern)
	if err != nil {
		return err
	}

	// map to latest
	if latest.Path != "" {
		logger.Infof("Found latest file, source %s link as %s...", latest.Path, dest)
		err = symlink.Create(latest.Path, dest)
		if err != nil {
			logger.Infof("Failed to create symlink. %+v", err)
		}
	}
	return
}
