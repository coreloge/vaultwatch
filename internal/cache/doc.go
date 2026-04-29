// Package cache provides a generic TTL-based in-memory cache used by
// vaultwatch to store transient lease metadata between polling cycles.
//
// The cache is safe for concurrent use. Expired entries are not evicted
// automatically; callers should invoke Purge periodically or rely on the
// Get method which skips expired entries transparently.
package cache
