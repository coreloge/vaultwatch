package dedup_test

import (
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/dedup"
	"github.com/your-org/vaultwatch/internal/lease"
)

// TestConcurrentRecordAndCheck verifies that concurrent Record / IsDuplicate
// calls do not race or panic.
func TestConcurrentRecordAndCheck(t *testing.T) {
	dd := dedup.New(time.Minute)
	const goroutines = 50

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			dd.Record("lease-concurrent", lease.StatusWarning)
		}()
		go func() {
			defer wg.Done()
			_ = dd.IsDuplicate("lease-concurrent", lease.StatusWarning)
		}()
	}

	wg.Wait()
}

// TestConcurrentPurge verifies that Purge is safe to call alongside writes.
func TestConcurrentPurge(t *testing.T) {
	dd := dedup.New(5 * time.Millisecond)
	const goroutines = 30

	var wg sync.WaitGroup
	wg.Add(goroutines + 1)

	for i := 0; i < goroutines; i++ {
		go func(n int) {
			defer wg.Done()
			dd.Record("lease-purge", lease.StatusCritical)
		}(i)
	}

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		dd.Purge()
	}()

	wg.Wait()
}
