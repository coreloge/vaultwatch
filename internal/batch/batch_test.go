package batch_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/batch"
	"github.com/yourusername/vaultwatch/internal/lease"
)

func sampleInfo(id string) lease.Info {
	return lease.Info{LeaseID: id, Status: lease.StatusWarning}
}

func TestAdd_FlushesOnMaxSize(t *testing.T) {
	var mu sync.Mutex
	var got []lease.Info

	cfg := batch.Config{MaxSize: 3, Window: 5 * time.Second}
	c := batch.New(cfg, func(_ context.Context, items []lease.Info) {
		mu.Lock()
		got = append(got, items...)
		mu.Unlock()
	})

	ctx := context.Background()
	c.Add(ctx, sampleInfo("a"))
	c.Add(ctx, sampleInfo("b"))
	c.Add(ctx, sampleInfo("c")) // triggers flush

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 3 {
		t.Fatalf("expected 3 items, got %d", len(got))
	}
}

func TestAdd_FlushesAfterWindow(t *testing.T) {
	var mu sync.Mutex
	var got []lease.Info

	cfg := batch.Config{MaxSize: 100, Window: 50 * time.Millisecond}
	c := batch.New(cfg, func(_ context.Context, items []lease.Info) {
		mu.Lock()
		got = append(got, items...)
		mu.Unlock()
	})

	ctx := context.Background()
	c.Add(ctx, sampleInfo("x"))

	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("expected 1 item after window, got %d", len(got))
	}
}

func TestFlush_ForcesImmediateDelivery(t *testing.T) {
	var mu sync.Mutex
	var got []lease.Info

	cfg := batch.Config{MaxSize: 100, Window: 10 * time.Second}
	c := batch.New(cfg, func(_ context.Context, items []lease.Info) {
		mu.Lock()
		got = append(got, items...)
		mu.Unlock()
	})

	ctx := context.Background()
	c.Add(ctx, sampleInfo("force1"))
	c.Add(ctx, sampleInfo("force2"))
	c.Flush(ctx)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 items after forced flush, got %d", len(got))
	}
}

func TestFlush_EmptyBatchIsNoop(t *testing.T) {
	called := false
	cfg := batch.DefaultConfig()
	c := batch.New(cfg, func(_ context.Context, _ []lease.Info) {
		called = true
	})
	c.Flush(context.Background())
	time.Sleep(30 * time.Millisecond)
	if called {
		t.Fatal("handler should not be called for empty batch")
	}
}

func TestDefaultConfig_Sensible(t *testing.T) {
	cfg := batch.DefaultConfig()
	if cfg.MaxSize <= 0 {
		t.Errorf("expected positive MaxSize, got %d", cfg.MaxSize)
	}
	if cfg.Window <= 0 {
		t.Errorf("expected positive Window, got %v", cfg.Window)
	}
}
