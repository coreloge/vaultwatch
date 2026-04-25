// Package snapshot provides functionality for capturing and comparing
// lease state over time. It allows vaultwatch to detect changes in
// lease status between monitor cycles, enabling targeted alerting
// only when lease conditions change rather than on every check.
//
// The Store maintains an in-memory map of lease snapshots keyed by
// lease ID. Callers can use Compare to determine whether a lease
// has transitioned to a new classification state.
package snapshot
