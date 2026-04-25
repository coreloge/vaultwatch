// Package dedup provides alert deduplication to prevent sending
// repeated notifications for the same lease state within a configurable window.
package dedup

import (
	"sync"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Entry holds the last-seen state for a lease.
type Entry struct {
	Status    lease.Status
	SeenAt    time.Time
}

// Deduplicator tracks previously alerted lease states.
type Deduplicator struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]Entry
}

// New returns a Deduplicator with the given deduplication window.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		window:  window,
		entries: make(map[string]Entry),
	}
}

// IsDuplicate reports whether an alert for leaseID with the given status
// was already sent within the deduplication window.
func (d *Deduplicator) IsDuplicate(leaseID string, status lease.Status) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	e, ok := d.entries[leaseID]
	if !ok {
		return false
	}
	if e.Status != status {
		return false
	}
	return time.Since(e.SeenAt) < d.window
}

// Record marks leaseID + status as seen at the current time.
func (d *Deduplicator) Record(leaseID string, status lease.Status) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.entries[leaseID] = Entry{
		Status: status,
		SeenAt: time.Now(),
	}
}

// Evict removes the deduplication entry for leaseID.
func (d *Deduplicator) Evict(leaseID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, leaseID)
}

// Purge removes all entries whose window has elapsed.
func (d *Deduplicator) Purge() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for id, e := range d.entries {
		if time.Since(e.SeenAt) >= d.window {
			delete(d.entries, id)
		}
	}
}
