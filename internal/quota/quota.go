// Package quota enforces per-lease-prefix alert dispatch limits within a
// rolling time window, preventing alert storms when many leases share a
// common path prefix.
package quota

import (
	"sync"
	"time"
)

// Config holds quota enforcement settings.
type Config struct {
	// MaxAlerts is the maximum number of alerts allowed per prefix per Window.
	MaxAlerts int
	// Window is the rolling duration over which MaxAlerts is enforced.
	Window time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxAlerts: 10,
		Window:    time.Minute * 5,
	}
}

type entry struct {
	count     int
	windowEnd time.Time
}

// Quota enforces alert rate limits per prefix.
type Quota struct {
	cfg     Config
	mu      sync.Mutex
	buckets map[string]*entry
}

// New creates a Quota enforcer with the given configuration.
func New(cfg Config) *Quota {
	if cfg.MaxAlerts <= 0 {
		cfg.MaxAlerts = DefaultConfig().MaxAlerts
	}
	if cfg.Window <= 0 {
		cfg.Window = DefaultConfig().Window
	}
	return &Quota{
		cfg:     cfg,
		buckets: make(map[string]*entry),
	}
}

// Allow reports whether an alert for the given prefix is within quota.
// It increments the counter for the prefix if allowed.
func (q *Quota) Allow(prefix string) bool {
	now := time.Now()
	q.mu.Lock()
	defer q.mu.Unlock()

	e, ok := q.buckets[prefix]
	if !ok || now.After(e.windowEnd) {
		q.buckets[prefix] = &entry{count: 1, windowEnd: now.Add(q.cfg.Window)}
		return true
	}
	if e.count >= q.cfg.MaxAlerts {
		return false
	}
	e.count++
	return true
}

// Remaining returns the number of alerts still permitted for the prefix
// within the current window, and the time at which the window resets.
func (q *Quota) Remaining(prefix string) (int, time.Time) {
	now := time.Now()
	q.mu.Lock()
	defer q.mu.Unlock()

	e, ok := q.buckets[prefix]
	if !ok || now.After(e.windowEnd) {
		return q.cfg.MaxAlerts, now.Add(q.cfg.Window)
	}
	remaining := q.cfg.MaxAlerts - e.count
	if remaining < 0 {
		remaining = 0
	}
	return remaining, e.windowEnd
}

// Reset clears the quota bucket for the given prefix.
func (q *Quota) Reset(prefix string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.buckets, prefix)
}
