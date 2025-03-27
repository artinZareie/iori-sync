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

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, files)
		},
	},
}
