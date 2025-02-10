package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/grandcat/zeroconf"
)

type ServerInfo struct {
	IP         string
	Port       int
	UUID       string
	DeviceName string
}

func getServers(timeout int) ([]ServerInfo, error) {
	resolver, err := zeroconf.NewResolver(nil)

	if err != nil {
		return nil, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	servers := make([]ServerInfo, 0)
	done := make(chan bool)

	go func() {
		for entry := range entries {
			for _, ip := range entry.AddrIPv4 {
				url := fmt.Sprintf("http://%s:%d/who", ip, entry.Port)
				resp, err := http.Get(url)

				if err != nil {
					continue
				}

				data, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				var deviceInfo DeviceInfo
				err = json.Unmarshal(data, &deviceInfo)

				if err != nil {
					continue
				}

				// TODO: Uncomment. Commented because of debugging purposes.
				// if deviceInfo.UUID == cfg.UUID {
				// 	continue
				// }

				servers = append(servers, ServerInfo{
					IP:         ip.String(),
					Port:       entry.Port,
					UUID:       deviceInfo.UUID,
					DeviceName: deviceInfo.DeviceName,
				})
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Second)
	defer cancel()

	err = resolver.Browse(ctx, "_http._tcp", "local.", entries)

	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
	case <-done:
	}

	return servers, nil
}
