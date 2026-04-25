// Package suppress provides alert suppression windows to prevent
// repeated notifications for the same lease within a configurable duration.
package suppress

import (
	"sync"
	"time"
)

// Window represents an active suppression entry.
type Window struct {
	SuppressedAt time.Time
	ExpiresAt    time.Time
}

// Suppressor tracks which lease IDs are currently suppressed.
type Suppressor struct {
	mu       sync.Mutex
	windows  map[string]Window
	duration time.Duration
}

// New creates a Suppressor with the given suppression duration.
func New(d time.Duration) *Suppressor {
	return &Suppressor{
		windows:  make(map[string]Window),
		duration: d,
	}
}

// IsSuppressed returns true if the given lease ID is currently suppressed.
func (s *Suppressor) IsSuppressed(leaseID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.windows[leaseID]
	if !ok {
		return false
	}
	if time.Now().After(w.ExpiresAt) {
		delete(s.windows, leaseID)
		return false
	}
	return true
}

// Suppress marks the given lease ID as suppressed for the configured duration.
func (s *Suppressor) Suppress(leaseID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.windows[leaseID] = Window{
		SuppressedAt: now,
		ExpiresAt:    now.Add(s.duration),
	}
}

// Release removes any suppression window for the given lease ID.
func (s *Suppressor) Release(leaseID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.windows, leaseID)
}

// Active returns the number of currently active suppression windows.
func (s *Suppressor) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	count := 0
	for id, w := range s.windows {
		if now.After(w.ExpiresAt) {
			delete(s.windows, id)
		} else {
			count++
		}
	}
	return count
}
