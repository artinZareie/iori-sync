package main

import (
	"encoding/json"
	"log"
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

	result := db.Where(Device{UUID: device.UUID}).Assign(Device{Name: device.Name}).FirstOrCreate(&device)

	if result.Error != nil {
		log.Printf("Database error: %v\n", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
