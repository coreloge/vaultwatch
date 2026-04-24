package webhook

import (
	"context"
	"fmt"
	"time"
)

// RetryConfig holds parameters for retry behaviour.
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       2 * time.Second,
	}
}

// SendWithRetry attempts to send the payload up to cfg.MaxAttempts times,
// sleeping cfg.Delay between failures. It returns the last error on exhaustion.
func (s *Sender) SendWithRetry(ctx context.Context, p Payload, cfg RetryConfig) error {
	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("webhook: context cancelled before attempt %d: %w", attempt, err)
		}
		if err := s.Send(ctx, p); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if attempt < cfg.MaxAttempts {
			select {
			case <-time.After(cfg.Delay):
			case <-ctx.Done():
				return fmt.Errorf("webhook: context cancelled during retry delay: %w", ctx.Err())
			}
		}
	}
	return fmt.Errorf("webhook: all %d attempts failed, last error: %w", cfg.MaxAttempts, lastErr)
}
