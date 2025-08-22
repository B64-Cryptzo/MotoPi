package rfid

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"firmware/hal"
	"firmware/hal/gpio"
)

// Configuration constants
const (
	PM3Client      = "/home/paniq/proxmark3/client/proxmark3"
	PM3Port        = "/dev/ttyACM0"
	TargetString   = "enzogenovese.com"
	ScanInterval   = 10 * time.Millisecond
	SnippetPadding = 10
)

// RFIDScanner implements hal.Device
type RFIDScanner struct {
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	running    bool
}

// Ensure RFIDScanner implements hal.Device
var _ hal.Device = (*RFIDScanner)(nil)

// Init starts the scanning routine
func (r *RFIDScanner) Init() error {
	if r.running {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFunc = cancel
	r.running = true

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				r.scanOnce()
				time.Sleep(ScanInterval)
			}
		}
	}()

	return nil
}

// Close stops the scanning routine
func (r *RFIDScanner) Close() error {
	if !r.running {
		return nil
	}
	r.cancelFunc()
	r.wg.Wait()
	r.running = false
	return nil
}

// Info returns the scanner status
func (r *RFIDScanner) Info() string {
	if r.running {
		return "RFID scanner running"
	}
	return "RFID scanner stopped"
}

// journalLog logs a message with the gimo-events tag
func journalLog(msg string) {
	cmd := exec.Command("logger", "-t", "gimo-events", msg)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to write to journal:", err)
	}
}

// scanOnce performs a single UID/memory check
func (r *RFIDScanner) scanOnce() {
	mem := readTagMemory()
	snippet := extractASCIISnippet(mem, []byte(TargetString), SnippetPadding)
	validRFIDTag := snippet != nil && strings.Contains(string(snippet), "enzogenovese.com")
	if validRFIDTag {
		journalLog("[FOUND_VALID_RFID]")
	}

	gpio.MomentarySwitch(validRFIDTag)
}

func readTagMemory() []byte {
	cmd := exec.Command(PM3Client, PM3Port, "-c", "hf 15 rdmulti -* -b 3 --cnt 6")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil
	}
	re := regexp.MustCompile(`([0-9A-Fa-f]{2})`)
	matches := re.FindAllString(out.String(), -1)
	memory := make([]byte, len(matches))
	for i, b := range matches {
		fmt.Sscanf(b, "%02X", &memory[i])
	}
	return memory
}

func extractASCIISnippet(memory []byte, target []byte, padding int) []byte {
	idx := bytes.Index(memory, target)
	if idx == -1 {
		return nil
	}
	start := idx - padding
	if start < 0 {
		start = 0
	}
	end := idx + len(target) + padding
	if end > len(memory) {
		end = len(memory)
	}
	snippet := memory[start:end]
	printable := make([]byte, 0, len(snippet))
	for _, b := range snippet {
		if b >= 32 && b <= 126 {
			printable = append(printable, b)
		}
	}
	return printable
}
