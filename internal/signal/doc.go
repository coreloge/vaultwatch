// Package signal wraps os/signal to provide context-based graceful shutdown
// for the vaultwatch daemon.
//
// Usage:
//
//	h := signal.New() // defaults to SIGINT and SIGTERM
//	ctx, stop := h.Notify(context.Background())
//	defer stop()
//
//	// pass ctx to long-running components; they will unwind when a
//	// shutdown signal is received.
package signal
