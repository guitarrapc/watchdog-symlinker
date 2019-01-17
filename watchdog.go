package main

import (
	"fmt"
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

type watchdog struct {
	pattern     string
	watchFolder string
	symlinkName string
	dest        string
	latest      latestFile
}

func (e *watchdog) initialize() (err error) {
	// check folder exists
	if !anyfileExists(e.watchFolder) {
		fmt.Printf("%s is empty, skip initialize symlink.\n", e.watchFolder)
		return nil
	}

	// remove exisiting symlink (because re-link to latest log file, existing is waste)
	if symlinkExists(e.dest) {
		fmt.Printf("Removing current Symlink %s.\n", e.dest)
		deleteSymlink(e.dest)
	} else {
		fmt.Printf("Symlink %s not found.\n", e.dest)
	}

	// list files
	fmt.Println("Checking latest file.")
	latest, err := getLatestFile(e.watchFolder, e.pattern)
	if err != nil {
		return err
	}
	// map to latest
	if latest.path != "" {
		fmt.Printf("Found latest file %s, start try link.\n", latest.path)
		makeSymlink(latest.path, e.dest)
	}
	return
}

// runWatcher
// @summary: file watcher to replace symlink to latest
func (e *watchdog) runWatcher() {
	// watcher
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create)
	r := regexp.MustCompile(e.pattern)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))
	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.Name() != e.symlinkName {
					fmt.Println(event)
					source := path.Join(e.watchFolder, event.Name())
					// replace symlink
					replaceSymlink(source, e.dest)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Add(e.watchFolder); err != nil {
		log.Fatalln(err)
	}

	for path, f := range w.WatchedFiles() {
		if !f.IsDir() {
			fmt.Printf("%s: %s\n", path, f.Name())
		}
	}

	go func() {
		fmt.Printf("Filewatcher started: %s\n", e.watchFolder)
		w.Wait()
	}()

	if err := w.Start(time.Millisecond * 1000); err != nil {
		log.Fatalln(err)
	}
}

func replaceSymlink(filePath string, symlinkPath string) {
	deleteSymlink(symlinkPath)
	makeSymlink(filePath, symlinkPath)
}

func deleteSymlink(symlinkPath string) {
	if _, err := os.Lstat(symlinkPath); err == nil {
		if err := os.Remove(symlinkPath); err != nil {
			fmt.Printf("failed to unlink: %+v\n", err)
		}
	} else if os.IsNotExist(err) {
		fmt.Printf("symlink not found, no need to unlink.")
	}
}
func makeSymlink(filePath string, symlinkPath string) {
	fmt.Printf("try link %s to %s\n", filePath, symlinkPath)
	err := os.Symlink(filePath, symlinkPath)
	if err != nil {
		fmt.Printf("Failed to create symlink. %+v\n", err)
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
		return latest, fmt.Errorf("no file exists")
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
