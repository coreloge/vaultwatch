package signal_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/signal"
)

func TestNotify_CancelledBySignal(t *testing.T) {
	h := signal.New(syscall.SIGUSR1)
	ctx, stop := h.Notify(context.Background())
	defer stop()

	// Send the signal to the current process.
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1) //nolint:errcheck

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context was not cancelled after signal")
	}
}

func TestNotify_CancelledByParent(t *testing.T) {
	parent, parentCancel := context.WithCancel(context.Background())

	h := signal.New(syscall.SIGUSR1)
	ctx, stop := h.Notify(parent)
	defer stop()

	parentCancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context was not cancelled after parent cancellation")
	}
}

func TestNotify_StopReleasesResources(t *testing.T) {
	h := signal.New(syscall.SIGUSR1)
	_, stop := h.Notify(context.Background())

	// stop should not panic and should be safe to call multiple times.
	stop()
	stop()
}

func TestNew_DefaultSignals(t *testing.T) {
	// Ensure New() without arguments does not panic.
	h := signal.New()
	ctx, stop := h.Notify(context.Background())
	defer stop()

	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}
