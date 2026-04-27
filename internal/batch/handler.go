package batch

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
)

// LoggingHandler returns a Handler that logs each flushed batch using the
// standard logger. It is intended as a lightweight default or debug handler.
func LoggingHandler(logger *log.Logger) Handler {
	if logger == nil {
		logger = log.Default()
	}
	return func(_ context.Context, items []lease.Info) {
		logger.Printf("batch: flushed %d lease event(s) at %s", len(items), time.Now().UTC().Format(time.RFC3339))
	}
}

// MultiHandler returns a Handler that fans out to all provided handlers in
// sequence. If a handler panics the remaining handlers are still called.
func MultiHandler(handlers ...Handler) Handler {
	return func(ctx context.Context, items []lease.Info) {
		for _, h := range handlers {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("batch: handler panic recovered: %v", r)
					}
				}()
				h(ctx, items)
			}()
		}
	}
}
