package server

import "github.com/gin-gonic/gin"

type RouteConfig struct {
	Methods []string
	Handler func(*gin.Context)
}
