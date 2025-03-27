package filesystem

import (
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WalkAsList(path string) ([]File, error) {
	folderList := make([]File, 0)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		folderList = append(folderList, File{
			Path: path,
			Info: info,
		})
		return nil
	})

	return folderList, err
}

func WalkAsListGuarded(path string, guards []FileGuard) ([]File, error) {
	folderList := make([]File, 0)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		folderList = append(folderList, File{
			Path: path,
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
