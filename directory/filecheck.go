package directory

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"time"
)

// LatestFile express pathinfo
type LatestFile struct {
	Path    string
	ModTime time.Time
}

// GetLatestFile return latest file in the path
func GetLatestFile(dir string, pattern string) (latest LatestFile, err error) {
	return getLatestFile(dir, pattern)
}

// ContainsFile check directory where it contains file or not
func ContainsFile(dir string) bool {
	return containsFile(dir)
}

func containsFile(dir string) bool {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, file := range files {
		if !file.IsDir() {
			return true
		}
	}
	return false
}

func getLatestFile(dir string, pattern string) (latest LatestFile, err error) {
	latest.ModTime = time.Time{}
	latest.Path = ""
	err = nil

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return latest, err
	}
	if len(files) == 0 {
		return latest, fmt.Errorf("no file exists in %s ... ", dir)
	}
	r := regexp.MustCompile(pattern)
	for _, fi := range files {
		if fi.Mode().IsRegular() && r.MatchString(fi.Name()) {
			if fi.ModTime().After(latest.ModTime) {
				latest.ModTime = fi.ModTime()
				latest.Path = path.Join(dir, fi.Name())
			}
		}
	}
	if latest.Path != "" {
		return latest, nil
	}
	return
}
