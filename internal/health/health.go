// Package health provides an HTTP handler for exposing daemon liveness
// and readiness status, including Vault connectivity.
package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// Checker is implemented by anything that can report its health.
type Checker interface {
	IsHealthy() (bool, error)
}

// Status represents the JSON response body for health endpoints.
type Status struct {
	OK        bool      `json:"ok"`
	VaultOK   bool      `json:"vault_ok"`
	CheckedAt time.Time `json:"checked_at"`
	Error     string    `json:"error,omitempty"`
}

// Handler holds dependencies for the health HTTP handler.
type Handler struct {
	vault Checker
}

// New creates a Handler with the provided Vault health checker.
func New(vault Checker) *Handler {
	return &Handler{vault: vault}
}

// ServeHTTP writes a JSON health status. Returns 200 when healthy, 503 otherwise.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ok, err := h.vault.IsHealthy()

	s := Status{
		OK:        ok,
		VaultOK:   ok,
		CheckedAt: time.Now().UTC(),
	}
	if err != nil {
		s.Error = err.Error()
	}

	code := http.StatusOK
	if !ok {
		code = http.StatusServiceUnavailable
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(s)
}
