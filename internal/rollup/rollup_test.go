package rollup_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
	"github.com/yourusername/vaultwatch/internal/rollup"
)

func sampleInfo(id string) lease.Info {
	return lease.Info{
		LeaseID: id,
		TTL:     lease.NewTTLFromSeconds(300),
	}
}

func TestNew_DefaultConfig(t *testing.T) {
	r := rollup.New(rollup.Config{})
	if r == nil {
		t.Fatal("expected non-nil Rollup")
	}
}

func TestAdd_FlushesAfterWindow(t *testing.T) {
	cfg := rollup.Config{
		Window:  50 * time.Millisecond,
		MaxSize: 100,
	}
	r := rollup.New(cfg)

	r.Add(sampleInfo("lease/a"))
	r.Add(sampleInfo("lease/b"))

	select {
	case batch := <-r.Batches():
		if len(batch.Events) != 2 {
			t.Fatalf("expected 2 events, got %d", len(batch.Events))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for batch flush")
	}
}

func TestAdd_FlushesOnMaxSize(t *testing.T) {
	cfg := rollup.Config{
		Window:  10 * time.Second, // long window — should flush on size
		MaxSize: 3,
	}
	r := rollup.New(cfg)

	r.Add(sampleInfo("lease/1"))
	r.Add(sampleInfo("lease/2"))
	r.Add(sampleInfo("lease/3"))

	select {
	case batch := <-r.Batches():
		if len(batch.Events) != 3 {
			t.Fatalf("expected 3 events, got %d", len(batch.Events))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for size-triggered flush")
	}
}

func TestBatch_WindowEndSet(t *testing.T) {
	cfg := rollup.Config{
		Window:  30 * time.Millisecond,
		MaxSize: 100,
	}
	r := rollup.New(cfg)
	before := time.Now()
	r.Add(sampleInfo("lease/x"))

	select {
	case batch := <-r.Batches():
		if batch.WindowEnd.Before(before) {
			t.Error("WindowEnd should be after test start")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out")
	}
}

func TestAdd_EmptyBufNoFlush(t *testing.T) {
	cfg := rollup.Config{
		Window:  20 * time.Millisecond,
		MaxSize: 10,
	}
	r := rollup.New(cfg)
	// Do not add anything; no flush should arrive.
	select {
	case <-r.Batches():
		t.Fatal("unexpected batch from empty buffer")
	case <-time.After(60 * time.Millisecond):
		// expected
	}
}
