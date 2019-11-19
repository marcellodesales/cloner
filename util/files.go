package util

import (
	"os"
	"path/filepath"
)

/**
 * Delete a directory with contents
 * https://stackoverflow.com/questions/33450980/how-to-remove-all-contents-of-a-directory-using-golang/33451503#33451503
 */
func DeleteDir(dirPath string) error {
	d, err := os.Open(dirPath)

	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)

	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dirPath, name))
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * Delete file
 */
func DeleteFile(filePath string) error {
	var err = os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}