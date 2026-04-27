// Package window provides a sliding time-window counter for tracking
// event frequency over a rolling duration.
package window

import (
	"sync"
	"time"
)

// Counter tracks how many events have occurred within a sliding window.
type Counter struct {
	mu       sync.Mutex
	duration time.Duration
	events   []time.Time
}

// New creates a Counter with the given sliding window duration.
// Panics if duration is zero or negative.
func New(duration time.Duration) *Counter {
	if duration <= 0 {
		panic("window: duration must be positive")
	}
	return &Counter{duration: duration}
}

// Add records a new event at the current time.
func (c *Counter) Add() {
	c.AddAt(time.Now())
}

// AddAt records a new event at the specified time.
func (c *Counter) AddAt(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, t)
	c.evict(t)
}

// Count returns the number of events within the current sliding window.
func (c *Counter) Count() int {
	return c.CountAt(time.Now())
}

// CountAt returns the number of events within the sliding window ending at t.
func (c *Counter) CountAt(t time.Time) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict(t)
	return len(c.events)
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = c.events[:0]
}

// evict removes events older than the window duration relative to t.
// Must be called with c.mu held.
func (c *Counter) evict(t time.Time) {
	cutoff := t.Add(-c.duration)
	i := 0
	for i < len(c.events) && c.events[i].Before(cutoff) {
		i++
	}
	c.events = c.events[i:]
}
