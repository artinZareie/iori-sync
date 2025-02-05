package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
	"gopkg.in/yaml.v2"
)

type Config struct {
	UUID       string `yaml:"uuid"`
	DeviceName string `yaml:"device_name"`
}

const configFilePath = "config.yaml"

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

func getDeviceName() (string, error) {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func serve(port int, password string) {
	if password == "" {
		fmt.Println("Error: password is required")
		os.Exit(1)
	}

	server, err := zeroconf.Register("IoriSyncServer", "_http._tcp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	defer server.Shutdown()

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	}))

	fmt.Printf("Starting server on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func main() {
	cfg := loadConfig()
	fmt.Println("===================================")
	fmt.Printf("UUID       : %s\n", cfg.UUID)
	fmt.Printf("Device Name: %s\n", cfg.DeviceName)
	fmt.Println("===================================")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	port := serveCmd.Int("port", 8080, "Port to run the server on")
	password := serveCmd.String("password", "", "Password for authentication (required)")

	if len(os.Args) < 2 {
		fmt.Println("Please use -h to see available commands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		serveCmd.Parse(os.Args[2:])
		serve(*port, *password)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
