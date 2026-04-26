// Package pipeline wires together the lease processing stages:
// filter → dedup → suppress → throttle → notify.
package pipeline

import (
	"context"
	"log"

	"github.com/yourusername/vaultwatch/internal/dedup"
	"github.com/yourusername/vaultwatch/internal/filter"
	"github.com/yourusername/vaultwatch/internal/lease"
	"github.com/yourusername/vaultwatch/internal/notify"
	"github.com/yourusername/vaultwatch/internal/suppress"
	"github.com/yourusername/vaultwatch/internal/throttle"
)

// Pipeline processes a LeaseInfo through a series of gates before
// dispatching an alert via the notifier.
type Pipeline struct {
	filter    *filter.Filter
	dedup     *dedup.Deduplicator
	suppress  *suppress.Suppressor
	throttle  *throttle.Throttler
	dispatch  *notify.Dispatcher
}

// Config holds the dependencies required to build a Pipeline.
type Config struct {
	Filter   *filter.Filter
	Dedup    *dedup.Deduplicator
	Suppress *suppress.Suppressor
	Throttle *throttle.Throttler
	Dispatch *notify.Dispatcher
}

// New constructs a Pipeline from the provided Config.
func New(cfg Config) *Pipeline {
	return &Pipeline{
		filter:   cfg.Filter,
		dedup:    cfg.Dedup,
		suppress: cfg.Suppress,
		throttle: cfg.Throttle,
		dispatch: cfg.Dispatch,
	}
}

// Process runs info through each stage. It returns true if an alert
// was dispatched, false if any stage dropped the event.
func (p *Pipeline) Process(ctx context.Context, info lease.Info) bool {
	if !p.filter.Allow(info) {
		log.Printf("[pipeline] filtered out lease %s", info.LeaseID)
		return false
	}

	if p.dedup.IsDuplicate(info) {
		log.Printf("[pipeline] duplicate lease %s, skipping", info.LeaseID)
		return false
	}

	if p.suppress.IsSuppressed(info.LeaseID) {
		log.Printf("[pipeline] suppressed lease %s, skipping", info.LeaseID)
		return false
	}

	if !p.throttle.Allow(info.LeaseID) {
		log.Printf("[pipeline] throttled lease %s, skipping", info.LeaseID)
		return false
	}

	if err := p.dispatch.Dispatch(ctx, info); err != nil {
		log.Printf("[pipeline] dispatch error for lease %s: %v", info.LeaseID, err)
		return false
	}

	p.dedup.Record(info)
	return true
}
