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

	"github.com/B64-Cryptzo/MotoPi/backend/Firmware/hal"
	"github.com/B64-Cryptzo/MotoPi/backend/Firmware/hal/gpio"
)

// Configuration constants
const (
	PM3Client      = "proxmark3"
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
	r.running = false // set false until first successful scan

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := r.safeScanOnce(); err != nil {
					fmt.Println("RFIDScanner encountered error:", err)
					r.running = false
					return
				}
				r.running = true // mark running if scanOnce succeeded
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
		return "online"
	}
	return "offline"
}

// journalLog logs a message with the gimo-events tag
func journalLog(msg string) {
	cmd := exec.Command("logger", "-t", "gimo-events", msg)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to write to journal:", err)
	}
}

// safeScanOnce wraps scanOnce and returns any error encountered
func (r *RFIDScanner) safeScanOnce() (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic in scanOnce: %v", rec)
		}
	}()
	err = r.scanOnce()
	return err
}

// scanOnce performs a single UID/memory check
func (r *RFIDScanner) scanOnce() error {
	mem, err := readTagMemory()
	if err != nil {
		return err
	}

	if mem == nil {
		return nil
	}

	snippet := extractASCIISnippet(mem, []byte(TargetString), SnippetPadding)
	if snippet == nil {
		return nil
	}

	validRFIDTag := strings.Contains(string(snippet), "enzogenovese.com")
	if validRFIDTag {
		journalLog("[FOUND_VALID_RFID]")
		gpio.MomentarySwitch()
	}

	return nil
}

func readTagMemory() ([]byte, error) {
	cmd := exec.Command(PM3Client, PM3Port, "-c", "hf 15 rdmulti -* -b 3 --cnt 6")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run PM3 command: %w", err)
	}

	re := regexp.MustCompile(`\|\s+([0-9A-Fa-f]{2}(?:\s[0-9A-Fa-f]{2})*)\s+\|`)
	matches := re.FindAllStringSubmatch(out.String(), -1)
	if matches == nil {
		return nil, nil
	}

	var memory []byte
	for _, m := range matches {
		bytesStr := strings.Split(m[1], " ")
		for _, bStr := range bytesStr {
			var b byte
			if _, err := fmt.Sscanf(bStr, "%02X", &b); err != nil {
				return nil, nil
			}
			memory = append(memory, b)
		}
	}

	return memory, nil
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
