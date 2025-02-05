package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

const configFilePath = "config.yaml"

var cfg Config

type Config struct {
	UUID       string `yaml:"uuid"`
	DeviceName string `yaml:"device_name"`
	Timeout    int    `yaml:"timeout"`
}

type DeviceInfo struct {
	UUID       string `json:"uuid"`
	DeviceName string `json:"device_name"`
}

func loadConfig() Config {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		defaultUUID := uuid.New().String()
		defaultDeviceName, err := getDeviceName()

		if err != nil {
			fmt.Println("Error getting device name:", err)
			os.Exit(1)
		}

		defaultConfig := Config{
			UUID:       defaultUUID,
			DeviceName: defaultDeviceName,
			Timeout:    60,
		}

		saveConfig(defaultConfig)
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		os.Exit(1)
	}

	return config
}

func saveConfig(config Config) {
	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(&config)
	if err != nil {
		fmt.Println("Error encoding config file:", err)
		os.Exit(1)
	}
}
