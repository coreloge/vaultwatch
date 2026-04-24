package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/vaultwatch/internal/vault"
)

func newMockVaultServer(t *testing.T, leaseID string, ttl float64, path string) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/sys/leases/lookup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"id":   leaseID,
				"path": path,
				"ttl":  ttl,
			},
		}
		_ = json.NewEncoder(w).Encode(payload)
	})

	mux.HandleFunc("/v1/sys/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"initialized": true,
			"sealed":      false,
			"standby":     false,
		})
	})

	return httptest.NewServer(mux)
}

func TestLookupLease(t *testing.T) {
	leaseID := "database/creds/my-role/abc123"
	srv := newMockVaultServer(t, leaseID, 3600, "database/creds/my-role")
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := client.LookupLease(context.Background(), leaseID)
	if err != nil {
		t.Fatalf("LookupLease: %v", err)
	}

	if info.LeaseID != leaseID {
		t.Errorf("expected lease ID %q, got %q", leaseID, info.LeaseID)
	}
	if info.TTL.Seconds() != 3600 {
		t.Errorf("expected TTL 3600s, got %v", info.TTL)
	}
	if !strings.Contains(info.Path, "database") {
		t.Errorf("unexpected path: %q", info.Path)
	}
}

func TestIsHealthy(t *testing.T) {
	srv := newMockVaultServer(t, "", 0, "")
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if err := client.IsHealthy(context.Background()); err != nil {
		t.Errorf("expected healthy vault, got error: %v", err)
	}
}

func TestNewClient_InvalidAddress(t *testing.T) {
	_, err := vault.NewClient("://bad-url", "token")
	if err == nil {
		t.Error("expected error for invalid address, got nil")
	}
}
