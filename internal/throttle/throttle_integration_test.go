package throttle_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vaultwatch/internal/throttle"
)

// TestAllow_ConcurrentSafety verifies that concurrent Allow calls for the
// same lease ID do not cause data races and that exactly one goroutine is
// permitted per window.
func TestAllow_ConcurrentSafety(t *testing.T) {
	th := throttle.New(time.Minute)

	var allowed int64
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if th.Allow("shared-lease") {
				atomic.AddInt64(&allowed, 1)
			}
		}()
	}

	wg.Wait()

	if allowed != 1 {
		t.Fatalf("expected exactly 1 allowed call, got %d", allowed)
	}
}

// TestPurge_ConcurrentSafety checks that Purge and Allow can run
// concurrently without panicking.
func TestPurge_ConcurrentSafety(t *testing.T) {
	th := throttle.New(time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			th.Allow("lease-concurrent")
		}(i)
		go func() {
			defer wg.Done()
			th.Purge()
		}()
	}
	wg.Wait()
}
