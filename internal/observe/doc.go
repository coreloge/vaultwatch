// Package observe implements a thread-safe event observer for lease lifecycle
// notifications within vaultwatch.
//
// Multiple handlers can be registered and will each receive every emitted
// lease.Info event in the order they were registered. The observer is safe
// for concurrent use; Register, Emit, Len, and Reset may be called from
// multiple goroutines simultaneously.
package observe
