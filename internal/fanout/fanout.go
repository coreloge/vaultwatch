// Package fanout distributes a single lease alert to multiple downstream handlers concurrently.
package fanout

import (
	"context"
	"sync"

	"github.com/youorg/vaultwatch/internal/lease"
)

// Handler is any component that can receive a lease info event.
type Handler interface {
	Handle(ctx context.Context, info lease.Info) error
}

// Fanout sends each event to all registered handlers concurrently and
// collects any errors that occur.
type Fanout struct {
	handlers []Handler
}

// New returns a Fanout that will broadcast to the provided handlers.
func New(handlers ...Handler) *Fanout {
	return &Fanout{handlers: handlers}
}

// Send broadcasts info to every handler concurrently and returns a combined
// slice of all non-nil errors. A nil return means all handlers succeeded.
func (f *Fanout) Send(ctx context.Context, info lease.Info) []error {
	if len(f.handlers) == 0 {
		return nil
	}

	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for _, h := range f.handlers {
		wg.Add(1)
		go func(h Handler) {
			defer wg.Done()
			if err := h.Handle(ctx, info); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(h)
	}

	wg.Wait()
	return errs
}

// Add appends a handler to the fanout at runtime.
func (f *Fanout) Add(h Handler) {
	f.handlers = append(f.handlers, h)
}

// Len returns the number of registered handlers.
func (f *Fanout) Len() int {
	return len(f.handlers)
}
