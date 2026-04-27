// Package replay provides a mechanism to replay missed or failed lease
// alert events, ensuring no critical notifications are lost during downtime.
package replay

import (
	"sync"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Entry holds a lease event that failed to dispatch and is pending replay.
type Entry struct {
	LeaseInfo lease.Info
	RecordedAt time.Time
	Attempts   int
}

// Store holds pending replay entries.
type Store struct {
	mu      sync.Mutex
	entries []*Entry
	maxAge  time.Duration
}

// New returns a new Store. maxAge defines how long an entry is eligible
// for replay before it is considered stale and discarded.
func New(maxAge time.Duration) *Store {
	if maxAge <= 0 {
		maxAge = 1 * time.Hour
	}
	return &Store{maxAge: maxAge}
}

// Add records a failed lease event for later replay.
func (s *Store) Add(info lease.Info) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, &Entry{
		LeaseInfo:  info,
		RecordedAt: time.Now(),
		Attempts:   0,
	})
}

// Drain returns all non-stale entries and removes them from the store.
// Callers are responsible for re-adding entries that fail again.
func (s *Store) Drain() []*Entry {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	var eligible []*Entry
	var remaining []*Entry

	for _, e := range s.entries {
		if now.Sub(e.RecordedAt) <= s.maxAge {
			eligible = append(eligible, e)
		}
		// stale entries are silently dropped
	}
	s.entries = remaining
	return eligible
}

// Len returns the current number of pending entries.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}

// Purge removes all entries from the store.
func (s *Store) Purge() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = nil
}
