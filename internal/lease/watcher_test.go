package lease_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
)

func newTestWatcher(t *testing.T, store *lease.Store, handler lease.WatchHandler) *lease.Watcher {
	t.Helper()
	checker := lease.NewExpiryChecker(lease.DefaultExpiryWindow())
	return lease.NewWatcher(lease.WatcherConfig{
		Store:    store,
		Checker:  checker,
		Interval: 50 * time.Millisecond,
		OnExpiry: handler,
	})
}

func TestWatcher_TriggersOnCriticalLease(t *testing.T) {
	store := lease.NewStore()
	info := lease.Info{
		LeaseID: "secret/data/db#abc",
		TTL:     lease.NewTTLFromSeconds(30), // within critical threshold
		Status:  lease.StatusCritical,
	}
	store.Set(info)

	var mu sync.Mutex
	var received []lease.Info

	watcher := newTestWatcher(t, store, func(i lease.Info) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, i)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	watcher.Run(ctx)

	mu.Lock()
	defer mu.Unlock()
	if len(received) == 0 {
		t.Fatal("expected handler to be called for critical lease, got none")
	}
	if received[0].LeaseID != info.LeaseID {
		t.Errorf("expected lease ID %q, got %q", info.LeaseID, received[0].LeaseID)
	}
}

func TestWatcher_SkipsHealthyLease(t *testing.T) {
	store := lease.NewStore()
	info := lease.Info{
		LeaseID: "secret/data/ok#xyz",
		TTL:     lease.NewTTLFromSeconds(3600),
		Status:  lease.StatusOK,
	}
	store.Set(info)

	var mu sync.Mutex
	var received []lease.Info

	watcher := newTestWatcher(t, store, func(i lease.Info) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, i)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	watcher.Run(ctx)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 0 {
		t.Errorf("expected no handler calls for healthy lease, got %d", len(received))
	}
}

func TestWatcher_DefaultInterval(t *testing.T) {
	store := lease.NewStore()
	checker := lease.NewExpiryChecker(lease.DefaultExpiryWindow())
	w := lease.NewWatcher(lease.WatcherConfig{
		Store:   store,
		Checker: checker,
	})
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}
