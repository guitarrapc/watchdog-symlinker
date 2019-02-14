package filewatch

import (
	"context"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/guitarrapc/watchdog-symlinker/directory"
	"github.com/guitarrapc/watchdog-symlinker/symlink"
	"github.com/kardianos/service"
	"github.com/radovskyb/watcher"
)

// Handler keep for filewatcher required parameters
type Handler struct {
	Logger       service.Logger
	FilePattern  string
	Directory    string
	SymlinkName  string
	Dest         string
	UseFileEvent bool
}

func (e *Handler) initSymlink(directoryPath string, pattern string, dest string) (err error) {
	// check directory exists
	if !directory.ContainsFile(directoryPath) {
		e.Logger.Infof("%s not contains files, skip initialize symlink ...", directoryPath)
		return nil
	}

	// remove exisiting symlink (because re-link to latest log file, existing is waste)
	if symlink.Exists(dest) {
		e.Logger.Infof("Removing current Symlink: %s", dest)
		symlink.Delete(dest)
	} else {
		e.Logger.Infof("Symlink %s not found ...", dest)
	}

	// list files
	e.Logger.Infof("Checking latest file ...")
	latest, err := directory.GetLatestFile(directoryPath, pattern)
	if err != nil {
		return err
	}

	// map to latest
	if latest.Path != "" {
		e.Logger.Infof("Found latest file, source %s link as %s...", latest.Path, dest)
		err = symlink.Create(latest.Path, dest)
		if err != nil {
			e.Logger.Infof("Failed to create symlink. %+v", err)
		}
	}
	return
}

// Run will trigger with directory walk polling
func (e *Handler) Run(ctx context.Context, exit chan<- struct{}, exitError chan<- error) {

	defer e.Logger.Info("exit file non-event Watcher run ...")

	// initialize existing symlink
	err := e.initSymlink(e.Directory, e.FilePattern, e.Dest)
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
	if !directory.IsExists(e.Directory) {
		e.Logger.Info("target directory not found, generating %s...", e.Directory)
		os.MkdirAll(e.Directory, os.ModePerm)
	}

	// add watch directory
	if err = w.Add(e.Directory); err != nil {
		e.Logger.Error(err)
		exitError <- err
		return
	}

	// list watch directory contents
	e.Logger.Info("List watching files ...")
	var fileList []string
	for path, f := range w.WatchedFiles() {
		if !f.IsDir() && f.Name() != e.SymlinkName {
			fileList = append(fileList, path)
		}
	}
	e.Logger.Info(strings.Join(fileList, "\n"))

	// monitor handler
	r := regexp.MustCompile(e.FilePattern)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))
	go func() {
		for {
			select {
			case <-ctx.Done():
				e.Logger.Info("cancel called in filewatcher ...")
				return
			case event := <-w.Event:
				// replace symlink to generated file = latest
				if event.Name() != e.SymlinkName {
					e.Logger.Info(event)
					source := path.Join(e.Directory, event.Name())
					e.Logger.Infof("Create/Replace symlink new: %s, old: %s ...", source, e.Dest)
					err = symlink.Replace(source, e.Dest)
					if err != nil {
						exitError <- err
					}
				}
			case err := <-w.Error:
				e.Logger.Error(err)
				e.Logger.Info("Restarting new filewatcher")
				go e.Run(ctx, exit, exitError)
				return
			case <-w.Closed:
				e.Logger.Info("file watcher ended because of watcher closed ...")
				var complete struct{}
				exit <- complete
				return
			}
		}
	}()

	go w.Wait()

	e.Logger.Infof("successfully start filewatcher %s ...", e.Directory)
	if err := w.Start(time.Second * 1); err != nil {
		e.Logger.Error(err)
		exitError <- err
	}
}
