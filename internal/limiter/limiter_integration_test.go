package limiter_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/yourusername/vaultwatch/internal/limiter"
)

func TestConcurrentAcquire_NeverExceedsCapacity(t *testing.T) {
	const cap = 4
	const goroutines = 32

	l, _ := limiter.New(cap)
	var peak int64
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Acquire(context.Background()); err != nil {
				return
			}
			current := int64(l.Active())
			for {
				old := atomic.LoadInt64(&peak)
				if current <= old || atomic.CompareAndSwapInt64(&peak, old, current) {
					break
				}
			}
			l.Release()
		}()
	}
	wg.Wait()

	if peak > int64(cap) {
		t.Errorf("peak active = %d, must not exceed capacity %d", peak, cap)
	}
}

func TestConcurrentRelease_NoPanic(t *testing.T) {
	l, _ := limiter.New(8)
	var wg sync.WaitGroup

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire(context.Background())
			l.Release()
		}()
	}
	wg.Wait()

	if got := l.Active(); got != 0 {
		t.Errorf("Active() after all releases = %d, want 0", got)
	}
}
