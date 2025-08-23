package gps

import (
	"bufio"
	"fmt"
	"sync"
	"time"

	"github.com/adrianmo/go-nmea"
	"go.bug.st/serial"
)

// GPSData holds the parsed values from the GPS
type GPSData struct {
	Time       string
	Latitude   float64
	Longitude  float64
	Altitude   float64
	Satellites int
	SpeedKph   float64
	TrackAngle float64
	ValidFix   bool
}

// GPS implements a background-reading GPS receiver
type GPS struct {
	port       serial.Port
	portName   string
	baudRate   int
	data       GPSData
	mu         sync.RWMutex
	running    bool
	cancelFunc func()
	wg         sync.WaitGroup
}

// NewGPS constructs GPS instance
func NewGPS(portName string, baudRate int) *GPS {
	return &GPS{
		portName: portName,
		baudRate: baudRate,
	}
}

// Init opens the serial port and starts background reading
func (g *GPS) Init() error {
	if g.running {
		return nil
	}

	mode := &serial.Mode{BaudRate: g.baudRate}
	port, err := serial.Open(g.portName, mode)
	if err != nil {
		fmt.Println("Warning: failed to open GPS port:", err)
		g.running = false
		return nil
	}
	g.port = port

	ctxDone := make(chan struct{})
	g.cancelFunc = func() { close(ctxDone) }

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		defer func() { g.running = false }() // reset when goroutine exits

		scanner := bufio.NewScanner(g.port)
		scanner.Split(bufio.ScanLines)

		for {
			select {
			case <-ctxDone:
				return
			default:
				if !scanner.Scan() {
					if err := scanner.Err(); err != nil {
						fmt.Println("GPS read error:", err)
					}
					time.Sleep(100 * time.Millisecond)
					continue
				}

				line := scanner.Text()
				if len(line) == 0 || line[0] != '$' {
					continue
				}

				msg, err := nmea.Parse(line)
				if err != nil {
					continue
				}

				g.mu.Lock()
				switch m := msg.(type) {
				case nmea.GGA:
					g.data.Time = m.Time.String()
					g.data.Latitude = m.Latitude
					g.data.Longitude = m.Longitude
					g.data.Altitude = m.Altitude
					g.data.Satellites = int(m.NumSatellites)
					g.data.ValidFix = m.FixQuality > nmea.Invalid
				case nmea.RMC:
					g.data.Time = m.Time.String()
					g.data.Latitude = m.Latitude
					g.data.Longitude = m.Longitude
					g.data.SpeedKph = m.Speed * 1.852
					g.data.TrackAngle = m.Course
					g.data.ValidFix = m.Validity == "A"
				}
				g.mu.Unlock()

				g.running = true
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	return nil
}

// Read returns the latest GPS data as a map
func (g *GPS) Read() (map[string]any, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return map[string]any{
		"time":       g.data.Time,
		"latitude":   g.data.Latitude,
		"longitude":  g.data.Longitude,
		"altitude":   g.data.Altitude,
		"satellites": g.data.Satellites,
		"speed_kph":  g.data.SpeedKph,
		"track_deg":  g.data.TrackAngle,
		"valid_fix":  g.data.ValidFix,
	}, nil
}

// Info returns online/offline status
func (g *GPS) Info() string {
	if !g.running {
		return "offline"
	}
	if g.data.ValidFix {
		return "online (fix)"
	}
	return "online (no fix)"
}

// Close stops background reading and closes serial port
func (g *GPS) Close() error {
	if g.cancelFunc != nil {
		g.cancelFunc()
	}
	g.wg.Wait()
	if g.port != nil {
		return g.port.Close()
	}
	return nil
}
