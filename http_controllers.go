package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWho(c *gin.Context) {
	deviceInfo := DeviceInfo{
		UUID:       cfg.UUID,
		DeviceName: cfg.DeviceName,
	}

	c.JSON(http.StatusOK, deviceInfo)
}

func HandleRegister(c *gin.Context) {
	if !CheckPassword(c.PostForm("password")) {
		c.Status(http.StatusUnauthorized)
		return
	}

	device := Device{
		UUID: c.PostForm("uuid"),
		Name: c.PostForm("name"),
	}

	if db == nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	result := db.Where(Device{UUID: device.UUID}).
		Assign(Device{Name: device.Name}).FirstOrCreate(&device)

	if result.Error != nil {
		log.Printf("Database error: %v\n", result.Error)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
