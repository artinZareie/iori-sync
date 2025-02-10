package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func interactiveConfigure() {
	fmt.Printf("Please enter the new device name, or hit enter: ")
	var name string
	fmt.Scanln(&name)

	if name == "" {
		name = cfg.DeviceName
	}

	cfg.DeviceName = name

	fmt.Printf("Please enter the new timeout in seconds, or hit enter: ")
	var timeout int
	fmt.Scanln(&timeout)

	if timeout == 0 {
		timeout = cfg.Timeout
	}

	cfg.Timeout = timeout

	fmt.Printf("Do you want to change password? (y/N)")
	var changePassword string
	fmt.Scanln(&changePassword)

	if changePassword == "Y" || changePassword == "y" {
		fmt.Printf("Please enter the new password: ")
		var password string
		fmt.Scanln(&password)

		cfg.Password = password
	}

	saveConfig(cfg)
}

func interactiveMe() {
	fmt.Printf("UUID       : %s\n", cfg.UUID)
	fmt.Printf("Device Name: %s\n", cfg.DeviceName)
	fmt.Printf("Timeout    : %d\n", cfg.Timeout)
	fmt.Printf("Password   : %s\n", cfg.Password)
}

func interactive() {
	err := obtainLock()

	if err != nil {
		fmt.Println("Another client is already running.")
		os.Exit(0)
	}

	defer releaseLock()

	// Interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		releaseLock()
		os.Exit(0)
	}()

	for {
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)
		cmd := strings.Split(input, " ")[0]
		cmd = strings.ToLower(cmd)

		switch cmd {
		case "exit":
			os.Exit(0)

		case "list-servers":
			listServers(5)

		case "ls":
			listServers(5)

		case "configure":
			interactiveConfigure()

		case "me":
			interactiveMe()

		default:
			fmt.Println("Ambiguous command:", cmd)
		}
	}
}
