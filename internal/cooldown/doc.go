// Package cooldown implements a per-lease quiet period that prevents the same
// lease from generating repeated alerts within a configurable window.
//
// After an alert is dispatched for a lease, Record should be called to start
// the cooldown. Subsequent calls to Allow will return false until the period
// has elapsed, at which point alerting resumes automatically.
//
// Purge can be called periodically to reclaim memory for expired records.
package cooldown
