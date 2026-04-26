package backoff_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/backoff"
)

func newConfig(jitter bool) backoff.Config {
	c := backoff.DefaultConfig()
	c.Jitter = jitter
	return c
}

func TestBackoff_FirstAttemptEqualsInitial(t *testing.T) {
	c := newConfig(false)
	got := c.Backoff(0)
	if got != c.InitialInterval {
		t.Errorf("expected %v, got %v", c.InitialInterval, got)
	}
}

func TestBackoff_GrowsExponentially(t *testing.T) {
	c := newConfig(false)
	prev := c.Backoff(0)
	for i := 1; i <= 3; i++ {
		curr := c.Backoff(i)
		if curr <= prev {
			t.Errorf("attempt %d: expected interval to grow, got %v <= %v", i, curr, prev)
		}
		prev = curr
	}
}

func TestBackoff_CapsAtMaxInterval(t *testing.T) {
	c := newConfig(false)
	for i := 0; i < 20; i++ {
		got := c.Backoff(i)
		if got > c.MaxInterval {
			t.Errorf("attempt %d: %v exceeds MaxInterval %v", i, got, c.MaxInterval)
		}
	}
}

func TestBackoff_NegativeAttemptClamped(t *testing.T) {
	c := newConfig(false)
	got := c.Backoff(-5)
	if got != c.InitialInterval {
		t.Errorf("expected initial interval for negative attempt, got %v", got)
	}
}

func TestBackoff_JitterVariesOutput(t *testing.T) {
	c := newConfig(true)
	seen := make(map[time.Duration]bool)
	for i := 0; i < 20; i++ {
		seen[c.Backoff(2)] = true
	}
	if len(seen) < 2 {
		t.Error("expected jitter to produce varied durations, got identical values")
	}
}

func TestExceeded_BelowMax(t *testing.T) {
	c := backoff.DefaultConfig() // MaxAttempts = 5
	if c.Exceeded(4) {
		t.Error("attempt 4 should not exceed max of 5")
	}
}

func TestExceeded_AtMax(t *testing.T) {
	c := backoff.DefaultConfig()
	if !c.Exceeded(5) {
		t.Error("attempt 5 should exceed max of 5")
	}
}

func TestExceeded_ZeroMaxNeverExceeds(t *testing.T) {
	c := backoff.DefaultConfig()
	c.MaxAttempts = 0
	for _, attempt := range []int{0, 10, 1000} {
		if c.Exceeded(attempt) {
			t.Errorf("attempt %d should never exceed when MaxAttempts=0", attempt)
		}
	}
}
