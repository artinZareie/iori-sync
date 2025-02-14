package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getDeviceName() (string, error) {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func main() {
	cfg = loadConfig()
	initDB()
	fmt.Println("===================================")
	fmt.Printf("UUID       : %s\n", cfg.UUID)
	fmt.Printf("Device Name: %s\n", cfg.DeviceName)
	fmt.Println("===================================")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	port := serveCmd.Int("port", 8080, "Port to run the server on")
	password := serveCmd.String("password", "", "Password for authentication (required)")

	listServerCmd := flag.NewFlagSet("list-servers", flag.ExitOnError)

	timeout := listServerCmd.Int("timeout", cfg.Timeout, "Timeout in seconds for server discovery")

	interactiveCmd := flag.NewFlagSet("interactive", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Please use -h to see available commands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		serveCmd.Parse(os.Args[2:])
		serve(*port, *password)

	case "list-servers":
		listServerCmd.Parse(os.Args[2:])
		if *timeout != 0 && *timeout != cfg.Timeout {
			cfg.Timeout = *timeout
			saveConfig(cfg)
		}

		listServers(*timeout)

	case "interactive":
		interactiveCmd.Parse(os.Args[2:])
		interactive()

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
