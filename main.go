package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grandcat/zeroconf"
)

func getDeviceName() (string, error) {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func serve(port int, password string) {
	if password != "" {
		cfg.Password = password
		saveConfig(cfg)
	} else if cfg.Password == "" {
		fmt.Println("Error: password is required")
		os.Exit(1)
	}

	server, err := zeroconf.Register("IoriSyncServer",
		"_http._tcp", "local.",
		port,
		[]string{"txtv=0", "lo=1", "la=2"},
		nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	defer server.Shutdown()

	// Initialize Gin router
	router := gin.Default()

	// Route handlers
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/who")
	})

	router.GET("/who", func(c *gin.Context) {
		HandleWho(c)
	})

	router.POST("/register", func(c *gin.Context) {
		HandleRegister(c)
	})

	fmt.Printf("Starting server on port %d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}

func listServers(timeout int) {
	resolver, err := zeroconf.NewResolver(nil)

	if err != nil {
		fmt.Println("Error creating resolver:", err)
		os.Exit(1)
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		defer w.Flush()
		fmt.Fprintln(w, "IP\tPort\tUUID\tDevice Name")

		for entry := range entries {
			if entry.ServiceRecord.Instance == "IoriSyncServer" {
				for _, ip := range entry.AddrIPv4 {
					url := fmt.Sprintf("http://%s:%d/who", ip, entry.Port)
					resp, err := http.Get(url)

					if err != nil {
						fmt.Printf("Error fetching /who from %s: %v\n", url, err)
						continue
					}

					data, _ := io.ReadAll(resp.Body)
					resp.Body.Close()

					var deviceInfo DeviceInfo
					err = json.Unmarshal(data, &deviceInfo)

					// TODO: UNCOMMENT
					// Commented due to debugging.
					/*
						if deviceInfo.UUID == cfg.UUID {
							continue // Skip self
						}
					*/

					if err != nil {
						fmt.Printf("Error decoding /who response from %s: %v\n", url, err)
						continue
					}

					fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", ip, entry.Port, deviceInfo.UUID, deviceInfo.DeviceName)
					w.Flush()
				}
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	err = resolver.Browse(ctx, "_http._tcp", "local.", entries)

	if err != nil {
		fmt.Println("Error browsing for servers:", err)
		os.Exit(1)
	}

	<-ctx.Done()
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
