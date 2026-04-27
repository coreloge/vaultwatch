package lease

import (
	"context"
	"log"
	"time"
)

// WatchHandler is called when a lease event is detected during watching.
type WatchHandler func(info Info)

// Watcher polls the lease store and invokes handlers when leases are
// approaching expiry or have changed status.
type Watcher struct {
	store    *Store
	checker  *ExpiryChecker
	interval time.Duration
	onExpiry WatchHandler
	log      *log.Logger
}

// WatcherConfig holds configuration for a Watcher.
type WatcherConfig struct {
	Store    *Store
	Checker  *ExpiryChecker
	Interval time.Duration
	OnExpiry WatchHandler
	Logger   *log.Logger
}

// NewWatcher creates a Watcher with the provided configuration.
// If Interval is zero it defaults to 30 seconds.
func NewWatcher(cfg WatcherConfig) *Watcher {
	if cfg.Interval <= 0 {
		cfg.Interval = 30 * time.Second
	}
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}
	return &Watcher{
		store:    cfg.Store,
		checker:  cfg.Checker,
		interval: cfg.Interval,
		onExpiry: cfg.OnExpiry,
		log:      cfg.Logger,
	}
}

// Run starts the watch loop. It returns when ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.scan()

	for {
		select {
		case <-ctx.Done():
			w.log.Println("lease watcher: stopping")
			return
		case <-ticker.C:
			w.scan()
		}
	}
}

// scan iterates all leases and fires the handler for actionable ones.
func (w *Watcher) scan() {
	now := time.Now()
	for _, info := range w.store.All() {
		if w.checker.IsCritical(info, now) || w.checker.IsWarning(info, now) {
			if w.onExpiry != nil {
				w.onExpiry(info)
			}
		}
	}
}
