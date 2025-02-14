package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/grandcat/zeroconf"
)

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
