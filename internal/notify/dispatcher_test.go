package notify_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/vaultwatch/internal/lease"
	"github.com/example/vaultwatch/internal/notify"
)

func newTestDispatcher(t *testing.T, serverURL, secret, format string) *notify.Dispatcher {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	d, err := notify.New(notify.Config{
		WebhookURL:    serverURL,
		WebhookSecret: secret,
		Format:        format,
	}, logger)
	if err != nil {
		t.Fatalf("failed to create dispatcher: %v", err)
	}
	return d
}

func sampleLeaseInfo() lease.Info {
	return lease.Info{
		LeaseID:   "secret/data/db#abc123",
		Path:      "secret/data/db",
		TTL:       30 * time.Minute,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
}

func TestDispatch_Success(t *testing.T) {
	received := make(chan []byte, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received <- body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := newTestDispatcher(t, server.URL, "", "json")
	err := d.Dispatch(context.Background(), sampleLeaseInfo())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	body := <-received
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("expected valid JSON body, got: %s", body)
	}
	if _, ok := payload["lease_id"]; !ok {
		t.Error("expected 'lease_id' field in payload")
	}
}

func TestDispatch_WebhookError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	d := newTestDispatcher(t, server.URL, "", "json")
	err := d.Dispatch(context.Background(), sampleLeaseInfo())
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestNew_EmptyURL(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	_, err := notify.New(notify.Config{WebhookURL: ""}, logger)
	if err == nil {
		t.Fatal("expected error for empty webhook URL")
	}
}

func TestDispatch_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	d := newTestDispatcher(t, server.URL, "", "json")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := d.Dispatch(ctx, sampleLeaseInfo())
	if err == nil {
		t.Fatal("expected error with cancelled context")
	}
}
