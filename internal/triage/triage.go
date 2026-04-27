// Package triage provides prioritisation logic for lease alerts,
// ranking pending notifications by severity and age so the most
// critical leases are dispatched first.
package triage

import (
	"sort"
	"sync"
	"time"

	"github.com/youorg/vaultwatch/internal/lease"
)

// Priority represents the dispatch order weight of a lease event.
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
)

// Entry is a queued lease event with its computed priority.
type Entry struct {
	Info      lease.Info
	Priority  Priority
	QueuedAt  time.Time
}

// Queue holds pending lease entries ordered by priority then age.
type Queue struct {
	mu      sync.Mutex
	entries []Entry
}

// New returns an empty triage Queue.
func New() *Queue {
	return &Queue{}
}

// Add inserts a lease.Info into the queue, computing its priority
// from the lease status.
func (q *Queue) Add(info lease.Info) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.entries = append(q.entries, Entry{
		Info:     info,
		Priority: priorityFor(info.Status),
		QueuedAt: time.Now(),
	})
}

// Drain returns all queued entries sorted by descending priority,
// then ascending queue time, and clears the queue.
func (q *Queue) Drain() []Entry {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]Entry, len(q.entries))
	copy(out, q.entries)
	q.entries = q.entries[:0]
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority != out[j].Priority {
			return out[i].Priority > out[j].Priority
		}
		return out[i].QueuedAt.Before(out[j].QueuedAt)
	})
	return out
}

// Len returns the current number of queued entries.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.entries)
}

func priorityFor(status lease.Status) Priority {
	switch status {
	case lease.StatusCritical, lease.StatusExpired:
		return PriorityHigh
	case lease.StatusWarning:
		return PriorityMedium
	default:
		return PriorityLow
	}
}
