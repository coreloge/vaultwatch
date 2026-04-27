package rollup_test

import (
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/rollup"
)

func TestConcurrentAdd_NoPanic(t *testing.T) {
	cfg := rollup.Config{
		Window:  20 * time.Millisecond,
		MaxSize: 10,
	}
	r := rollup.New(cfg)

	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			r.Add(sampleInfo("lease/concurrent"))
		}(i)
	}

	// Drain batches while goroutines are running.
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	timeout := time.After(500 * time.Millisecond)
	for {
		select {
		case <-r.Batches():
		case <-done:
			return
		case <-timeout:
			t.Fatal("timed out waiting for concurrent adds to complete")
		}
	}
}

func TestMultipleBatches_AllEventsDelivered(t *testing.T) {
	cfg := rollup.Config{
		Window:  10 * time.Millisecond,
		MaxSize: 5,
	}
	r := rollup.New(cfg)

	total := 10
	for i := 0; i < total; i++ {
		r.Add(sampleInfo("lease/multi"))
		time.Sleep(2 * time.Millisecond)
	}

	collected := 0
	timeout := time.After(300 * time.Millisecond)
	for collected < total {
		select {
		case batch := <-r.Batches():
			collected += len(batch.Events)
		case <-timeout:
			t.Fatalf("only collected %d/%d events before timeout", collected, total)
		}
	}
}
