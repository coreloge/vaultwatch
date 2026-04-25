// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently webhook alerts can be dispatched per lease.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls dispatch frequency using a per-key token bucket.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]time.Time
	cooldown time.Duration
}

// New creates a Limiter with the given cooldown period between allowed events
// for the same key.
func New(cooldown time.Duration) *Limiter {
	if cooldown <= 0 {
		cooldown = time.Minute
	}
	return &Limiter{
		buckets:  make(map[string]time.Time),
		cooldown: cooldown,
	}
}

// Allow returns true if the key is permitted to proceed, and records the
// current time as the last-seen timestamp. Returns false if the key was
// seen within the cooldown window.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if last, ok := l.buckets[key]; ok {
		if now.Sub(last) < l.cooldown {
			return false
		}
	}
	l.buckets[key] = now
	return true
}

// Reset removes the rate-limit record for a key, allowing it to proceed
// immediately on the next call to Allow.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Len returns the number of tracked keys.
func (l *Limiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.buckets)
}

// Purge removes all keys whose last-seen time is older than the cooldown,
// freeing memory for leases that are no longer active.
func (l *Limiter) Purge() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for k, t := range l.buckets {
		if now.Sub(t) >= l.cooldown {
			delete(l.buckets, k)
		}
	}
}
