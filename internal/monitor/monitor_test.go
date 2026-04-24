package monitor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vaultwatch/internal/config"
	"github.com/vaultwatch/internal/vault"
)

func newTestServer(ttl int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"id":  "test/lease/abc123",
				"ttl": ttl,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

func newTestMonitor(t *testing.T, serverURL string, warningThreshold time.Duration, leaseIDs []string) *Monitor {
	t.Helper()
	cfg := &config.Config{
		VaultAddress:     serverURL,
		VaultToken:       "test-token",
		PollInterval:     time.Second,
		WarningThreshold: warningThreshold,
		LeaseIDs:         leaseIDs,
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	return New(cfg, client)
}

func TestCheckLease_Expiring(t *testing.T) {
	srv := newTestServer(30) // 30 second TTL
	defer srv.Close()

	m := newTestMonitor(t, srv.URL, 5*time.Minute, []string{"test/lease/abc123"})

	status, err := m.checkLease(context.Background(), "test/lease/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Expiring {
		t.Errorf("expected lease to be flagged as expiring (TTL=%s, threshold=%s)",
			status.TTL, m.cfg.WarningThreshold)
	}
}

func TestCheckLease_NotExpiring(t *testing.T) {
	srv := newTestServer(7200) // 2 hour TTL
	defer srv.Close()

	m := newTestMonitor(t, srv.URL, 5*time.Minute, []string{"test/lease/abc123"})

	status, err := m.checkLease(context.Background(), "test/lease/abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Expiring {
		t.Errorf("expected lease NOT to be flagged as expiring (TTL=%s, threshold=%s)",
			status.TTL, m.cfg.WarningThreshold)
	}
}

func TestRun_CancelContext(t *testing.T) {
	srv := newTestServer(3600)
	defer srv.Close()

	m := newTestMonitor(t, srv.URL, 5*time.Minute, []string{"test/lease/abc123"})
	m.cfg.PollInterval = 50 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := m.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
