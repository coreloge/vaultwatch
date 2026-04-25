// Package suppress implements alert suppression windows for VaultWatch.
//
// A Suppressor prevents duplicate alerts from being dispatched for the same
// lease within a configurable cooldown duration. Once a lease alert is sent,
// the lease ID is suppressed until the window expires or is explicitly released.
//
// Typical usage:
//
//	suppressor := suppress.New(15 * time.Minute)
//	if !suppressor.IsSuppressed(leaseID) {
//		dispatcher.Dispatch(ctx, info)
//		suppressor.Suppress(leaseID)
//	}
package suppress
