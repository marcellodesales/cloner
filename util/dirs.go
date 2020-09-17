package util

import (
	"github.com/disiqueira/gotree"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

/**
 * If the dir is empty or not
 * https://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty/30708914#30708914
 */
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

/**
 * Get the dir tree printed as the tree command, excluding some dirs
 */
func GetDirTree(dirPath string, excludedList []string) (string, error) {
	gotree, err := dirTree(dirPath, excludedList)

	if err != nil {
		return "", err
	}

	return gotree.Print(), nil
}

// Based on https://github.com/FalcoSuessgott/gitget/blob/39b928b65e5c1bdeb0c109ab4e741c2ea86a35f2/tree/tree.go
func dirTree(path string, excludedList []string) (gotree.Tree, error) {
	name := path[strings.LastIndex(path, "/")+1:]
	tree := gotree.New(path)
	err := filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Exclude files/dirs in the provided list
		if len(excludedList) > 0 {
			for _, excludeFile := range excludedList {
				if strings.Contains(dir, excludeFile) {
					return nil
				}
			}
		}

		if name == path[strings.LastIndex(path, "/")+1:] && !info.IsDir() {
			tree.Add(info.Name())
		}

		if info.IsDir() && info.Name() != name {
			tmpTree := buildSubdirectoryTree(dir)
			tree.AddTree(tmpTree)
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tree, nil
}

// Based on https://github.com/FalcoSuessgott/gitget/blob/39b928b65e5c1bdeb0c109ab4e741c2ea86a35f2/tree/tree.go
func buildSubdirectoryTree(dir string) gotree.Tree {
	dirName := dir[strings.LastIndex(dir, "/")+1:]
	tree := gotree.New(dirName)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// if directory, step into and build tree
		if info.IsDir() && dirName != info.Name() {
			tree.AddTree(buildSubdirectoryTree(path))
			return filepath.SkipDir
		}

		// only add nodes to tree with the same depth
		if len(strings.Split(dir, "/"))+1 == len(strings.Split(path, "/")) &&
			info.Name() != dirName && !info.IsDir() {
			tree.Add(info.Name())
		}

		return nil
	})

	if err != nil {
		return nil
	}

	return tree
}