package main

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

type latestFile struct {
	path    string
	modTime time.Time
}

type fileWatcher struct {
	pattern     string
	watchFolder string
	symlinkName string
	dest        string
	latest      latestFile
}

// runWatcher
// @summary: file watcher to replace symlink to latest
func (e *fileWatcher) run(ctx context.Context, exit chan<- struct{}, exitError chan<- error) {

	logger.Info("starting filewatcher ...")

	// initialize existing symlink
	err := e.initializeSymlink()
	if err != nil {
		exitError <- err
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
	}

	logger.Info("List current files in watchfolder ...")
	var fileList []string
	for path, f := range w.WatchedFiles() {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
	}
	logger.Info(strings.Join(fileList, "\n"))
	fileList = nil

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
					e.replaceSymlink(source, e.dest)
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

	go func() {
		w.Wait()
	}()

	logger.Infof("successfully start filewatcher ... %s\n", e.watchFolder)
	if err := w.Start(time.Second * 1); err != nil {
		logger.Error(err)
		exitError <- err
	}
	return
}

func (e *fileWatcher) initializeSymlink() (err error) {
	// check folder exists
	if !e.containsFile(e.watchFolder) {
		logger.Infof("%s is empty, skip initialize symlink ...\n", e.watchFolder)
		return nil
	}

	// remove exisiting symlink (because re-link to latest log file, existing is waste)
	if e.existsSymlink(e.dest) {
		logger.Infof("Removing current Symlink: %s\n", e.dest)
		e.deleteSymlink(e.dest)
	} else {
		logger.Infof("Symlink %s not found ...\n", e.dest)
	}

	// list files
	logger.Infof("Checking latest file ...\n")
	latest, err := e.getLatestFile(e.watchFolder, e.pattern)
	if err != nil {
		return err
	}
	// map to latest
	if latest.path != "" {
		logger.Infof("Found latest file: %s\n", latest.path)
		e.createSymlink(latest.path, e.dest)
	}
	return
}

// symlink
func (e *fileWatcher) replaceSymlink(filePath string, symlinkPath string) {
	e.deleteSymlink(symlinkPath)
	e.createSymlink(filePath, symlinkPath)
}

func (e *fileWatcher) deleteSymlink(symlinkPath string) {
	if _, err := os.Lstat(symlinkPath); err == nil {
		if err := os.Remove(symlinkPath); err != nil {
			logger.Infof("failed to unlink: %+v\n", err)
		}
	} else if os.IsNotExist(err) {
		logger.Infof("symlink not found, no need to unlink ...")
	}
}

func (e *fileWatcher) existsSymlink(linkPath string) bool {
	info, err := os.Lstat(linkPath)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}

func (e *fileWatcher) createSymlink(filePath string, symlinkPath string) {
	logger.Infof("link %s with source %s\n", symlinkPath, filePath)
	err := os.Symlink(filePath, symlinkPath)
	if err != nil {
		logger.Infof("Failed to create symlink. %+v\n", err)
	}
}

// directory
func (e *fileWatcher) containsFile(dir string) bool {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false
	}
	return len(files) > 0
}

func (e *fileWatcher) getLatestFile(dir string, pattern string) (latest latestFile, err error) {
	latest.modTime = time.Time{}
	latest.path = ""
	err = nil

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return latest, err
	}
	if len(files) == 0 {
		return latest, logger.Errorf("no file exists ...")
	}
	r := regexp.MustCompile(pattern)
	for _, fi := range files {
		if fi.Mode().IsRegular() && r.MatchString(fi.Name()) {
			if fi.ModTime().After(latest.modTime) {
				latest.modTime = fi.ModTime()
				latest.path = path.Join(dir, fi.Name())
			}
		}
	}
	if latest.path != "" {
		return latest, nil
	}
	return
}

// file
func (e *fileWatcher) existsFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
