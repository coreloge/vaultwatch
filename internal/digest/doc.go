// Package digest provides content-based fingerprinting for lease alert payloads.
//
// A Digester produces short, stable hex strings from the fields that
// uniquely identify a lease event (lease ID, status, expiry minute).
// These fingerprints are used by deduplication and caching layers to
// avoid re-alerting on equivalent events within a time window.
//
// Usage:
//
//	d := digest.New(16) // 16-character truncated digest
//	fp := d.Compute(info)
package digest
