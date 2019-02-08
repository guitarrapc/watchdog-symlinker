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
	"github.com/radovskyb/watcher"
)

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
