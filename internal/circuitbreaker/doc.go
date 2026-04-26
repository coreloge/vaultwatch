// Package circuitbreaker provides a thread-safe circuit breaker for protecting
// outbound webhook calls in vaultwatch.
//
// The circuit breaker operates in three states:
//
//   - Closed: normal operation; all requests are permitted.
//   - Open: the failure threshold has been exceeded; requests are blocked
//     and ErrCircuitOpen is returned immediately.
//   - Half-Open: after the open duration elapses, one request is allowed
//     through as a probe. A success closes the circuit; a failure re-opens it.
//
// Usage:
//
//	cb := circuitbreaker.New(circuitbreaker.DefaultConfig())
//	if err := cb.Allow(); err != nil {
//	    // skip delivery
//	}
//	// ... attempt delivery ...
//	cb.RecordSuccess() // or cb.RecordFailure()
package circuitbreaker
