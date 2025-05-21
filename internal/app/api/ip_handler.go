package api

import (
	"encoding/json"
	"net/http"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/services"
)

// IPHandler handles HTTP requests for IP-related operations
type IPHandler struct {
	service services.IPService
}

// NewIPHandler creates a new instance of IPHandler
func NewIPHandler(service services.IPService) *IPHandler {
	return &IPHandler{service: service}
}

// FindByIP handles GET requests to find IP information
func (h *IPHandler) FindByIP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get IP from query parameter
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		http.Error(w, "IP parameter is required", http.StatusBadRequest)
		return
	}

	// Call service to get IP information
	info, err := h.service.FindByIP(r.Context(), ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// RegisterRoutes registers the IP handler routes
func (h *IPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/ip", h.FindByIP)
}
