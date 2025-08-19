package hal

// Device is a generic HAL device
type Device interface {
    Init() error
    Close() error
    Info() string
}

// Sensor is a device that provides readings
type Sensor interface {
    Device
    Read() (map[string]any, error) // generic key/value reading
}

// Actuator is a device that performs actions
type Actuator interface {
    Device
    Command(cmd string, args ...any) error
}

