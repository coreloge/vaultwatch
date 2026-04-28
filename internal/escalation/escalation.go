// Package escalation provides multi-tier alert escalation based on
// how long a lease has remained in a critical or warning state.
package escalation

import (
	"sync"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Tier represents an escalation level.
type Tier int

const (
	TierNone     Tier = iota
	TierWarning        // first escalation
	TierCritical       // second escalation
	TierEmergency      // final escalation
)

// Config holds thresholds for each escalation tier.
type Config struct {
	WarningAfter   time.Duration // time in bad state before Warning tier
	CriticalAfter  time.Duration // time in bad state before Critical tier
	EmergencyAfter time.Duration // time in bad state before Emergency tier
}

// DefaultConfig returns sensible default escalation thresholds.
func DefaultConfig() Config {
	return Config{
		WarningAfter:   5 * time.Minute,
		CriticalAfter:  15 * time.Minute,
		EmergencyAfter: 30 * time.Minute,
	}
}

type entry struct {
	firstSeen time.Time
	status    lease.Status
}

// Escalator tracks how long each lease has been in a degraded state
// and returns the appropriate escalation tier.
type Escalator struct {
	mu      sync.Mutex
	cfg     Config
	entries map[string]entry
}

// New creates a new Escalator with the given Config.
func New(cfg Config) *Escalator {
	return &Escalator{
		cfg:     cfg,
		entries: make(map[string]entry),
	}
}

// Evaluate returns the current escalation Tier for the given lease.
// If the lease is healthy, any tracked state for it is cleared.
func (e *Escalator) Evaluate(info lease.Info) Tier {
	e.mu.Lock()
	defer e.mu.Unlock()

	if info.Status == lease.StatusOK {
		delete(e.entries, info.LeaseID)
		return TierNone
	}

	ent, ok := e.entries[info.LeaseID]
	if !ok || ent.status != info.Status {
		ent = entry{firstSeen: time.Now(), status: info.Status}
		e.entries[info.LeaseID] = ent
	}

	duration := time.Since(ent.firstSeen)
	switch {
	case duration >= e.cfg.EmergencyAfter:
		return TierEmergency
	case duration >= e.cfg.CriticalAfter:
		return TierCritical
	case duration >= e.cfg.WarningAfter:
		return TierWarning
	default:
		return TierNone
	}
}

// Reset clears tracked state for a specific lease.
func (e *Escalator) Reset(leaseID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.entries, leaseID)
}
