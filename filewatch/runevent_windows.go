// +build windows

package filewatch

import (
	"context"
	"os"
	"path/filepath"
	"regexp"

	"github.com/guitarrapc/watchdog-symlinker/symlink"
	"github.com/rjeczalik/notify"
)

// RunEvent will trigger with file event
func (e *Handler) RunEvent(ctx context.Context, exit chan<- struct{}, exitError chan<- error) {

	defer e.Logger.Info("exit file event Watcher runEvent ...")

	// initialize existing symlink
	err := e.initSymlink(e.Directory, e.FilePattern, e.Dest)
	if err != nil {
		exitError <- err
		return
	}

	// Set up a watchpoint listening on events
	// Dispatch each create events separately to channel.
	e.Logger.Infof("create watcher for %s ...", e.Directory)
	fileCreate := make(chan notify.EventInfo, 1)
	if err := notify.Watch(e.Directory, fileCreate, notify.FileNotifyChangeFileName); err != nil {
		e.Logger.Error(err)
		exitError <- err
		return
	}
	defer func() {
		notify.Stop(fileCreate)
		close(fileCreate)
	}()

	r := regexp.MustCompile(e.FilePattern)

	// monitor handler
	var current os.FileInfo
	for {
		select {
		case <-ctx.Done():
			e.Logger.Info("cancel called in filewatcher ...")
			return
		case ei := <-fileCreate:
			e.Logger.Info("file event %s", ei)
			source := ei.Path()
			fileName := filepath.Base(source)
			if !r.MatchString(fileName) {
				e.Logger.Infof("event filename was not target, skip and wait next ...")
				break
			}
			switch ei.Event() {
			case notify.FileActionAdded:
				e.Logger.Info("file action added event, checking file exists ...")
				fi, err := os.Stat(source)
				if err != nil {
					e.Logger.Errorf("error happen when checking %s. %s", source, err)
					break
				}
				if fileName == e.SymlinkName {
					e.Logger.Infof("event filename was same as symlink, skip and wait next ...")
					break
				}
				// replace symlink to generated file = latest
				e.Logger.Infof("Create/Replace symlink new: %s, old: %s ...", source, e.Dest)
				err = symlink.Replace(source, e.Dest)
				if err != nil {
					e.Logger.Errorf("error happen when replacing symlink %s. %s", source, err)
					break
				}
				current = fi
			case notify.FileActionRenamedNewName:
				e.Logger.Info("file action added event, checking file exists ...")
				fi, err := os.Stat(source)
				if err != nil {
					e.Logger.Errorf("error happen when checking %s. %s", source, err)
					break
				}
				if fileName == e.SymlinkName {
					e.Logger.Infof("event filename was same as symlink, skip and wait next ...")
					break
				}
				// replace symlink to renamed file
				if current == nil || fi.ModTime().After(current.ModTime()) {
					e.Logger.Infof("Create/Replace symlink new: %s, old: %s ...", source, e.Dest)
					err = symlink.Replace(source, e.Dest)
					if err != nil {
						e.Logger.Errorf("error happen when replacing symlink %s. %s", source, err)
						break
					}
					current = fi
				}
			}
		}
	}
}
