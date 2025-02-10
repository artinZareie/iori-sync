package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func HangleWho(w http.ResponseWriter, r *http.Request) {
	deviceInfo := DeviceInfo{
		UUID:       cfg.UUID,
		DeviceName: cfg.DeviceName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deviceInfo)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if !CheckPassword(r.FormValue("password")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	device := Device{
		UUID: r.FormValue("uuid"),
		Name: r.FormValue("name"),
	}

	if db == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var existingDevice Device
	result := db.Where("uuid = ?", device.UUID).First(&existingDevice)

	if result != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if err := db.Create(&device).Error; err != nil {
				log.Printf("Error creating device: %v\n", err)
			}
		} else if result.Error != nil {
			log.Printf("Database error: %v\n", result.Error)
		}
	} else {
		if err := db.Model(&existingDevice).Update("name", device.Name).Error; err != nil {
			log.Printf("Error updating device: %v\n", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}
