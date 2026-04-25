package health_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/health"
)

// mockChecker implements health.Checker for testing.
type mockChecker struct {
	healthy bool
	err     error
}

func (m *mockChecker) IsHealthy() (bool, error) {
	return m.healthy, m.err
}

func TestServeHTTP_Healthy(t *testing.T) {
	h := health.New(&mockChecker{healthy: true})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var s health.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !s.OK || !s.VaultOK {
		t.Errorf("expected ok=true vault_ok=true, got %+v", s)
	}
}

func TestServeHTTP_Unhealthy(t *testing.T) {
	h := health.New(&mockChecker{healthy: false, err: errors.New("connection refused")})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var s health.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if s.OK || s.VaultOK {
		t.Errorf("expected ok=false, got %+v", s)
	}
	if s.Error == "" {
		t.Error("expected non-empty error field")
	}
}

func TestServeHTTP_ContentType(t *testing.T) {
	h := health.New(&mockChecker{healthy: true})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestServeHTTP_CheckedAtPresent(t *testing.T) {
	h := health.New(&mockChecker{healthy: true})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	var s health.Status
	_ = json.NewDecoder(rec.Body).Decode(&s)
	if s.CheckedAt.IsZero() {
		t.Error("expected non-zero checked_at timestamp")
	}
}
