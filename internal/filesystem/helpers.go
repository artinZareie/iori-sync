package filesystem

import (
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Lists all containing directories and files in a given path. This function
// returns a list of File structs. All returned paths are calculated relatively to the
// base path, for security reasons.
func WalkAsList(basePath string) ([]File, error) {
	basePath = filepath.Clean(basePath)

	folderList := make([]File, 0)

	err := filepath.Walk(basePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var relPath string
		relPath, err = filepath.Rel(basePath, path)

		if err != nil {
			return err
		}

		folderList = append(folderList, File{
			Path: relPath,
			Info: info,
		})
		return nil
	})

	return folderList, err
}

// Lists all containing directories and files in a given path. This function
// returns a list of File structs. All returned paths are calculated relatively to the
// base path, for security reasons. Additionally, this function filters the files
// based on the provided guards.
func WalkAsListGuarded(basePath string, guards []FileGuard) ([]File, error) {
	basePath = filepath.Clean(basePath)
	basePath, _ = filepath.Abs(basePath)

	folderList := make([]File, 0)

	err := filepath.Walk(basePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var relPath string
		relPath, err = filepath.Rel(basePath, path)

		if err != nil {
			return err
		}

		folderList = append(folderList, File{
			Path: relPath,
			Info: info,
		})
		return nil
	})

	return FileGuardFilterFiles(folderList, guards), err
}

// TODO: Add rules to each watch path.
// NOTE: This function is incomplete and requires additional implementation.
func Watch(paths []string, feedback func()) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	for _, path := range paths {
		watcher.Add(path)
	}

	<-make(chan struct{})
}
