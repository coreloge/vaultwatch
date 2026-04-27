package window_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/window"
)

func newCounter(d time.Duration) *window.Counter {
	return window.New(d)
}

func TestNew_PanicsOnZeroDuration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero duration")
		}
	}()
	window.New(0)
}

func TestCount_EmptyReturnsZero(t *testing.T) {
	c := newCounter(time.Minute)
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_IncrementsCount(t *testing.T) {
	c := newCounter(time.Minute)
	c.Add()
	c.Add()
	if got := c.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCountAt_ExcludesOldEvents(t *testing.T) {
	c := newCounter(time.Minute)
	now := time.Now()

	// event two minutes ago — outside the window
	c.AddAt(now.Add(-2 * time.Minute))
	// event 30 seconds ago — inside the window
	c.AddAt(now.Add(-30 * time.Second))

	if got := c.CountAt(now); got != 1 {
		t.Fatalf("expected 1 event in window, got %d", got)
	}
}

func TestCountAt_AllEventsWithinWindow(t *testing.T) {
	c := newCounter(time.Hour)
	now := time.Now()
	for i := 0; i < 5; i++ {
		c.AddAt(now.Add(-time.Duration(i) * time.Minute))
	}
	if got := c.CountAt(now); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestReset_ClearsAllEvents(t *testing.T) {
	c := newCounter(time.Minute)
	c.Add()
	c.Add()
	c.Reset()
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCountAt_ExactBoundaryExcluded(t *testing.T) {
	c := newCounter(time.Minute)
	now := time.Now()
	// event exactly at the cutoff boundary (strictly before cutoff → evicted)
	c.AddAt(now.Add(-time.Minute))
	if got := c.CountAt(now); got != 0 {
		t.Fatalf("expected 0 for event at exact boundary, got %d", got)
	}
}
