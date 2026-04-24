package webhook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/vaultwatch/internal/webhook"
)

func TestSendWithRetry_SucceedsOnSecondAttempt(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := webhook.New(ts.URL, "", 5*time.Second)
	cfg := webhook.RetryConfig{MaxAttempts: 3, Delay: 10 * time.Millisecond}

	if err := s.SendWithRetry(context.Background(), webhook.Payload{LeaseID: "retry-test"}, cfg); err != nil {
		t.Fatalf("expected success on retry, got: %v", err)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestSendWithRetry_ExhaustsAttempts(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	s := webhook.New(ts.URL, "", 5*time.Second)
	cfg := webhook.RetryConfig{MaxAttempts: 3, Delay: 10 * time.Millisecond}

	err := s.SendWithRetry(context.Background(), webhook.Payload{LeaseID: "fail-test"}, cfg)
	if err == nil {
		t.Fatal("expected error after exhausting retries, got nil")
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestSendWithRetry_ContextCancelled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	s := webhook.New(ts.URL, "", 5*time.Second)
	cfg := webhook.DefaultRetryConfig()

	err := s.SendWithRetry(ctx, webhook.Payload{LeaseID: "ctx-test"}, cfg)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}
