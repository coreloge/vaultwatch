// Package limiter provides a token-bucket style concurrency limiter
// for controlling the number of simultaneous webhook dispatches.
package limiter

import (
	"context"
	"fmt"
	"sync"
)

// Limiter controls concurrent access to a shared resource using a semaphore.
type Limiter struct {
	mu      sync.Mutex
	sem     chan struct{}
	cap     int
	active  int
}

// New creates a Limiter with the given concurrency capacity.
// capacity must be >= 1.
func New(capacity int) (*Limiter, error) {
	if capacity < 1 {
		return nil, fmt.Errorf("limiter: capacity must be at least 1, got %d", capacity)
	}
	return &Limiter{
		sem: make(chan struct{}, capacity),
		cap: capacity,
	}, nil
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns an error if the context expires before a slot is acquired.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		l.mu.Lock()
		l.active++
		l.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a previously acquired slot.
func (l *Limiter) Release() {
	select {
	case <-l.sem:
		l.mu.Lock()
		l.active--
		l.mu.Unlock()
	default:
	}
}

// Active returns the number of currently held slots.
func (l *Limiter) Active() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.active
}

// Capacity returns the maximum concurrency allowed.
func (l *Limiter) Capacity() int {
	return l.cap
}
