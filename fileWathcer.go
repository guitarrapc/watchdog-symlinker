package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
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
func (e *fileWatcher) run(ctx context.Context, exit chan<- error) {

	// initialize
	err := e.initialize()
	if err != nil {
		exit <- err
	}

	// watcher
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create)

	if err := w.Add(e.watchFolder); err != nil {
		log.Fatalln(err)
		exit <- err
	}

	logger.Info("List current files in watchfolder ...")
	for path, f := range w.WatchedFiles() {
		if !f.IsDir() {
			logger.Infof(" * %s: %s\n", path, f.Name())
		}
	}

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
				log.Fatalln(err)
				exit <- err
			case <-w.Closed:
				return
			}
		}
	}()

	logger.Infof("Filewatcher starting ... %s\n", e.watchFolder)
	if err := w.Start(time.Millisecond * 1000); err != nil {
		log.Fatalln(err)
		exit <- err
	}

	logger.Infof("Filewatcher waiting ... \n")
	w.Wait()
	return
}

func (e *fileWatcher) initialize() (err error) {
	// check folder exists
	if !anyfileExists(e.watchFolder) {
		logger.Infof("%s is empty, skip initialize symlink ...\n", e.watchFolder)
		return nil
	}

	// remove exisiting symlink (because re-link to latest log file, existing is waste)
	if symlinkExists(e.dest) {
		logger.Infof("Removing current Symlink: %s\n", e.dest)
		deleteSymlink(e.dest)
	} else {
		logger.Infof("Symlink %s not found ...\n", e.dest)
	}

	// list files
	logger.Infof("Checking latest file ...\n")
	latest, err := getLatestFile(e.watchFolder, e.pattern)
	if err != nil {
		return err
	}
	// map to latest
	if latest.path != "" {
		logger.Infof("Found latest file: %s\n", latest.path)
		makeSymlink(latest.path, e.dest)
	}
	return
}

func replaceSymlink(filePath string, symlinkPath string) {
	deleteSymlink(symlinkPath)
	makeSymlink(filePath, symlinkPath)
}

func deleteSymlink(symlinkPath string) {
	if _, err := os.Lstat(symlinkPath); err == nil {
		if err := os.Remove(symlinkPath); err != nil {
			logger.Infof("failed to unlink: %+v\n", err)
		}
	} else if os.IsNotExist(err) {
		logger.Infof("symlink not found, no need to unlink ...")
	}
}
func makeSymlink(filePath string, symlinkPath string) {
	logger.Infof("link %s with source %s\n", symlinkPath, filePath)
	err := os.Symlink(filePath, symlinkPath)
	if err != nil {
		logger.Infof("Failed to create symlink. %+v\n", err)
	}
}

func anyfileExists(dir string) bool {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false
	}
	return len(files) > 0
}

func getLatestFile(dir string, pattern string) (latest latestFile, err error) {
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func symlinkExists(linkPath string) bool {
	info, err := os.Lstat(linkPath)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}
