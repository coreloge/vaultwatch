// Package ratelimit implements a per-key cooldown limiter used by the
// vaultwatch dispatcher to suppress duplicate webhook alerts for the same
// Vault lease within a configurable time window.
//
// Usage:
//
//	limiter := ratelimit.New(5 * time.Minute)
//
//	if limiter.Allow(leaseID) {
//		// dispatch alert
//	}
//
// Purge should be called periodically (e.g. from the monitor loop) to
// reclaim memory for leases that are no longer being tracked.
package ratelimit
