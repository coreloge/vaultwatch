// Package quota provides per-prefix alert dispatch rate limiting for
// vaultwatch. It tracks how many alerts have been sent for a given lease
// path prefix within a rolling time window and rejects further dispatches
// once the configured maximum is reached.
//
// Usage:
//
//	q := quota.New(quota.DefaultConfig())
//	if q.Allow("secret/prod/") {
//		// dispatch alert
//	}
package quota
