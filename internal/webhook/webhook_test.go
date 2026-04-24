package webhook_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/vaultwatch/internal/webhook"
)

func newTestServer(t *testing.T, statusCode int, verify func(*http.Request)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if verify != nil {
			verify(r)
		}
		w.WriteHeader(statusCode)
	}))
}

func TestSend_Success(t *testing.T) {
	payload := webhook.Payload{
		LeaseID:   "lease/abc123",
		ExpiresAt: time.Now().Add(5 * time.Minute),
		TTL:       300,
		Message:   "lease expiring soon",
	}

	ts := newTestServer(t, http.StatusOK, func(r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		var got webhook.Payload
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.LeaseID != payload.LeaseID {
			t.Errorf("expected lease_id %q, got %q", payload.LeaseID, got.LeaseID)
		}
	})
	defer ts.Close()

	s := webhook.New(ts.URL, "", 5*time.Second)
	if err := s.Send(context.Background(), payload); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
}

func TestSend_WithSecret(t *testing.T) {
	const secret = "mysecret"
	ts := newTestServer(t, http.StatusOK, func(r *http.Request) {
		if got := r.Header.Get("X-VaultWatch-Secret"); got != secret {
			t.Errorf("expected secret header %q, got %q", secret, got)
		}
	})
	defer ts.Close()

	s := webhook.New(ts.URL, secret, 5*time.Second)
	if err := s.Send(context.Background(), webhook.Payload{LeaseID: "x"}); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
}

func TestSend_NonSuccessStatus(t *testing.T) {
	ts := newTestServer(t, http.StatusInternalServerError, nil)
	defer ts.Close()

	s := webhook.New(ts.URL, "", 5*time.Second)
	err := s.Send(context.Background(), webhook.Payload{LeaseID: "y"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_InvalidURL(t *testing.T) {
	s := webhook.New("http://127.0.0.1:0/nowhere", "", 1*time.Second)
	err := s.Send(context.Background(), webhook.Payload{LeaseID: "z"})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
