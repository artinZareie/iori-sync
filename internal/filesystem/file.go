package filesystem

import (
	"io/fs"
	"time"
)

type File struct {
	Path string
	Info fs.FileInfo
}

type FileJSON struct {
	Path    string      `json:"path"`
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	IsDir   bool        `json:"is_dir"`
	Mode    fs.FileMode `json:"mode"`
	ModTime time.Time   `json:"mod_time"`
}

func (f File) ToFileJSON() FileJSON {
	return FileJSON{
		Path:    f.Path,
		Name:    f.Info.Name(),
		Size:    f.Info.Size(),
		IsDir:   f.Info.IsDir(),
		Mode:    f.Info.Mode(),
		ModTime: f.Info.ModTime(),
	}
}
