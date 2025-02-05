package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"github.com/grandcat/zeroconf"
)

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
