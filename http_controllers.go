package main

import (
	"encoding/json"
	"net/http"
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

	if err := db.Where("uuid = ?", device.UUID).First(&Device{}).Error; err != nil {
		db.Create(&device)
	} else {
		db.Model(&Device{}).Where("uuid = ?", device.UUID).Update("name", device.Name)
	}

	w.WriteHeader(http.StatusOK)
}
