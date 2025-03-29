package server

import (
	"log"

	"github.com/artinZareie/iori-sync/internal/database"
	"github.com/artinZareie/iori-sync/internal/filesystem"
	"github.com/gin-gonic/gin"
)

var routeTable = map[string]RouteConfig{
	"/": {
		Methods: []string{"GET"},
		Handler: func(c *gin.Context) {
			// Show all directories in the db.
			db := database.GetDB()
			var watchDirs []database.WatchDirectory

			if err := db.Preload("PathRules").Find(&watchDirs).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to retrieve watch directories"})
				return
			}

			response := make([]database.WatchDirectoryJSON, len(watchDirs))

			for x := range watchDirs {
				guards := make([]filesystem.FileGuard, len(watchDirs[x].PathRules))
				for i, rule := range watchDirs[x].PathRules {
					guards[i] = rule.ToFileGuard()
				}

				files, err := filesystem.WalkAsListGuarded(watchDirs[x].Path, guards)

				if err != nil {
					c.JSON(500, gin.H{"error": "Failed to walk the directory"})
					log.Fatal("Error while walking the directory:", err)
					return
				}

				filesJSON := make([]filesystem.FileJSON, len(files))

				for i, file := range files {
					filesJSON[i] = file.ToFileJSON()
				}

				response = append(response, database.WatchDirectoryJSON{
					Path:        watchDirs[x].Path,
					ReadOnly:    watchDirs[x].ReadOnly,
					Containings: filesJSON,
				})
			}

			c.JSON(200, response)
		},
	},
}
