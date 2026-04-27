package batch_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/batch"
	"github.com/yourusername/vaultwatch/internal/lease"
)

func TestConcurrentAdd_NoPanic(t *testing.T) {
	cfg := batch.Config{MaxSize: 10, Window: 20 * time.Millisecond}
	var total int64
	c := batch.New(cfg, func(_ context.Context, items []lease.Info) {
		atomic.AddInt64(&total, int64(len(items)))
	})

	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Add(ctx, lease.Info{LeaseID: "concurrent", Status: lease.StatusCritical})
		}(i)
	}
	wg.Wait()
	c.Flush(ctx)

	time.Sleep(60 * time.Millisecond)

	if got := atomic.LoadInt64(&total); got != 50 {
		t.Errorf("expected 50 total items delivered, got %d", got)
	}
}

func TestMultipleFlushCycles_AllDelivered(t *testing.T) {
	cfg := batch.Config{MaxSize: 5, Window: 200 * time.Millisecond}
	var total int64
	c := batch.New(cfg, func(_ context.Context, items []lease.Info) {
		atomic.AddInt64(&total, int64(len(items)))
	})

	ctx := context.Background()
	for cycle := 0; cycle < 3; cycle++ {
		for i := 0; i < 5; i++ {
			c.Add(ctx, lease.Info{LeaseID: "cycle", Status: lease.StatusWarning})
		}
		time.Sleep(30 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)

	if got := atomic.LoadInt64(&total); got != 15 {
		t.Errorf("expected 15 total items across 3 cycles, got %d", got)
	}
}
