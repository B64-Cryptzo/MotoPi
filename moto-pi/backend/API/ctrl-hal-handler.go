package API

import (
	"encoding/json"
	"net/http"

	"github.com/B64-Cryptzo/MotoPi/backend/Firmware/hal/gps"
	"github.com/B64-Cryptzo/MotoPi/backend/Firmware/hal/rfid"
	"github.com/julienschmidt/httprouter"
)

// HALInterfaceHandler struct to hold interfaces for HAL handling
type HALInterfaceHandler struct {
	*httprouter.Router
	// Embed a HALService to separate stub/live logic
	service HALServiceInterface
}

// HALServiceInterface defines methods the HAL service must implement
type HALServiceInterface interface {
	GetStatus() map[string]interface{}
}

// StubHALService is a stub implementation
type StubHALService struct{}

func (s *StubHALService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"status": "online",
		"temp":   42,
	}
}

// LiveHALService will hit the real PI firmware
type LiveHALService struct {
	RFIDScanner *rfid.RFIDScanner
	GPS         *gps.GPS
}

func (s *LiveHALService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"Proxmark3 Reader": s.RFIDScanner.Info(),
		"GPS Module":       s.GPS.Info(),
	}
}

// NewHALInterfaceHandler creates a new HAL handler
func NewHALInterfaceHandler(service HALServiceInterface, router *httprouter.Router) *HALInterfaceHandler {
	h := &HALInterfaceHandler{
		Router:  router,
		service: service,
	}

	h.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	h.Router.GET("/v1/api/hal/status", h.GetHalStatus)

	return h
}

// GetHalStatus endpoint
func (h *HALInterfaceHandler) GetHalStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := h.service.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
