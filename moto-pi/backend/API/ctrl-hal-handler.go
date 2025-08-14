package API

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
type LiveHALService struct{}

func (s *LiveHALService) GetStatus() map[string]interface{} {
	// TODO: implement calls to the actual HAL/firmware
	return map[string]interface{}{
		"status": "offline",
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
