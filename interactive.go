package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/fsnotify/fsnotify"
)

type Command struct {
	Name string
	Abbr string
	Help string
	Func func()
}

var commands = []Command{
	{
		Name: "configure",
		Abbr: "config",
		Help: "Configure the device",
		Func: interactiveConfigure,
	},
	{
		Name: "me",
		Abbr: "m",
		Help: "Show current configuration",
		Func: interactiveMe,
	},
	{
		Name: "list-servers",
		Abbr: "l",
		Help: "List available servers",
		Func: interactiveListServers,
	},
	{
		Name: "exit",
		Abbr: "e",
		Help: "Exit the client",
		Func: interactiveExit,
	},
	{
		Name: "connect",
		Abbr: "c",
		Help: "Connect to a server",
		Func: interactiveConnect,
	},
	{
		Name: "test",
		Abbr: "t",
		Help: "For debugging purposes",
		Func: interactiveTest,
	},
}

func interactiveListServers() {
	listServers(5)
}

func interactiveExit() {
	releaseLock()
	os.Exit(0)
}

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

func interactiveHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Command\tShortcut\tDescription")
	fmt.Fprintln(w, "-------\t--------\t-----------")
	for _, c := range commands {
		fmt.Fprintf(w, "%s\t%s\t%s\n", c.Name, c.Abbr, c.Help)
	}
	w.Flush()
}

func interactiveConnect() {
	fmt.Println("Searching for servers...")
	fmt.Println("==========================================================================================")

	servers, err := getServers(cfg.Timeout)

	if err != nil {
		fmt.Println("Error getting servers:", err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "Row\tIP\tPort\tUUID\tDevice Name")
	fmt.Fprintln(w, "---\t--\t----\t----\t-----------")
	for i, s := range servers {
		fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\n", i+1, s.IP, s.Port, s.UUID, s.DeviceName)
	}
	w.Flush()

	fmt.Println("==========================================================================================")

	for {
		fmt.Printf("\nPlease enter the number of the server you want to connect to: ")
		var serverNumber int
		fmt.Scanln(&serverNumber)

		if serverNumber < 1 || serverNumber > len(servers) {
			fmt.Println("Invalid server number.")
			continue
		}

		server := servers[serverNumber-1]
		fmt.Printf("Connecting to %s:%d...\n", server.IP, server.Port)

		url := fmt.Sprintf("http://%s:%d/register", server.IP, server.Port)

		formData := fmt.Sprintf("uuid=%s&name=%s&password=%s", cfg.UUID, cfg.DeviceName, cfg.Password)
		resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(formData))

		if err != nil {
			fmt.Println("Error connecting to server: ", err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error in registration response: ", resp.Status)
			return
		}

		fmt.Println("Connected to server.")
		defer resp.Body.Close()
		return
	}
}

func interactiveTest() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
		os.Exit(1)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Println("Event:", event)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("Error:", err)
			}
		}
	}()

	watcher.Add("./test")

	if err != nil {
		log.Fatal("Error watching file:", err)
		os.Exit(1)
	}

	<-make(chan struct{})
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
		interactiveExit()
	}()

	for {
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)
		cmd := strings.Split(input, " ")[0]
		cmd = strings.ToLower(cmd)

		executed := false

		for _, c := range commands {
			if c.Name == cmd || c.Abbr == cmd {
				executed = true
				c.Func()
			}
		}

		if !executed && cmd != "help" {
			fmt.Println("Unknown command, please use \"help\" to see available commands.")
		}

		if cmd == "help" || cmd == "h" {
			interactiveHelp()
		}
	}
}
