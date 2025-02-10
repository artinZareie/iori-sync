package main

import (
	"fmt"
	"os"
)

func obtainLock() error {
	_, err := os.Stat("client.lock")
	if os.IsNotExist(err) {
		return os.WriteFile("client.lock", []byte("lock"), 0644)
	}

	return os.ErrExist
}

func releaseLock() {
	err := os.Remove("client.lock")

	if err != nil {
		fmt.Println("Error removing lock file:", err)
	}
}
