// Package sink provides a fan-out delivery layer that writes alert
// payloads to one or more registered output targets (webhooks, loggers,
// etc.) and collects per-target errors without short-circuiting delivery.
package sink

import (
	"context"
	"fmt"
	"strings"

	"github.com/youorg/vaultwatch/internal/alert"
)

// Target is anything that can receive an alert payload.
type Target interface {
	Name() string
	Send(ctx context.Context, p alert.Payload) error
}

// Sink fans an alert payload out to every registered Target.
type Sink struct {
	targets []Target
}

// New returns a Sink that will deliver to the provided targets.
func New(targets ...Target) *Sink {
	return &Sink{targets: targets}
}

// SendAll delivers p to every target. Errors from individual targets
// are collected and returned as a single combined error; successful
// targets are not affected by failures in others.
func (s *Sink) SendAll(ctx context.Context, p alert.Payload) error {
	var errs []string
	for _, t := range s.targets {
		if err := t.Send(ctx, p); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", t.Name(), err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("sink errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// Len returns the number of registered targets.
func (s *Sink) Len() int { return len(s.targets) }
