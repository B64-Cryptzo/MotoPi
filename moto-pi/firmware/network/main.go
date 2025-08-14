package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/B64-Cryptzo/moto-pi-network/monitor"
	"github.com/B64-Cryptzo/moto-pi-network/scan"
)

func cleanup() {
	err := monitor.ResetAllInterfacesToManaged()
	if err != nil {
		log.Printf("Failed to reset interface on cleanup: %v", err)
	}
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	defer cleanup()

	err := monitor.ResetAllInterfacesToManaged()
	if err != nil {
		log.Fatalf("Failed to reset interfaces: %v", err)
	}

	scanner := &scan.RealScanner{Interface: "wlan1"}

	aps, err := scanner.ScanNetworks()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d Access Points", len(aps))
	for _, ap := range aps {
		fmt.Printf("SSID: %s - Strength: %d - MAC Address: %s\n", ap.SSID, ap.SignalStrength, ap.MAC)
	}

	// TODO: start AP hosting, deauth mode, network monitoring...
}
