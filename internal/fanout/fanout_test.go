package fanout_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/youorg/vaultwatch/internal/fanout"
	"github.com/youorg/vaultwatch/internal/lease"
)

type countHandler struct {
	calls atomic.Int64
	err   error
}

func (h *countHandler) Handle(_ context.Context, _ lease.Info) error {
	h.calls.Add(1)
	return h.err
}

func sampleInfo() lease.Info {
	return lease.Info{LeaseID: "secret/data/test/abc123", Status: lease.StatusWarning}
}

func TestSend_CallsAllHandlers(t *testing.T) {
	a, b := &countHandler{}, &countHandler{}
	f := fanout.New(a, b)

	errs := f.Send(context.Background(), sampleInfo())

	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if a.calls.Load() != 1 || b.calls.Load() != 1 {
		t.Errorf("expected each handler called once")
	}
}

func TestSend_CollectsErrors(t *testing.T) {
	sentinel := errors.New("handler failed")
	a := &countHandler{err: sentinel}
	b := &countHandler{}
	f := fanout.New(a, b)

	errs := f.Send(context.Background(), sampleInfo())

	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if !errors.Is(errs[0], sentinel) {
		t.Errorf("unexpected error: %v", errs[0])
	}
}

func TestSend_NoHandlers_ReturnsNil(t *testing.T) {
	f := fanout.New()
	errs := f.Send(context.Background(), sampleInfo())
	if errs != nil {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	f := fanout.New()
	if f.Len() != 0 {
		t.Fatalf("expected 0, got %d", f.Len())
	}
	f.Add(&countHandler{})
	if f.Len() != 1 {
		t.Errorf("expected 1, got %d", f.Len())
	}
}

func TestSend_ConcurrentSafety(t *testing.T) {
	const n = 50
	handlers := make([]fanout.Handler, n)
	for i := range handlers {
		handlers[i] = &countHandler{}
	}
	f := fanout.New(handlers...)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for i := 0; i < 20; i++ {
		f.Send(ctx, sampleInfo())
	}
}
