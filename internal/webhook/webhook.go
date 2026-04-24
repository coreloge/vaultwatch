package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the JSON body sent to a webhook endpoint.
type Payload struct {
	LeaseID   string    `json:"lease_id"`
	ExpiresAt time.Time `json:"expires_at"`
	TTL       int64     `json:"ttl_seconds"`
	Message   string    `json:"message"`
}

// Sender sends webhook notifications.
type Sender struct {
	client  *http.Client
	url     string
	secret  string
}

// New creates a new Sender with the given webhook URL and optional signing secret.
func New(url, secret string, timeout time.Duration) *Sender {
	return &Sender{
		client: &http.Client{Timeout: timeout},
		url:    url,
		secret: secret,
	}
}

// Send marshals the payload and POSTs it to the configured webhook URL.
func (s *Sender) Send(ctx context.Context, p Payload) error {
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if s.secret != "" {
		req.Header.Set("X-VaultWatch-Secret", s.secret)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, s.url)
	}
	return nil
}
