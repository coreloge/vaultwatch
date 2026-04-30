// Package signal provides OS signal handling for graceful shutdown.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for OS signals and cancels a context on shutdown signals.
type Handler struct {
	signals []os.Signal
}

// New returns a Handler that responds to SIGINT and SIGTERM by default.
func New(sigs ...os.Signal) *Handler {
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	return &Handler{signals: sigs}
}

// Notify returns a derived context that is cancelled when one of the
// configured OS signals is received. The returned stop function releases
// the signal channel and should be called when the context is no longer
// needed (e.g. after the main goroutine exits).
func (h *Handler) Notify(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.signals...)

	go func() {
		defer signal.Stop(ch)
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}
