package filewatch

import (
	"github.com/guitarrapc/watchdog-symlinker/directory"
	"github.com/guitarrapc/watchdog-symlinker/symlink"
	"github.com/kardianos/service"
)

// Handler keep for filewatcher required parameters
type Handler struct {
	Logger      service.Logger
	FilePattern string
	Directory   string
	SymlinkName string
	Dest        string
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
