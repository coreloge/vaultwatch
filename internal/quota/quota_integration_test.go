package quota_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/quota"
)

func TestConcurrentAllow_NoPanic(t *testing.T) {
	q := quota.New(quota.Config{MaxAlerts: 50, Window: time.Second})
	const goroutines = 20
	const callsEach = 10

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < callsEach; j++ {
				q.Allow("secret/concurrent/")
			}
		}()
	}
	wg.Wait()
}

func TestConcurrentAllow_RespectsMax(t *testing.T) {
	const max = 10
	q := quota.New(quota.Config{MaxAlerts: max, Window: time.Minute})

	var allowed atomic.Int64
	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if q.Allow("secret/stress/") {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got > max {
		t.Fatalf("allowed %d calls, expected at most %d", got, max)
	}
}
