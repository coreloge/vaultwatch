package sink

import (
	"context"

	"github.com/youorg/vaultwatch/internal/alert"
	"github.com/youorg/vaultwatch/internal/webhook"
)

// WebhookTarget wraps a webhook.Sender so it satisfies the Target interface.
type WebhookTarget struct {
	name   string
	sender *webhook.Sender
}

// NewWebhookTarget creates a Target that delivers alerts via the given
// webhook.Sender. name is used in error messages to identify the target.
func NewWebhookTarget(name string, sender *webhook.Sender) *WebhookTarget {
	return &WebhookTarget{name: name, sender: sender}
}

// Name returns the human-readable identifier for this target.
func (w *WebhookTarget) Name() string { return w.name }

// Send marshals the payload and posts it to the configured webhook URL.
func (w *WebhookTarget) Send(ctx context.Context, p alert.Payload) error {
	return w.sender.Send(ctx, p)
}
