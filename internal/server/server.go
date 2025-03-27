package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartServer(port int) {
	router := gin.Default()

	for route, config := range routeTable {
		for _, method := range config.Methods {
			router.Handle(method, route, config.Handler)
		}
	}

	router.Run(fmt.Sprintf(":%d", port))
}
