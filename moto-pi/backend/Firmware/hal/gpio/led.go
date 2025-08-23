package gpio

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

// MomentarySwitch handles the GPIO behavior described
func MomentarySwitch() {
	if _, err := host.Init(); err != nil {
		fmt.Println("Failed to init GPIO:", err)
		return
	}

	gpio21 := rpi.P1_40 // GPIO21
	gpio26 := rpi.P1_37 // GPIO26

	// Ensure GPIO26 is ON initially
	gpio26.Out(gpio.High)
	gpio21.Out(gpio.Low)

	// Step 1: turn off GPIO26
	gpio26.Out(gpio.Low)
	// Step 2: turn on GPIO21 for 3 seconds
	gpio21.Out(gpio.High)
	time.Sleep(3 * time.Second)
	// Step 3: restore normal state
	gpio21.Out(gpio.Low)
	gpio26.Out(gpio.High)
}
