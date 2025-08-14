package API

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// NetworkInterfaceHandler struct to hold interfaces for Network handling
type NetworkInterfaceHandler struct {
	*httprouter.Router
	// Embed a NetworkService to separate stub/live logic
	service NetworkServiceInterface
}

// NetworkServiceInterface defines methods the Network service must implement
type NetworkServiceInterface interface {
	GetStatus() map[string]interface{}
}

// StubNetworkService is a stub implementation
type StubNetworkService struct{}

func (s *StubNetworkService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"status": "online",
	}
}

// LiveNetworkService will hit the real PI firmware
type LiveNetworkService struct{}

func (s *LiveNetworkService) GetStatus() map[string]interface{} {
	// TODO: implement calls to the actual Network/firmware
	return map[string]interface{}{
		"status": "offline",
	}
}

// NewNetworkInterfaceHandler creates a new Network handler
func NewNetworkInterfaceHandler(service NetworkServiceInterface, router *httprouter.Router) *NetworkInterfaceHandler {
	h := &NetworkInterfaceHandler{
		Router:  router,
		service: service,
	}

	h.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	h.Router.GET("/v1/api/network/status", h.GetNetworkStatus)

	return h
}

// GetNetworkStatus endpoint
func (h *NetworkInterfaceHandler) GetNetworkStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := h.service.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
