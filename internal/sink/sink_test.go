package sink_test

import (
	"context"
	"errors"
	"testing"

	"github.com/youorg/vaultwatch/internal/alert"
	"github.com/youorg/vaultwatch/internal/sink"
)

// stubTarget is a Target that records calls and optionally returns an error.
type stubTarget struct {
	name    string
	err     error
	called  int
	lastPayload alert.Payload
}

func (s *stubTarget) Name() string { return s.name }
func (s *stubTarget) Send(_ context.Context, p alert.Payload) error {
	s.called++
	s.lastPayload = p
	return s.err
}

func newPayload() alert.Payload {
	return alert.Payload{
		"lease_id": "secret/data/db#abc123",
		"status":   "critical",
	}
}

func TestSendAll_DeliversToAllTargets(t *testing.T) {
	a := &stubTarget{name: "a"}
	b := &stubTarget{name: "b"}
	s := sink.New(a, b)

	if err := s.SendAll(context.Background(), newPayload()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.called != 1 || b.called != 1 {
		t.Errorf("expected each target called once, got a=%d b=%d", a.called, b.called)
	}
}

func TestSendAll_CollectsErrors(t *testing.T) {
	a := &stubTarget{name: "a", err: errors.New("timeout")}
	b := &stubTarget{name: "b"}
	s := sink.New(a, b)

	err := s.SendAll(context.Background(), newPayload())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// b should still have been called despite a failing
	if b.called != 1 {
		t.Errorf("expected b to be called even after a failed, got %d", b.called)
	}
}

func TestSendAll_AllFail_ReturnsError(t *testing.T) {
	a := &stubTarget{name: "a", err: errors.New("err-a")}
	b := &stubTarget{name: "b", err: errors.New("err-b")}
	s := sink.New(a, b)

	err := s.SendAll(context.Background(), newPayload())
	if err == nil {
		t.Fatal("expected combined error")
	}
}

func TestSendAll_NoTargets_ReturnsNil(t *testing.T) {
	s := sink.New()
	if err := s.SendAll(context.Background(), newPayload()); err != nil {
		t.Fatalf("unexpected error with no targets: %v", err)
	}
}

func TestLen_ReturnsTargetCount(t *testing.T) {
	s := sink.New(&stubTarget{name: "x"}, &stubTarget{name: "y"})
	if got := s.Len(); got != 2 {
		t.Errorf("expected Len()=2, got %d", got)
	}
}
