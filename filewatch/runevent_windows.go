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

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)
	defer close(c)

	// Set up a watchpoint listening on events
	// Dispatch each create events separately to c.
	if err := notify.Watch(e.Directory, c, notify.FileNotifyChangeFileName, notify.FileNotifyChangeCreation); err != nil {
		e.Logger.Error(err)
		exitError <- err
		return
	}
	defer notify.Stop(c)

	r := regexp.MustCompile(e.FilePattern)

	// monitor handler
	var current os.FileInfo
	for {
		select {
		case <-ctx.Done():
			e.Logger.Info("cancel called in filewatcher ...")
			return
		case event := <-c:
			source := event.Path()
			fileName := filepath.Base(source)
			if !r.MatchString(fileName) {
				return
			}
			switch event.Event() {
			case notify.FileActionAdded:
				e.Logger.Info(event)
				fi, err := os.Stat(source)
				if err != nil {
					e.Logger.Errorf("error happen when checking %s. %s", source, err)
					return
				}
				// replace symlink to generated file = latest
				if fileName != e.SymlinkName {
					e.Logger.Info(event)
					e.Logger.Infof("Create/Replace symlink new: %s, old: %s ...", source, e.Dest)
					err = symlink.Replace(source, e.Dest)
					if err != nil {
						exitError <- err
					}
					current = fi
				}
			case notify.FileActionRenamedNewName:
				e.Logger.Info(event)
				fi, err := os.Stat(event.Path())
				if err != nil {
					e.Logger.Errorf("error happen when checking %s. %s", event.Path(), err)
					return
				}
				if current == nil || fi.ModTime().After(current.ModTime()) {
					// replace symlink to renamed file
					if fileName != e.SymlinkName {
						e.Logger.Info(event)
						e.Logger.Infof("Create/Replace symlink new: %s, old: %s ...", source, e.Dest)
						err = symlink.Replace(source, e.Dest)
						if err != nil {
							exitError <- err
						}
						current = fi
					}
				}
			}
		}
	}
}
