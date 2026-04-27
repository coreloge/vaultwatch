// Package limiter implements a semaphore-based concurrency limiter for
// VaultWatch's webhook dispatch pipeline.
//
// It prevents thundering-herd conditions when many leases expire
// simultaneously by capping the number of in-flight webhook calls.
//
// Usage:
//
//	l, _ := limiter.New(10)
//	if err := l.Acquire(ctx); err != nil {
//	    return err
//	}
//	defer l.Release()
//	// ... perform webhook dispatch
package limiter
