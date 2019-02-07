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
	"github.com/radovskyb/watcher"
)

type fileWatcher struct {
	pattern            string
	watchFolderPattern string
	watchFolder        string
	symlinkName        string
	dest               string
}

// runWatcher
// @summary: file watcher to replace symlink to latest
func (e *fileWatcher) run(ctx context.Context, exit chan<- struct{}, exitError chan<- error) {

	logger.Info("starting filewatcher ...")

	// extract base path
	logger.Infof("extract base path for %s ...", e.watchFolderPattern)
	basePath, err := directory.GetBasePath(e.watchFolderPattern)
	if err != nil {
		exitError <- err
		return
	}

	// loop until target folder found
	var directories []string
	pattern := regexp.MustCompile(e.watchFolderPattern)
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	found := false
L:
	for {
		select {
		case <-t.C:
			logger.Infof("walking directories in %s ...", basePath)
			directories, err = directory.Dirwalk(basePath)
			directories = append(directories, basePath)
			if err != nil {
				logger.Error(err)
				logger.Info("retrying to find target folder check ...")
				break
			}

			// check each directory
			logger.Infof("matching directories %s ...", directories)
			for _, directory := range directories {
				isMatch := pattern.MatchString(directory)
				logger.Infof("result %s: %s ...", strconv.FormatBool(isMatch), directory)
				if isMatch {
					d := path.Join(directory, e.symlinkName)
					logger.Infof("start checking %s ...", d)
					_e := fileWatcher{dest: d, pattern: e.pattern, symlinkName: e.symlinkName, watchFolder: directory}
					go _e.mainHandler(ctx, exit, exitError)
					found = true
				}
			}

			if found {
				break L
			}
		}
	}
}

func (e *fileWatcher) mainHandler(ctx context.Context, exit chan<- struct{}, exitError chan<- error) {

	defer logger.Info("exit fileWatcher mainhandler ...")

	// initialize existing symlink
	err := initSymlink(e.watchFolder, e.pattern, e.dest)
	if err != nil {
		exitError <- err
		return
	}

	// watcher
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create)
	defer w.Close()

	// add watch folder
	if err := w.Add(e.watchFolder); err != nil {
		logger.Error(err)
		exitError <- err
		return
	}

	logger.Info("List watching files ...")
	var fileList []string
	for path, f := range w.WatchedFiles() {
		if !f.IsDir() && f.Name() != e.symlinkName {
			fileList = append(fileList, path)
		}
	}
	logger.Info(strings.Join(fileList, "\n"))

	r := regexp.MustCompile(e.pattern)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("cancel called in filewatcher ...")
				return
			case event := <-w.Event:
				if event.Name() != e.symlinkName {
					logger.Info(event)
					source := path.Join(e.watchFolder, event.Name())
					replaceSymlink(source, e.dest)
				}
			case err := <-w.Error:
				logger.Error(err)
				exitError <- err
			case <-w.Closed:
				logger.Info("file watcher ended because of watcher closed.")
				var complete struct{}
				exit <- complete
				return
			}
		}
	}()

	go w.Wait()

	logger.Infof("successfully start filewatcher %s ...", e.watchFolder)
	if err := w.Start(time.Second * 1); err != nil {
		logger.Error(err)
		exitError <- err
	}
}

func initSymlink(folderPath string, pattern string, dest string) (err error) {
	// check folder exists
	if !directory.ContainsFile(folderPath) {
		logger.Infof("%s not contains files, skip initialize symlink ...\n", folderPath)
		return nil
	}

	// remove exisiting symlink (because re-link to latest log file, existing is waste)
	if existsSymlink(dest) {
		logger.Infof("Removing current Symlink: %s\n", dest)
		deleteSymlink(dest)
	} else {
		logger.Infof("Symlink %s not found ...\n", dest)
	}

	// list files
	logger.Infof("Checking latest file ...\n")
	latest, err := directory.GetLatestFile(folderPath, pattern)
	if err != nil {
		return err
	}
	// map to latest
	if latest.Path != "" {
		logger.Infof("Found latest file: %s\n", latest.Path)
		createSymlink(latest.Path, dest)
	}
	return
}

// symlink
func replaceSymlink(filePath string, symlinkPath string) {
	deleteSymlink(symlinkPath)
	createSymlink(filePath, symlinkPath)
}

func deleteSymlink(symlinkPath string) {
	if _, err := os.Lstat(symlinkPath); err == nil {
		if err = os.Remove(symlinkPath); err != nil {
			logger.Infof("failed to unlink: %+v\n", err)
		}
	} else if os.IsNotExist(err) {
		logger.Infof("symlink not found, no need to unlink ...")
	}
}

func existsSymlink(linkPath string) bool {
	info, err := os.Lstat(linkPath)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}

func createSymlink(filePath string, symlinkPath string) {
	logger.Infof("link %s with source %s\n", symlinkPath, filePath)
	err := os.Symlink(filePath, symlinkPath)
	if err != nil {
		logger.Infof("Failed to create symlink. %+v\n", err)
	}
}
