// Package sampler implements probabilistic sampling for vaultwatch alert events.
//
// Sampling reduces downstream load by forwarding only a fraction of events
// that match a given set of lease statuses. Events whose status is not
// targeted by the sampler are always forwarded unchanged.
//
// Example usage:
//
//	cfg := sampler.Config{
//		Rate:     0.25,
//		Statuses: []lease.Status{lease.StatusWarning},
//	}
//	s := sampler.New(cfg, nil)
//	if s.Allow(info) {
//		// forward the alert
//	}
package sampler
