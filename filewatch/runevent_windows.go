// +build windows

package filewatch

import (
	"context"
	"os"
	"path/filepath"
	"regexp"

	"github.com/guitarrapc/watchdog-symlinker/directory"

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

	create := make(chan notify.EventInfo, 1)
	defer close(create)

	// Set up a watchpoint listening on events
	// Dispatch each create events separately to channel.
	if err := notify.Watch(e.Directory, c, notify.FileNotifyChangeFileName, notify.FileNotifyChangeCreation); err != nil {
		e.Logger.Error(err)
		exitError <- err
		return
	}
	defer notify.Stop(c)
	if err = notify.Watch(e.Directory, create, notify.Create); err != nil {
		e.Logger.Error(err)
		exitError <- err
		return
	}
	defer notify.Stop(create)

	r := regexp.MustCompile(e.FilePattern)

	// monitor handler
	var current os.FileInfo
	for {
		select {
		case <-ctx.Done():
			e.Logger.Info("cancel called in filewatcher ...")
			return
		case event := <-create:
			e.Logger.Infof("create event detected: %s", event)
			if directory.IsExists(event.Path()) {
				e.Logger.Info("event was directory, skip and wait next ...")
				break
			}
			source := event.Path()
			fileName := filepath.Base(source)
			if !r.MatchString(fileName) {
				break
			}
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
		case event := <-c:
			e.Logger.Infof("file event detected: %s", event)
			if !e.UseFileEvent {
				return
			}
			source := event.Path()
			fileName := filepath.Base(source)
			if !r.MatchString(fileName) {
				e.Logger.Infof("event filename was not target, skip and wait next ...")
				break
			}
			switch event.Event() {
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
