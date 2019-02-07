package symlink

import "os"

// Replace will replace existing symlink
func Replace(filePath string, symlinkPath string) error {
	return replaceSymlink(filePath, symlinkPath)
}

// Delete will remove symlink
func Delete(symlinkPath string) error {
	return deleteSymlink(symlinkPath)
}

// Exists check symlink is exists or not
func Exists(symlinkPath string) bool {
	return existsSymlink(symlinkPath)
}

// Create will create new symlink
func Create(filePath string, symlinkPath string) error {
	return createSymlink(filePath, symlinkPath)
}

func replaceSymlink(filePath string, symlinkPath string) error {
	deleteSymlink(symlinkPath)
	return createSymlink(filePath, symlinkPath)
}

func deleteSymlink(symlinkPath string) error {
	if existsSymlink(symlinkPath) {
		return os.Remove(symlinkPath)
	}
	return nil
}

func existsSymlink(symlinkPath string) bool {
	info, err := os.Lstat(symlinkPath)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}

func createSymlink(filePath string, symlinkPath string) error {
	err := os.Symlink(filePath, symlinkPath)
	return err
}
