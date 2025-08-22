package rfid

import (
	"bufio"
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

const (
	PM3Client      = "/home/paniq/proxmark3/client/proxmark3"
	PM3Port        = "/dev/ttyACM0"
	TargetString   = "enzogenovese.com"
	ScanInterval   = 50 * time.Millisecond // realistic interval
	SnippetPadding = 10
)

type RFIDScanner struct {
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	running    bool
	lastUID    string
	pm3Cmd     *exec.Cmd
	stdin      *bufio.Writer
	stdout     *bufio.Reader
	mu         sync.Mutex
}

var _ hal.Device = (*RFIDScanner)(nil)

func (r *RFIDScanner) Init() error {
	if r.running {
		return nil
	}

	// Start persistent Proxmark3 session
	cmd := exec.Command(PM3Client, PM3Port)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	r.stdin = bufio.NewWriter(stdin)
	r.stdout = bufio.NewReader(stdout)
	r.pm3Cmd = cmd

	if err := cmd.Start(); err != nil {
		return err
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

func (r *RFIDScanner) Close() error {
	if !r.running {
		return nil
	}
	r.cancelFunc()
	r.wg.Wait()
	r.running = false

	if r.pm3Cmd != nil && r.pm3Cmd.Process != nil {
		r.pm3Cmd.Process.Kill()
	}

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
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to write to journal:", err)
	}
}

func (r *RFIDScanner) scanOnce() {
	uid := r.getTagUID()
	if uid != "" && uid != r.lastUID {
		r.lastUID = uid
		journalLog("[SCANNING_RFID] UID=" + uid)

		mem := r.readTagBlocks([]int{0, 1, 2, 3}) // only necessary blocks
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

// ---------------- Helper Functions ----------------
func (r *RFIDScanner) sendPM3Command(cmd string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := r.stdin.WriteString(cmd + "\n"); err != nil {
		return "", err
	}
	r.stdin.Flush()

	var out bytes.Buffer
	for {
		line, err := r.stdout.ReadString('\n')
		if err != nil {
			break
		}
		out.WriteString(line)
		if strings.Contains(line, "pm3>") { // prompt signals command end
			break
		}
	}
	return out.String(), nil
}

func (r *RFIDScanner) getTagUID() string {
	out, err := r.sendPM3Command("hf 15 info")
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`UID\.{3,}\s*([0-9A-F ]+)`)
	match := re.FindStringSubmatch(out)
	if match != nil {
		return strings.ReplaceAll(match[1], " ", "")
	}
	return ""
}

func (r *RFIDScanner) readTagBlocks(blocks []int) []byte {
	var memory []byte
	for _, block := range blocks {
		cmd := fmt.Sprintf("hf 15 readblock %d", block)
		out, err := r.sendPM3Command(cmd)
		if err != nil {
			continue
		}
		re := regexp.MustCompile(`([0-9A-Fa-f]{2})`)
		matches := re.FindAllString(out, -1)
		for _, b := range matches {
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
