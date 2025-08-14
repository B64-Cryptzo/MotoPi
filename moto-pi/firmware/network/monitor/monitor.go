package monitor

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ResetAllInterfacesToManaged sets all wireless interfaces out of monitor mode into managed mode.
func ResetAllInterfacesToManaged() error {
	ifaces, err := getWirelessInterfaces()
	if err != nil {
		return err
	}

	for _, iface := range ifaces {
		mode, err := getInterfaceMode(iface)
		if err != nil {
			return fmt.Errorf("failed to get mode for %s: %w", iface, err)
		}

		if mode == "Monitor" {
			if err := setInterfaceModeManaged(iface); err != nil {
				return fmt.Errorf("failed to reset %s to managed mode: %w", iface, err)
			}
			fmt.Printf("Set interface %s to managed mode\n", iface)
		}
	}

	return nil
}

func getWirelessInterfaces() ([]string, error) {
	cmd := exec.Command("iw", "dev")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ifaces []string
	lines := bytes.Split(out, []byte{'\n'})
	for _, line := range lines {
		lineStr := strings.TrimSpace(string(line))
		if strings.HasPrefix(lineStr, "Interface") {
			parts := strings.Fields(lineStr)
			if len(parts) == 2 {
				ifaces = append(ifaces, parts[1])
			}
		}
	}

	return ifaces, nil
}

func getInterfaceMode(iface string) (string, error) {
	cmd := exec.Command("iwconfig", iface)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	outputStr := string(out)
	if strings.Contains(outputStr, "Mode:Monitor") {
		return "Monitor", nil
	}

	return "Managed", nil
}

func setInterfaceModeManaged(iface string) error {
	for i := 0; i < 3; i++ {
		err := exec.Command("ip", "link", "set", iface, "down").Run()
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		err = exec.Command("iw", iface, "set", "type", "managed").Run()
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		err = exec.Command("ip", "link", "set", iface, "up").Run()
		if err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return errors.New("failed to set interface to managed mode after retries")
}
