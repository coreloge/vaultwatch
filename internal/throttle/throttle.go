// Package throttle provides per-lease alert throttling to prevent
// alert storms when many leases expire simultaneously.
package throttle

import (
	"sync"
	"time"
)

// Throttler limits how frequently alerts are emitted for a given lease ID.
type Throttler struct {
	mu       sync.Mutex
	window   time.Duration
	lastSent map[string]time.Time
}

// New creates a Throttler that enforces a minimum duration between alerts
// for the same lease ID.
func New(window time.Duration) *Throttler {
	return &Throttler{
		window:   window,
		lastSent: make(map[string]time.Time),
	}
}

// Allow reports whether an alert for leaseID should be allowed through.
// It returns true the first time a leaseID is seen, and again only after
// the configured window has elapsed since the last allowed alert.
func (t *Throttler) Allow(leaseID string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	last, seen := t.lastSent[leaseID]
	if !seen || now.Sub(last) >= t.window {
		t.lastSent[leaseID] = now
		return true
	}
	return false
}

// Reset clears the throttle state for a specific lease ID.
func (t *Throttler) Reset(leaseID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSent, leaseID)
}

// Purge removes all entries whose last-sent time is older than the window,
// preventing unbounded memory growth for expired leases.
func (t *Throttler) Purge() {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := time.Now().Add(-t.window)
	for id, ts := range t.lastSent {
		if ts.Before(cutoff) {
			delete(t.lastSent, id)
		}
	}
}

// Len returns the number of tracked lease IDs.
func (t *Throttler) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.lastSent)
}
