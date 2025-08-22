package gpio

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

func ToggleLED() {
	if _, err := host.Init(); err != nil {
		fmt.Println("Failed to init GPIO:", err)
		return
	}

	led := rpi.P1_40 // GPIO21 physical pin
	if err := led.Out(gpio.Low); err != nil {
		fmt.Println("Failed to set initial LED state:", err)
		return
	}

	// pulse high for 50ms to toggle T flip-flop
	led.Out(gpio.High)
	time.Sleep(50 * time.Millisecond)
	led.Out(gpio.Low)
}
