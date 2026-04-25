package notify

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/lease"
	"github.com/yourusername/vaultwatch/internal/suppress"
	"github.com/yourusername/vaultwatch/internal/webhook"
)

// Dispatcher sends alert payloads to a configured webhook endpoint.
type Dispatcher struct {
	hook       *webhook.Webhook
	builder    *alert.Builder
	suppressor *suppress.Suppressor
}

// New creates a Dispatcher targeting the given webhook URL.
// suppressDuration controls how long a lease alert is silenced after firing.
func New(webhookURL, secret string, suppressDuration time.Duration) (*Dispatcher, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("notify: webhook URL must not be empty")
	}
	hook, err := webhook.New(webhookURL, secret)
	if err != nil {
		return nil, fmt.Errorf("notify: %w", err)
	}
	return &Dispatcher{
		hook:       hook,
		builder:    alert.NewBuilder(),
		suppressor: suppress.New(suppressDuration),
	}, nil
}

// Dispatch builds an alert payload and sends it via the webhook.
// If the lease is currently suppressed, the dispatch is skipped.
func (d *Dispatcher) Dispatch(ctx context.Context, info lease.Info) error {
	if d.suppressor.IsSuppressed(info.LeaseID) {
		log.Printf("notify: suppressed alert for lease %s", info.LeaseID)
		return nil
	}

	payload, err := d.builder.Build(info)
	if err != nil {
		return fmt.Errorf("notify: build alert: %w", err)
	}

	if err := d.hook.Send(ctx, payload); err != nil {
		return fmt.Errorf("notify: send webhook: %w", err)
	}

	d.suppressor.Suppress(info.LeaseID)
	return nil
}

// Release clears any active suppression for the given lease ID,
// allowing the next alert to be dispatched immediately.
func (d *Dispatcher) Release(leaseID string) {
	d.suppressor.Release(leaseID)
}
