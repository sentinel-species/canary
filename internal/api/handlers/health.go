package handlers

import (
	"encoding/json"
	"net/http"

	"canary/internal/logger"
	"canary/internal/types"
)

type HealthHandler types.HealthHandler

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("health check requested")
		RespondJSON(w, http.StatusOK, types.HealthCheckResponse{Status: "ok"})
	}
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
