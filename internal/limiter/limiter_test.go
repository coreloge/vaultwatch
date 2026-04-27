package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/limiter"
)

func newLimiter(t *testing.T, cap int) *limiter.Limiter {
	t.Helper()
	l, err := limiter.New(cap)
	if err != nil {
		t.Fatalf("limiter.New: %v", err)
	}
	return l
}

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := limiter.New(0)
	if err == nil {
		t.Fatal("expected error for capacity 0")
	}
}

func TestAcquire_Release_UpdatesActive(t *testing.T) {
	l := newLimiter(t, 3)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if got := l.Active(); got != 1 {
		t.Errorf("Active() = %d, want 1", got)
	}

	l.Release()
	if got := l.Active(); got != 0 {
		t.Errorf("Active() after Release = %d, want 0", got)
	}
}

func TestAcquire_BlocksAtCapacity(t *testing.T) {
	l := newLimiter(t, 1)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first Acquire: %v", err)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Acquire(ctxTimeout)
	if err == nil {
		t.Fatal("expected error when capacity exhausted")
	}
}

func TestCapacity_ReturnsConfiguredValue(t *testing.T) {
	l := newLimiter(t, 5)
	if got := l.Capacity(); got != 5 {
		t.Errorf("Capacity() = %d, want 5", got)
	}
}

func TestRelease_ExtraCallIsNoop(t *testing.T) {
	l := newLimiter(t, 2)
	// Should not panic or block.
	l.Release()
	if got := l.Active(); got != 0 {
		t.Errorf("Active() = %d, want 0", got)
	}
}
