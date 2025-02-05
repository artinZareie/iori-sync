package main

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/grandcat/zeroconf"
)

func TestServeWithValidPassword(t *testing.T) {
	go serve(8081, "testpass")
	time.Sleep(time.Second)

	resp, err := http.Get("http://localhost:8081/")
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello, world!\n" {
		t.Errorf("unexpected response body: %s", string(body))
	}
}

func TestZeroconfRegistration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		t.Fatalf("Failed to create resolver: %v", err)
	}

	go serve(8082, "testpass")
	time.Sleep(time.Second)

	entries := make(chan *zeroconf.ServiceEntry)
	err = resolver.Browse(ctx, "_http._tcp", "local.", entries)
	if err != nil {
		t.Fatalf("Failed to browse: %v", err)
	}

	found := false
	timeout := time.After(3 * time.Second)

	for {
		select {
		case entry := <-entries:
			if entry.ServiceInstanceName() == "IoriSyncServer._http._tcp.local." {
				found = true
				cancel()
			}
		case <-timeout:
			cancel()
		}
		if found {
			break
		}
	}

	if !found {
		t.Error("Did not find zeroconf service for IoriSyncServer")
	}
}
