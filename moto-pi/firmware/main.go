package main

import (
	"fmt"
	"firmware/hal/gps"
)

func main() {
	g, err := gps.NewGPS("/dev/ttyAMA0", 9600)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	for {
		data, _ := g.Read()
		if data != nil {
			fmt.Printf("%+v\n", data)
		}
	}
}

