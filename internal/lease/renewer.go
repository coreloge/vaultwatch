package lease

import (
	"context"
	"log"
	"time"
)

// RenewFunc is a function that attempts to renew a lease by ID.
// It returns the new TTL on success, or an error.
type RenewFunc func(ctx context.Context, leaseID string) (time.Duration, error)

// Renewer attempts to renew leases that are approaching expiration.
type Renewer struct {
	store   *Store
	renewFn RenewFunc
	threshold time.Duration
	logger  *log.Logger
}

// NewRenewer creates a Renewer that will attempt renewal for leases
// whose remaining TTL is below the given threshold.
func NewRenewer(store *Store, renewFn RenewFunc, threshold time.Duration, logger *log.Logger) *Renewer {
	return &Renewer{
		store:     store,
		renewFn:   renewFn,
		threshold: threshold,
		logger:    logger,
	}
}

// RenewEligible iterates all tracked leases and attempts renewal for any
// whose remaining TTL is at or below the configured threshold.
// It returns the number of successful renewals.
func (r *Renewer) RenewEligible(ctx context.Context) int {
	leases := r.store.All()
	successful := 0

	for _, info := range leases {
		remaining := time.Until(info.ExpireTime)
		if remaining > r.threshold {
			continue
		}

		newTTL, err := r.renewFn(ctx, info.LeaseID)
		if err != nil {
			r.logger.Printf("[renewer] failed to renew lease %s: %v", info.LeaseID, err)
			continue
		}

		updated := info
		updated.ExpireTime = time.Now().Add(newTTL)
		updated.TTL = newTTL
		r.store.Set(updated)

		r.logger.Printf("[renewer] renewed lease %s, new TTL: %s", info.LeaseID, newTTL)
		successful++
	}

	return successful
}
