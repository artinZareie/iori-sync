package database

import (
	"github.com/artinZareie/iori-sync/internal/filesystem"
	"gorm.io/gorm"
)

type WatchDirectory struct {
	gorm.Model
	Path      string     `gorm:"not null"`
	ReadOnly  bool       `gorm:"default:false;not null"`
	PathRules []PathRule `gorm:"foreignKey:WatchDirectoryID;constraint:OnDelete:CASCADE;"`
}

type WatchDirectoryJSON struct {
	Path        string                `json:"path"`
	ReadOnly    bool                  `json:"read_only"`
	Containings []filesystem.FileJSON `json:"containings,omitempty"`
}
