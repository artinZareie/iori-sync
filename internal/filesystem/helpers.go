package filesystem

import (
	"io/fs"
	"path/filepath"
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
