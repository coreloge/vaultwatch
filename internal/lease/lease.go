package lease

import "time"

// Status represents the urgency level of a lease expiration.
type Status int

const (
	StatusOK       Status = iota
	StatusWarning         // within warning threshold
	StatusCritical        // within critical threshold
	StatusExpired         // already expired
)

// Info holds metadata about a Vault secret lease.
type Info struct {
	LeaseID        string
	Path           string
	Renewable      bool
	TTL            time.Duration
	ExpiresAt      time.Time
	CreatedAt      time.Time
}

// Classify returns the Status of a lease given warning and critical thresholds.
func Classify(l Info, warnThreshold, critThreshold time.Duration) Status {
	now := time.Now()
	if now.After(l.ExpiresAt) {
		return StatusExpired
	}
	remaining := l.ExpiresAt.Sub(now)
	if remaining <= critThreshold {
		return StatusCritical
	}
	if remaining <= warnThreshold {
		return StatusWarning
	}
	return StatusOK
}

// Remaining returns the duration until the lease expires.
// Returns zero if already expired.
func (l Info) Remaining() time.Duration {
	d := time.Until(l.ExpiresAt)
	if d < 0 {
		return 0
	}
	return d
}

// IsExpired reports whether the lease has already expired.
func (l Info) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}
