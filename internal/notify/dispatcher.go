package notify

import (
	"context"
	"fmt"
	"log"

	"github.com/vaultwatch/internal/alert"
	"github.com/vaultwatch/internal/lease"
	"github.com/vaultwatch/internal/throttle"
	"github.com/vaultwatch/internal/webhook"
)

// Dispatcher sends alert payloads to a configured webhook endpoint.
type Dispatcher struct {
	sender   *webhook.Sender
	builder  *alert.Builder
	throttle *throttle.Throttler
	url      string
}

// New creates a Dispatcher targeting the given webhook URL.
// throttleWindow controls the minimum interval between repeated alerts
// for the same lease ID (pass 0 to disable throttling).
func New(url string, th *throttle.Throttler) (*Dispatcher, error) {
	if url == "" {
		return nil, fmt.Errorf("notify: webhook URL must not be empty")
	}
	s, err := webhook.New(url, "")
	if err != nil {
		return nil, fmt.Errorf("notify: %w", err)
	}
	return &Dispatcher{
		sender:  s,
		builder: alert.NewBuilder("vaultwatch"),
		throttle: th,
		url:     url,
	}, nil
}

// Dispatch builds and sends an alert for the given lease info.
// If the throttle rejects the lease ID, the call is a no-op.
func (d *Dispatcher) Dispatch(ctx context.Context, info lease.Info) error {
	if d.throttle != nil && !d.throttle.Allow(info.LeaseID) {
		log.Printf("notify: throttled alert for lease %s", info.LeaseID)
		return nil
	}

	payload, err := d.builder.Build(info)
	if err != nil {
		return fmt.Errorf("notify: build payload: %w", err)
	}

	if err := d.sender.Send(ctx, payload); err != nil {
		return fmt.Errorf("notify: send: %w", err)
	}
	return nil
}
