package monitor

import (
	"context"
	"log"
	"time"

	"github.com/vaultwatch/internal/config"
	"github.com/vaultwatch/internal/vault"
)

// LeaseStatus represents the current state of a monitored lease.
type LeaseStatus struct {
	LeaseID   string
	ExpiresAt time.Time
	TTL       time.Duration
	Expiring  bool
}

// Monitor polls Vault for lease expiration and triggers alerts.
type Monitor struct {
	cfg    *config.Config
	client *vault.Client
}

// New creates a new Monitor instance.
func New(cfg *config.Config, client *vault.Client) *Monitor {
	return &Monitor{
		cfg:    cfg,
		client: client,
	}
}

// Run starts the monitoring loop, blocking until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.cfg.PollInterval)
	defer ticker.Stop()

	log.Printf("monitor: starting poll loop (interval=%s, warning_threshold=%s)",
		m.cfg.PollInterval, m.cfg.WarningThreshold)

	for {
		select {
		case <-ctx.Done():
			log.Println("monitor: shutting down")
			return ctx.Err()
		case <-ticker.C:
			if err := m.poll(ctx); err != nil {
				log.Printf("monitor: poll error: %v", err)
			}
		}
	}
}

// poll checks each configured lease and evaluates expiration.
func (m *Monitor) poll(ctx context.Context) error {
	for _, leaseID := range m.cfg.LeaseIDs {
		status, err := m.checkLease(ctx, leaseID)
		if err != nil {
			log.Printf("monitor: failed to check lease %s: %v", leaseID, err)
			continue
		}
		if status.Expiring {
			log.Printf("monitor: lease %s expiring in %s", leaseID, status.TTL.Round(time.Second))
		}
	}
	return nil
}

// checkLease looks up a lease and determines whether it is near expiration.
func (m *Monitor) checkLease(ctx context.Context, leaseID string) (*LeaseStatus, error) {
	lease, err := m.client.LookupLease(ctx, leaseID)
	if err != nil {
		return nil, err
	}

	ttl := time.Duration(lease.TTL) * time.Second
	expiresAt := time.Now().Add(ttl)

	return &LeaseStatus{
		LeaseID:   leaseID,
		ExpiresAt: expiresAt,
		TTL:       ttl,
		Expiring:  ttl <= m.cfg.WarningThreshold,
	}, nil
}
