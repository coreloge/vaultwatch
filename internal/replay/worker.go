package replay

import (
	"context"
	"log"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Dispatcher is the interface used by the Worker to re-dispatch lease events.
type Dispatcher interface {
	Dispatch(ctx context.Context, info lease.Info) error
}

// Worker periodically drains the replay store and attempts re-dispatch.
type Worker struct {
	store    *Store
	dispatch Dispatcher
	interval time.Duration
}

// NewWorker returns a Worker that drains store on each tick of interval.
func NewWorker(store *Store, d Dispatcher, interval time.Duration) *Worker {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Worker{store: store, dispatch: d, interval: interval}
}

// Run starts the replay loop. It blocks until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.flush(ctx)
		}
	}
}

func (w *Worker) flush(ctx context.Context) {
	entries := w.store.Drain()
	for _, e := range entries {
		e.Attempts++
		if err := w.dispatch.Dispatch(ctx, e.LeaseInfo); err != nil {
			log.Printf("replay: re-dispatch failed for lease %s (attempt %d): %v",
				e.LeaseInfo.LeaseID, e.Attempts, err)
			// Re-add only if context is still alive
			if ctx.Err() == nil {
				w.store.Add(e.LeaseInfo)
			}
		}
	}
}
