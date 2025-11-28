package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// HealthResponse represents health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// ReadinessResponse represents readiness check response
type ReadinessResponse struct {
	Status  string            `json:"status"`
	Checks  map[string]string `json:"checks"`
	Ready   bool              `json:"ready"`
	Version string            `json:"version,omitempty"`
}

// health handles basic health check (liveness probe)
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) { //nolint:revive // r unused in health check
	response := HealthResponse{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode health response", slog.String("error", err.Error()))
	}
}

// ready handles readiness check
func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)
	ready := true

	// Check database connectivity
	ctx := r.Context()
	if err := h.services.HealthCheck(ctx); err != nil {
		checks["database"] = "unhealthy: " + err.Error()
		ready = false
		h.log.Error("database health check failed", slog.String("error", err.Error()))
	} else {
		checks["database"] = "healthy"
	}

	// Prepare response
	response := ReadinessResponse{
		Status: "ok",
		Checks: checks,
		Ready:  ready,
	}

	w.Header().Set("Content-Type", "application/json")

	if ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode readiness response", slog.String("error", err.Error()))
	}
}
