package API

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// MotorcycleInterfaceHandler struct to hold interfaces for Motorcycle handling
type MotorcycleInterfaceHandler struct {
	*httprouter.Router
	// Embed a MotorcycleService to separate stub/live logic
	service MotorcycleServiceInterface
}

// NewMotorcycleInterfaceHandler creates a new Motorcycle handler
func NewMotorcycleInterfaceHandler(service MotorcycleServiceInterface, router *httprouter.Router) *MotorcycleInterfaceHandler {
	h := &MotorcycleInterfaceHandler{
		Router:  router,
		service: service,
	}

	h.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	h.Router.GET("/v1/api/motorcycle/status", h.GetMotorcycleStatus)
	h.Router.GET("/v1/api/motorcycle/gps", h.GetMotorcycleGPSData)

	return h
}

// MotorcycleServiceInterface defines methods the Motorcycle service must implement
type MotorcycleServiceInterface interface {
	GetStatus() map[string]interface{}
	GetGPSData() map[string]interface{}
}

// StubMotorcycleService is a stub implementation
type StubMotorcycleService struct{}

func (s *StubMotorcycleService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"status": "online",
	}
}

func (s *StubMotorcycleService) GetGPSData() map[string]interface{} {
	return map[string]interface{}{
		"lat": "1.111",
		"lng": "2.222",
	}
}

// LiveMotorcycleService will hit the real PI firmware
type LiveMotorcycleService struct{}

func (s *LiveMotorcycleService) GetStatus() map[string]interface{} {
	// TODO: implement calls to the actual Motorcycle/firmware
	return map[string]interface{}{
		"status": "offline",
	}
}

func (s *LiveMotorcycleService) GetGPSData() map[string]interface{} {
	return map[string]interface{}{
		"lat": "0.0",
		"lng": "0.0",
	}
}

// GetMotorcycleStatus endpoint
func (h *MotorcycleInterfaceHandler) GetMotorcycleStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := h.service.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetMotorcycleGPSData endpoint
func (h *MotorcycleInterfaceHandler) GetMotorcycleGPSData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := h.service.GetGPSData()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
