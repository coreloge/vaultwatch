// Package cooldown provides a per-lease cooldown mechanism that prevents
// re-alerting within a configurable quiet period after an alert has been sent.
package cooldown

import (
	"sync"
	"time"
)

// DefaultCooldown is the default quiet period between alerts for the same lease.
const DefaultCooldown = 5 * time.Minute

// Cooldown tracks the last alert time per lease ID and enforces a minimum
// interval before the same lease may trigger another alert.
type Cooldown struct {
	mu       sync.Mutex
	records  map[string]time.Time
	period   time.Duration
	now      func() time.Time
}

// New returns a Cooldown with the given quiet period. If period is zero,
// DefaultCooldown is used.
func New(period time.Duration) *Cooldown {
	if period <= 0 {
		period = DefaultCooldown
	}
	return &Cooldown{
		records: make(map[string]time.Time),
		period:  period,
		now:     time.Now,
	}
}

// Allow reports whether leaseID may trigger an alert. It returns true if no
// alert has been recorded for the lease, or if the cooldown period has elapsed
// since the last alert.
func (c *Cooldown) Allow(leaseID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	last, ok := c.records[leaseID]
	if !ok {
		return true
	}
	return c.now().Sub(last) >= c.period
}

// Record marks leaseID as having just triggered an alert, resetting its
// cooldown window.
func (c *Cooldown) Record(leaseID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.records[leaseID] = c.now()
}

// Reset removes the cooldown record for leaseID, allowing it to alert
// immediately on the next check.
func (c *Cooldown) Reset(leaseID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.records, leaseID)
}

// Purge removes all records whose cooldown period has already elapsed,
// freeing memory for leases that are no longer active.
func (c *Cooldown) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	for id, last := range c.records {
		if now.Sub(last) >= c.period {
			delete(c.records, id)
		}
	}
}
