package database

import (
	"github.com/artinZareie/iori-sync/internal/filesystem"
	"gorm.io/gorm"
)

type PathRule struct {
	gorm.Model
	Regex            string `gorm:"not null"`
	IsExclude        bool   `gorm:"default:false;not null"`
	IsName           bool   `gorm:"default:false;not null"`
	WatchDirectoryID uint   `gorm:"not null"`
}

func (pr *PathRule) ToFileGuard() filesystem.FileGuard {
	if pr.IsName {
		if pr.IsExclude {
			return &filesystem.RegexNameRejectorGuard{
				Regex: pr.Regex,
			}
		} else {
			return &filesystem.RegexNameAcceptorGuard{
				Regex: pr.Regex,
			}
		}
	} else {
		if pr.IsExclude {
			return &filesystem.RegexPathRejectorGuard{
				Regex: pr.Regex,
			}
		} else {
			return &filesystem.RegexPathAcceptorGuard{
				Regex: pr.Regex,
			}
		}
	}
}
