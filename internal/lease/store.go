package lease

import (
	"fmt"
	"sync"
)

// Store is a thread-safe in-memory registry of tracked leases.
type Store struct {
	mu     sync.RWMutex
	leases map[string]Info
}

// NewStore creates an empty lease Store.
func NewStore() *Store {
	return &Store{
		leases: make(map[string]Info),
	}
}

// Set adds or updates a lease in the store.
func (s *Store) Set(l Info) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.leases[l.LeaseID] = l
}

// Get retrieves a lease by ID. Returns an error if not found.
func (s *Store) Get(leaseID string) (Info, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.leases[leaseID]
	if !ok {
		return Info{}, fmt.Errorf("lease %q not found", leaseID)
	}
	return l, nil
}

// Delete removes a lease from the store.
func (s *Store) Delete(leaseID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.leases, leaseID)
}

// All returns a snapshot of all leases currently tracked.
func (s *Store) All() []Info {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Info, 0, len(s.leases))
	for _, l := range s.leases {
		result = append(result, l)
	}
	return result
}

// Count returns the number of leases in the store.
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.leases)
}
