// Package replay provides a persistent-in-memory store for lease alert events
// that could not be dispatched at the time they were generated.
//
// When a webhook call fails and the circuit breaker is open, or when the
// dispatcher is temporarily unavailable, events can be added to the replay
// store. A background worker can periodically drain the store and attempt
// re-dispatch, ensuring no critical lease expiry notifications are silently
// lost.
//
// Entries older than the configured maxAge are automatically discarded on
// the next Drain call to prevent unbounded growth.
package replay
