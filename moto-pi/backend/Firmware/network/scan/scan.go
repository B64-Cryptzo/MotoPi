package scan

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// AccessPoint represents one WiFi network found in a scan.
type AccessPoint struct {
	SSID           string
	SignalStrength int    // e.g. RSSI in dBm
	Encryption     string // e.g. WPA2, Open, WEP
	MAC            string // BSSID (MAC address)
}

// NetworkInterface defines the interface for scanning networks.
type NetworkInterface interface {
	ScanNetworks() ([]AccessPoint, error)
}

// StubScanner implements NetworkInterface and returns mock data.
type StubScanner struct{}

func (s *StubScanner) ScanNetworks() ([]AccessPoint, error) {
	aps := []AccessPoint{
		{SSID: "HomeWiFi", SignalStrength: -40, Encryption: "WPA2", MAC: "00:11:22:33:44:55"},
		{SSID: "CafeNet", SignalStrength: -70, Encryption: "Open", MAC: "66:77:88:99:AA:BB"},
	}

	if len(aps) == 0 {
		return nil, errors.New("no access points found")
	}
	return aps, nil
}

// RealScanner implements NetworkInterface and uses nmcli to scan networks.
type RealScanner struct {
	Interface string // e.g. "wlan0"
}

// ensureMonitorMode switches the interface to monitor mode if it’s not already.
func (r *RealScanner) ensureMonitorMode() error {
	mode, err := r.getInterfaceMode()
	if err != nil {
		return err
	}

	if mode == "Monitor" {
		return nil // Already in monitor mode
	}

	// Switch to monitor mode
	cmds := [][]string{
		{"ip", "link", "set", r.Interface, "down"},
		{"iw", r.Interface, "set", "monitor", "control"},
		{"ip", "link", "set", r.Interface, "up"},
	}

	for _, cmdArgs := range cmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to run %v: %v - output: %s", cmdArgs, err, string(out))
		}
	}

	return nil
}

// getInterfaceMode returns the current mode of the wireless interface.
func (r *RealScanner) getInterfaceMode() (string, error) {
	cmd := exec.Command("iwconfig", r.Interface)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	outputStr := string(out)
	// Look for "Mode:Monitor" in the output
	for _, line := range strings.Split(outputStr, "\n") {
		if strings.Contains(line, "Mode:Monitor") {
			return "Monitor", nil
		}
	}

	// Not found, assume something else
	return "Managed", nil
}

func (r *RealScanner) ScanNetworks() ([]AccessPoint, error) {
	if err := r.ensureMonitorMode(); err != nil {
		return nil, err
	}

	cmd := exec.Command("iw", r.Interface, "scan")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseIwScanOutput(out), nil
}

// parseIwScanOutput parses output from `iw <iface> scan` into AccessPoint list
func parseIwScanOutput(output []byte) []AccessPoint {
	var aps []AccessPoint
	var bssMacRegex = regexp.MustCompile(`BSS ([0-9a-fA-F:]{17})`)

	scanner := bufio.NewScanner(bytes.NewReader(output))
	var currentAP AccessPoint

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "BSS ") {
			// If previous AP block has SSID, save it
			if currentAP.SSID != "" {
				aps = append(aps, currentAP)
				currentAP = AccessPoint{}
			}

			// Use regex to extract MAC address
			matches := bssMacRegex.FindStringSubmatch(line)
			if len(matches) == 2 {
				currentAP.MAC = matches[1]
			}
		}

		if strings.HasPrefix(line, "SSID:") {
			currentAP.SSID = strings.TrimSpace(strings.TrimPrefix(line, "SSID:"))
		}

		if strings.HasPrefix(line, "signal:") {
			// Example: signal: -40.00 dBm
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if val, err := strconv.ParseFloat(parts[1], 32); err == nil {
					currentAP.SignalStrength = int(val)
				}
			}
		}

		if strings.HasPrefix(line, "RSN:") || strings.HasPrefix(line, "WPA:") {
			// Presence of these indicates WPA/WPA2 encryption
			currentAP.Encryption = "WPA/WPA2"
		} else if strings.HasPrefix(line, "Privacy:") {
			// Older format indication for encryption
			currentAP.Encryption = "Encrypted"
		} else if strings.HasPrefix(line, "Capability:") && !strings.Contains(line, "Privacy") {
			currentAP.Encryption = "Open"
		}
	}

	// Append last AP if exists
	if currentAP.SSID != "" {
		aps = append(aps, currentAP)
	}

	return aps
}
