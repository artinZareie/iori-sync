package server

import (
	"github.com/artinZareie/iori-sync/internal/filesystem"
	"github.com/gin-gonic/gin"
)

var routeTable = map[string]RouteConfig{
	"/": {
		Methods: []string{"GET"},
		Handler: func(c *gin.Context) {
			files, err := filesystem.WalkAsList(".")
			var filesJson []filesystem.FileJSON = make([]filesystem.FileJSON, 0)

			for _, file := range files {
				filesJson = append(filesJson, file.ToFileJSON())
			}

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, filesJson)
		},
	},
}
