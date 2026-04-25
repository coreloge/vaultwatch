// Package snapshot provides point-in-time capture and diffing of lease state,
// allowing vaultwatch to detect changes between monitor cycles.
package snapshot

import (
	"sync"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
)

// Snapshot holds a captured view of all known leases at a specific moment.
type Snapshot struct {
	CapturedAt time.Time
	Leases     map[string]lease.Info
}

// Diff describes the changes between two snapshots.
type Diff struct {
	Added   []lease.Info
	Removed []lease.Info
	Changed []Change
}

// Change represents a lease whose state has changed between snapshots.
type Change struct {
	Previous lease.Info
	Current  lease.Info
}

// Store maintains the most recent snapshot and supports atomic replacement.
type Store struct {
	mu       sync.RWMutex
	current  *Snapshot
}

// NewStore returns an initialised snapshot Store with no current snapshot.
func NewStore() *Store {
	return &Store{}
}

// Capture records a new snapshot from the provided lease infos and returns it.
// The snapshot replaces any previously stored one.
func (s *Store) Capture(leases []lease.Info) *Snapshot {
	snap := &Snapshot{
		CapturedAt: time.Now().UTC(),
		Leases:     make(map[string]lease.Info, len(leases)),
	}
	for _, l := range leases {
		snap.Leases[l.LeaseID] = l
	}

	s.mu.Lock()
	s.current = snap
	s.mu.Unlock()

	return snap
}

// Current returns the most recently captured snapshot, or nil if none exists.
func (s *Store) Current() *Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Compare returns the diff between a previous and current snapshot.
// If prev is nil every lease in curr is treated as added.
func Compare(prev, curr *Snapshot) Diff {
	var d Diff

	if prev == nil {
		for _, l := range curr.Leases {
			d.Added = append(d.Added, l)
		}
		return d
	}

	// Detect added or changed leases.
	for id, cur := range curr.Leases {
		if old, ok := prev.Leases[id]; !ok {
			d.Added = append(d.Added, cur)
		} else if leaseChanged(old, cur) {
			d.Changed = append(d.Changed, Change{Previous: old, Current: cur})
		}
	}

	// Detect removed leases.
	for id, old := range prev.Leases {
		if _, ok := curr.Leases[id]; !ok {
			d.Removed = append(d.Removed, old)
		}
	}

	return d
}

// leaseChanged reports whether any meaningful field has changed between two
// Info values for the same lease ID.
func leaseChanged(a, b lease.Info) bool {
	return a.Status != b.Status ||
		a.Renewable != b.Renewable ||
		abs(a.TTL.Remaining()-b.TTL.Remaining()) > 5*time.Second
}

func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
