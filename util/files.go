package util

import (
	"os"
)

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