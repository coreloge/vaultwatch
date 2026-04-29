// Package cache provides a simple TTL-based in-memory cache for lease metadata.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached value alongside its expiry time.
type Entry[V any] struct {
	Value     V
	ExpiresAt time.Time
}

// IsExpired reports whether the entry has passed its expiry time.
func (e Entry[V]) IsExpired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// Cache is a generic TTL-based in-memory store.
type Cache[K comparable, V any] struct {
	mu  sync.RWMutex
	items map[K]Entry[V]
	ttl  time.Duration
}

// New creates a Cache with the given TTL applied to every stored entry.
func New[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	if ttl <= 0 {
		ttl = time.Minute
	}
	return &Cache[K, V]{
		items: make(map[K]Entry[V]),
		ttl:  ttl,
	}
}

// Set stores value under key, resetting its TTL.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry[V]{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get returns the cached value and true if present and not expired.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok || entry.IsExpired(time.Now()) {
		var zero V
		return zero, false
	}
	return entry.Value, true
}

// Delete removes the entry for key.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Purge removes all entries that have expired.
func (c *Cache[K, V]) Purge() int {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	removed := 0
	for k, entry := range c.items {
		if entry.IsExpired(now) {
			delete(c.items, k)
			removed++
		}
	}
	return removed
}

// Len returns the total number of entries, including expired ones.
func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
