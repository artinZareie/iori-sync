package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/grandcat/zeroconf"
)

func serve(port int, password string) {
	if password != "" {
		cfg.Password = password
		saveConfig(cfg)
	} else if cfg.Password == "" {
		fmt.Println("Error: password is required")
		os.Exit(1)
	}

	server, err := zeroconf.Register("IoriSyncServer",
		"_http._tcp", "local.",
		port,
		[]string{"txtv=0", "lo=1", "la=2"},
		nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	defer server.Shutdown()

	// Initialize Gin router
	router := gin.Default()

	// Route handlers
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/who")
	})

	router.GET("/who", func(c *gin.Context) {
		HandleWho(c)
	})

	router.POST("/register", func(c *gin.Context) {
		HandleRegister(c)
	})

	fmt.Printf("Starting server on port %d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}
