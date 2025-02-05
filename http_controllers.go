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
