package gps

import (
	"bufio"
	"fmt"
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

// GPS is the main struct for managing the GPS device
type GPS struct {
	port serial.Port
	data GPSData
}

// NewGPS opens a serial port and returns a GPS instance
func NewGPS(portName string, baudRate int) (*GPS, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open port %s: %w", portName, err)
	}

	return &GPS{
		port: port,
	}, nil
}

// Read reads one full NMEA sentence, updates internal state, and returns data as a map
func (g *GPS) Read() (map[string]any, error) {
	scanner := bufio.NewScanner(g.port)
	scanner.Split(bufio.ScanLines)

	if scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] != '$' {
			// skip junk/UBX
			return nil, nil
		}

		s, err := nmea.Parse(line)
		if err != nil {
			return nil, nil
		}

		switch msg := s.(type) {
		case nmea.GGA:
			g.data.Time = msg.Time.String()
			g.data.Latitude = msg.Latitude
			g.data.Longitude = msg.Longitude
			g.data.Altitude = msg.Altitude
			g.data.Satellites = int(msg.NumSatellites)
			g.data.ValidFix = msg.FixQuality > nmea.Invalid
		case nmea.RMC:
			g.data.Time = msg.Time.String()
			g.data.Latitude = msg.Latitude
			g.data.Longitude = msg.Longitude
			g.data.SpeedKph = msg.Speed * 1.852 // knots → kph
			g.data.TrackAngle = msg.Course
			g.data.ValidFix = msg.Validity == "A"
		}

		out := map[string]any{
			"time":       g.data.Time,
			"latitude":   g.data.Latitude,
			"longitude":  g.data.Longitude,
			"altitude":   g.data.Altitude,
			"satellites": g.data.Satellites,
			"speed_kph":  g.data.SpeedKph,
			"track_deg":  g.data.TrackAngle,
			"valid_fix":  g.data.ValidFix,
		}
		return out, nil
	}
	return nil, nil
}

// Close closes the GPS serial port
func (g *GPS) Close() error {
	if g.port != nil {
		return g.port.Close()
	}
	return nil
}
