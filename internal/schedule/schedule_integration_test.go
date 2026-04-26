package schedule_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/schedule"
)

// TestConcurrentStop verifies that calling Stop from multiple goroutines
// concurrently does not panic or deadlock.
func TestConcurrentStop(t *testing.T) {
	s := schedule.New(10 * time.Millisecond)
	ctx := context.Background()

	go s.Run(ctx, func(_ context.Context) {})
	time.Sleep(25 * time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Stop()
		}()
	}
	wg.Wait()
}

// TestRun_MultipleSchedulers ensures independent schedulers do not
// interfere with each other's tick counts.
func TestRun_MultipleSchedulers(t *testing.T) {
	const n = 4
	var counts [n]int64

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			s := schedule.New(30 * time.Millisecond)
			s.Run(ctx, func(_ context.Context) {
				atomic.AddInt64(&counts[idx], 1)
			})
		}()
	}
	wg.Wait()

	for i, c := range counts {
		if c < 2 {
			t.Errorf("scheduler %d: expected >=2 calls, got %d", i, c)
		}
	}
}
