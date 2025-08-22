package rfid

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"firmware/hal"
	"firmware/hal/gpio"
)

const (
	PM3Client      = "/home/paniq/proxmark3/client/proxmark3"
	PM3Port        = "/dev/ttyACM0"
	TargetString   = "enzogenovese.com"
	ScanInterval   = 10 * time.Millisecond // realistic interval for HF 15
	SnippetPadding = 10
)

type RFIDScanner struct {
	wg      sync.WaitGroup
	running bool
	lastUID string
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
			r.scanOnce()
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

func (r *RFIDScanner) scanOnce() {
	uid := getTagUID()
	if uid != "" && uid != r.lastUID {
		r.lastUID = uid
		journalLog("[SCANNING_RFID] UID=" + uid)

		// Only read first few blocks; adjust as needed
		mem := readTagBlocks([]int{0, 1, 2, 3})
		snippet := extractASCIISnippet(mem, []byte(TargetString), SnippetPadding)
		validRFIDTag := snippet != nil && strings.Contains(string(snippet), TargetString)
		if validRFIDTag {
			journalLog("[FOUND_VALID_RFID] UID=" + uid)
		}

		gpio.MomentarySwitch(validRFIDTag)

	} else if uid == "" {
		r.lastUID = ""
	}
}

func getTagUID() string {
	cmd := exec.Command(PM3Client, PM3Port, "-c", "hf 15 uid")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	re := regexp.MustCompile(`UID\s*:\s*([0-9A-F ]+)`)
	match := re.FindStringSubmatch(out.String())
	if match != nil {
		return strings.ReplaceAll(match[1], " ", "")
	}
	return ""
}

func readTagBlocks(blocks []int) []byte {
	var memory []byte
	for _, block := range blocks {
		cmd := exec.Command(PM3Client, PM3Port, "-c", fmt.Sprintf("hf 15 readblock %d", block))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			continue
		}
		re := regexp.MustCompile(`([0-9A-Fa-f]{2})`)
		for _, b := range re.FindAllString(out.String(), -1) {
			var val byte
			fmt.Sscanf(b, "%02X", &val)
			memory = append(memory, val)
		}
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
