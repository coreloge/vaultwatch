// Package lease provides types and utilities for tracking and classifying
// HashiCorp Vault secret lease expirations.
//
// The Info type represents a single lease with its metadata including TTL,
// expiry time, and path. Use Classify to determine the urgency level of a
// lease given configurable warning and critical thresholds.
//
// The Store type provides a thread-safe in-memory registry for tracking
// multiple leases concurrently, supporting set, get, delete, and enumeration
// operations.
package lease
