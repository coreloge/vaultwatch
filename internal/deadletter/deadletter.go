// Package deadletter provides a dead-letter store for lease alerts that
// could not be delivered after all retry attempts are exhausted.
package deadletter

import (
	"sync"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Entry holds a failed alert alongside metadata about the failure.
type Entry struct {
	LeaseInfo  lease.Info
	Reason     string
	Attempts   int
	FailedAt   time.Time
	ExpiresAt  time.Time
}

// Store is a bounded, thread-safe dead-letter queue.
type Store struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
	ttl     time.Duration
}

// New returns a Store that retains at most maxSize entries, each surviving for ttl.
// Sensible defaults are applied when zero values are supplied.
func New(maxSize int, ttl time.Duration) *Store {
	if maxSize <= 0 {
		maxSize = 256
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &Store{maxSize: maxSize, ttl: ttl}
}

// Add records a failed delivery. If the store is full the oldest entry is evicted.
func (s *Store) Add(info lease.Info, reason string, attempts int) {
	now := time.Now()
	e := Entry{
		LeaseInfo: info,
		Reason:    reason,
		Attempts:  attempts,
		FailedAt:  now,
		ExpiresAt: now.Add(s.ttl),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeExpiredLocked(now)
	if len(s.entries) >= s.maxSize {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, e)
}

// Drain returns all non-expired entries and clears the store.
func (s *Store) Drain() []Entry {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeExpiredLocked(now)
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	s.entries = s.entries[:0]
	return out
}

// Len returns the number of live entries currently held.
func (s *Store) Len() int {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeExpiredLocked(now)
	return len(s.entries)
}

// purgeExpiredLocked removes entries whose ExpiresAt is in the past.
// Caller must hold s.mu.
func (s *Store) purgeExpiredLocked(now time.Time) {
	keep := s.entries[:0]
	for _, e := range s.entries {
		if e.ExpiresAt.After(now) {
			keep = append(keep, e)
		}
	}
	s.entries = keep
}
