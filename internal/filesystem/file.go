package filesystem

import "io/fs"

type File struct {
	Path string
	Info fs.FileInfo
}
