package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"firmware/hal/rfid"
)

func main() {
	scanner := &rfid.RFIDScanner{}
	if err := scanner.Init(); err != nil {
		panic(err)
	}
	defer scanner.Close()

	fmt.Println("RFID scanner running. Press Ctrl+C to stop.")

	// Wait for Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("Stopping scanner...")
}
