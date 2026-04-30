package fanout

import (
	"context"
	"fmt"

	"github.com/youorg/vaultwatch/internal/lease"
)

// HandlerFunc is a function adapter that satisfies the Handler interface.
type HandlerFunc func(ctx context.Context, info lease.Info) error

// Handle calls the underlying function.
func (f HandlerFunc) Handle(ctx context.Context, info lease.Info) error {
	return f(ctx, info)
}

// LoggingHandler wraps a Handler and logs the lease ID and any error to the
// provided printf-compatible function before returning the original error.
func LoggingHandler(inner Handler, logf func(format string, args ...any)) Handler {
	return HandlerFunc(func(ctx context.Context, info lease.Info) error {
		err := inner.Handle(ctx, info)
		if err != nil {
			logf("fanout: handler error for lease %s: %v", info.LeaseID, err)
		} else {
			logf("fanout: delivered lease %s (%s)", info.LeaseID, info.Status)
		}
		return err
	})
}

// NoopHandler returns a Handler that always succeeds without doing anything.
// Useful as a placeholder in tests or disabled pipeline stages.
func NoopHandler() Handler {
	return HandlerFunc(func(_ context.Context, _ lease.Info) error {
		return nil
	})
}

// ErrorHandler returns a Handler that always returns the given error.
// Useful for testing fanout error-collection behaviour.
func ErrorHandler(err error) Handler {
	return HandlerFunc(func(_ context.Context, _ lease.Info) error {
		return fmt.Errorf("errorHandler: %w", err)
	})
}
