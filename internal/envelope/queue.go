package envelope

import (
	"sync"
)

// Queue is a thread-safe FIFO queue of Envelopes.
type Queue struct {
	mu    sync.Mutex
	items []*Envelope
}

// NewQueue returns an initialised, empty Queue.
func NewQueue() *Queue {
	return &Queue{}
}

// Push appends an envelope to the back of the queue.
func (q *Queue) Push(e *Envelope) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, e)
}

// Pop removes and returns the front envelope, or nil if empty.
func (q *Queue) Pop() *Envelope {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	e := q.items[0]
	q.items = q.items[1:]
	return e
}

// Len returns the current number of envelopes in the queue.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Drain removes and returns all envelopes currently in the queue.
func (q *Queue) Drain() []*Envelope {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]*Envelope, len(q.items))
	copy(out, q.items)
	q.items = q.items[:0]
	return out
}
