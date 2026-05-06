// Package clock provides a mockable time source for deterministic testing
// of time-sensitive components such as lease expiry and TTL calculations.
package clock

import (
	"sync"
	"time"
)

// Clock is an interface for obtaining the current time.
type Clock interface {
	Now() time.Time
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
}

// Real is a Clock backed by the system clock.
type Real struct{}

// New returns a Real clock backed by the system time.
func New() Clock {
	return &Real{}
}

// Now returns the current system time.
func (r *Real) Now() time.Time { return time.Now() }

// Since returns the duration elapsed since t.
func (r *Real) Since(t time.Time) time.Duration { return time.Since(t) }

// Until returns the duration until t.
func (r *Real) Until(t time.Time) time.Duration { return time.Until(t) }

// Mock is a Clock whose current time can be controlled programmatically.
type Mock struct {
	mu  sync.RWMutex
	now time.Time
}

// NewMock returns a Mock clock set to the given initial time.
func NewMock(initial time.Time) *Mock {
	return &Mock{now: initial}
}

// Now returns the mock's current time.
func (m *Mock) Now() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.now
}

// Since returns the duration elapsed since t relative to the mock clock.
func (m *Mock) Since(t time.Time) time.Duration {
	return m.Now().Sub(t)
}

// Until returns the duration until t relative to the mock clock.
func (m *Mock) Until(t time.Time) time.Duration {
	return t.Sub(m.Now())
}

// Advance moves the mock clock forward by the given duration.
func (m *Mock) Advance(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = m.now.Add(d)
}

// Set sets the mock clock to an absolute time.
func (m *Mock) Set(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = t
}
