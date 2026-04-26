// Package schedule provides a lightweight, context-aware Scheduler
// that drives periodic lease-check cycles inside the vaultwatch daemon.
//
// Usage:
//
//	s := schedule.New(30 * time.Second)
//	go s.Run(ctx, func(ctx context.Context) {
//		monitor.CheckAll(ctx)
//	})
//
// The callback is invoked once immediately and then on every tick.
// Cancelling the context or calling Stop() both cleanly terminate the
// loop without leaking goroutines.
package schedule
