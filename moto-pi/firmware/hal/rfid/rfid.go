package rfid

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"sync"
	"time"

	"firmware/hal"
	"firmware/hal/gpio"
)

const (
	PM3Client      = "/home/paniq/proxmark3/client/proxmark3"
	PM3Port        = "/dev/ttyACM0"
	TargetString   = "enzogenovese.com"
	ScanInterval   = 100 * time.Millisecond // HF 15 needs time to respond
	SnippetPadding = 10
	MemoryBlocks   = 4 // adjust number of blocks to scan
)

// RFIDScanner actively scans memory for a target string
type RFIDScanner struct {
	wg      sync.WaitGroup
	running bool
}

var _ hal.Device = (*RFIDScanner)(nil)

func (r *RFIDScanner) Init() error {
	if r.running {
		return nil
	}
	r.running = true

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for r.running {
			r.scanMemory()
			time.Sleep(ScanInterval)
		}
	}()

	return nil
}

func (r *RFIDScanner) Close() error {
	r.running = false
	r.wg.Wait()
	return nil
}

func (r *RFIDScanner) Info() string {
	if r.running {
		return "RFID scanner running"
	}
	return "RFID scanner stopped"
}

func journalLog(msg string) {
	cmd := exec.Command("logger", "-t", "gimo-events", msg)
	cmd.Run()
}

// scanMemory reads the blocks and looks for the target string
func (r *RFIDScanner) scanMemory() {
	mem := readTagBlocks(MemoryBlocks)
	snippet := extractASCIISnippet(mem, []byte(TargetString), SnippetPadding)
	found := snippet != nil
	if found {
		journalLog("[FOUND_VALID_RFID] " + string(snippet))
	}
	gpio.MomentarySwitch(found)
}

// readTagBlocks reads a given number of blocks from HF 15 tag
func readTagBlocks(numBlocks int) []byte {
	var memory []byte
	// read blocks in one command if supported
	cmd := exec.Command(PM3Client, PM3Port, "-c", fmt.Sprintf("hf 15 readblock 0-%d", numBlocks-1))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil
	}

	re := regexp.MustCompile(`([0-9A-Fa-f]{2})`)
	for _, b := range re.FindAllString(out.String(), -1) {
		var val byte
		fmt.Sscanf(b, "%02X", &val)
		memory = append(memory, val)
	}

	return memory
}

// extractASCIISnippet finds the target string in memory and returns surrounding bytes
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
