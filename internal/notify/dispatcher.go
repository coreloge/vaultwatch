package notify

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/example/vaultwatch/internal/alert"
	"github.com/example/vaultwatch/internal/lease"
	"github.com/example/vaultwatch/internal/webhook"
)

// Dispatcher routes lease alerts to configured webhook targets.
type Dispatcher struct {
	webhook   *webhook.Sender
	builder   *alert.Builder
	formatter *alert.Formatter
	logger    *slog.Logger
}

// Config holds dispatcher configuration.
type Config struct {
	WebhookURL    string
	WebhookSecret string
	Format        string
}

// New creates a Dispatcher with the given configuration.
func New(cfg Config, logger *slog.Logger) (*Dispatcher, error) {
	if cfg.WebhookURL == "" {
		return nil, fmt.Errorf("webhook URL must not be empty")
	}

	sender := webhook.New(cfg.WebhookURL, cfg.WebhookSecret)
	builder := alert.NewBuilder()
	formatter := alert.NewFormatter(cfg.Format)

	return &Dispatcher{
		webhook:   sender,
		builder:   builder,
		formatter: formatter,
		logger:    logger,
	}, nil
}

// Dispatch builds and sends an alert for the given lease info.
func (d *Dispatcher) Dispatch(ctx context.Context, info lease.Info) error {
	payload := d.builder.Build(info)

	body, err := d.formatter.Format(payload)
	if err != nil {
		return fmt.Errorf("formatting alert payload: %w", err)
	}

	if err := d.webhook.Send(ctx, body); err != nil {
		d.logger.Error("failed to dispatch alert",
			"lease_id", info.LeaseID,
			"severity", payload.Severity,
			"error", err,
		)
		return fmt.Errorf("sending webhook: %w", err)
	}

	d.logger.Info("alert dispatched",
		"lease_id", info.LeaseID,
		"severity", payload.Severity,
		"path", info.Path,
	)
	return nil
}
