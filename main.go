package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/radovskyb/watcher"
)

type latestFile struct {
	path    string
	modTime time.Time
}

var dest string

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, `Usage: %s filenamePattern FolderToWatch symlinkName
Sample: %s ^.*.log$ "C:/Users/ikiru.yoshizaki/Downloads/logfiles" current.log
`, os.Args[0], os.Args[0])
		os.Exit(1)
	}

	pattern := os.Args[1]
	watchFolder := os.Args[2]
	symlinkName := os.Args[3]

	dest = path.Join(watchFolder, symlinkName)

	// initialize
	err := initialize(watchFolder, symlinkName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: Windows Service として実行

	// TODO: go-datadog sdk で metrics を直接投げる
	// MEMO: jobrunner で60s に一回実行。

	// TODO: health check は外す
	// health check
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "health")
	})
	go func() {
		// localhost 以外は拒否
		routes.Run("127.0.0.1:8080")
	}()

	// file watcher
	runWatcher(pattern, watchFolder, symlinkName)
}

func initialize(watchFolder string, symlinkName string) (err error) {
	// check folder exists
	if !anyfileExists(watchFolder) {
		fmt.Printf("%s is empty, skip initialize symlink.\n", watchFolder)
		return nil
	}

	target := path.Join(watchFolder, symlinkName)
	// check symlink already exists
	if symlinkExists(target) {
		fmt.Printf("Removing current Symlink %s.\n", target)
		deleteSymlink(target)
	} else {
		fmt.Printf("Symlink %s not found.\n", target)
	}

	// map to latest file
	fmt.Println("Checking latest file.")
	latest, err := getLatestFile(watchFolder)
	if err != nil {
		return err
	}
	if latest.path != "" {
		fmt.Printf("Found latest file %s, start try link.\n", latest.path)
		makeSymlink(latest.path, dest)
	}
	return
}

// runWatcher
// @summary: file watcher to replace symlink to latest
func runWatcher(pattern string, watchFolder string, symlinkName string) {
	// watcher
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create)
	r := regexp.MustCompile(pattern)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))
	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.Name() != symlinkName {
					fmt.Println(event)
					source := path.Join(watchFolder, event.Name())
					// replace symlink
					replaceSymlink(source, dest)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Add(watchFolder); err != nil {
		log.Fatalln(err)
	}

	for path, f := range w.WatchedFiles() {
		if !f.IsDir() {
			fmt.Printf("%s: %s\n", path, f.Name())
		}
	}

	go func() {
		fmt.Printf("Filewatcher started: %s\n", watchFolder)
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
	return len(files) == 0
}

func getLatestFile(dir string) (latest latestFile, err error) {
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
	for _, fi := range files {
		if fi.Mode().IsRegular() {
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

func symlinkExists(filename string) bool {
	if !fileExists(filename) {
		return false
	}

	info, err := os.Lstat(filename)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}
