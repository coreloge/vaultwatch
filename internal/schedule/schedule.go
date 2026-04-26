// Package schedule provides interval-based ticker utilities for
// driving periodic lease checks within the vaultwatch daemon.
package schedule

import (
	"context"
	"time"
)

// Scheduler triggers a callback at a fixed interval until the context
// is cancelled or Stop is called.
type Scheduler struct {
	interval time.Duration
	stop     chan struct{}
}

// New returns a Scheduler that fires every interval.
func New(interval time.Duration) *Scheduler {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Scheduler{
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Run calls fn immediately and then on every tick until ctx is done
// or Stop is called. It blocks until the scheduler exits.
func (s *Scheduler) Run(ctx context.Context, fn func(context.Context)) {
	fn(ctx)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fn(ctx)
		case <-s.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

// Stop signals the scheduler to cease firing after the current
// invocation (if any) completes.
func (s *Scheduler) Stop() {
	select {
	case s.stop <- struct{}{}:
	default:
	}
}

// Interval returns the configured tick interval.
func (s *Scheduler) Interval() time.Duration {
	return s.interval
}
