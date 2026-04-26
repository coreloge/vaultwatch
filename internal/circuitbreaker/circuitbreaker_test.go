package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/circuitbreaker"
)

func newBreaker(maxFailures int, openDuration time.Duration) *circuitbreaker.CircuitBreaker {
	return circuitbreaker.New(circuitbreaker.Config{
		MaxFailures:  maxFailures,
		OpenDuration: openDuration,
	})
}

func TestAllow_ClosedByDefault(t *testing.T) {
	cb := newBreaker(3, time.Second)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_OpensAfterMaxFailures(t *testing.T) {
	cb := newBreaker(3, time.Minute)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if err := cb.Allow(); err != circuitbreaker.ErrCircuitOpen {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestAllow_ClosedAfterSuccess(t *testing.T) {
	cb := newBreaker(2, time.Minute)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected circuit closed after success, got %v", err)
	}
}

func TestAllow_HalfOpenAfterDuration(t *testing.T) {
	cb := newBreaker(1, 10*time.Millisecond)
	cb.RecordFailure()
	if err := cb.Allow(); err != circuitbreaker.ErrCircuitOpen {
		t.Fatalf("expected open immediately after failure")
	}
	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected half-open to permit request, got %v", err)
	}
	if cb.CurrentState() != circuitbreaker.StateHalfOpen {
		t.Fatalf("expected StateHalfOpen")
	}
}

func TestHalfOpen_FailureReopens(t *testing.T) {
	cb := newBreaker(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = cb.Allow() // transition to half-open
	cb.RecordFailure()
	if cb.CurrentState() != circuitbreaker.StateOpen {
		t.Fatalf("expected circuit to reopen after half-open failure")
	}
}

func TestDefaultConfig_ReturnsNonZero(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	if cfg.MaxFailures <= 0 {
		t.Errorf("expected positive MaxFailures, got %d", cfg.MaxFailures)
	}
	if cfg.OpenDuration <= 0 {
		t.Errorf("expected positive OpenDuration, got %v", cfg.OpenDuration)
	}
}
