// Package circuitbreaker implements a simple circuit breaker for webhook delivery.
// It tracks consecutive failures and opens the circuit to prevent cascading errors.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing, requests blocked
	StateHalfOpen              // testing if backend recovered
)

// ErrCircuitOpen is returned when the circuit is open and requests are blocked.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// Config holds the configuration for a CircuitBreaker.
type Config struct {
	MaxFailures  int
	OpenDuration time.Duration
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		MaxFailures:  5,
		OpenDuration: 30 * time.Second,
	}
}

// CircuitBreaker tracks failures for a named endpoint and opens the circuit
// after MaxFailures consecutive failures.
type CircuitBreaker struct {
	mu           sync.Mutex
	cfg          Config
	state        State
	failures     int
	openedAt     time.Time
}

// New creates a new CircuitBreaker with the given config.
func New(cfg Config) *CircuitBreaker {
	return &CircuitBreaker{cfg: cfg}
}

// Allow returns nil if the request is permitted, or ErrCircuitOpen if blocked.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Since(cb.openedAt) >= cb.cfg.OpenDuration {
			cb.state = StateHalfOpen
			return nil
		}
		return ErrCircuitOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess resets the failure count and closes the circuit.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = StateClosed
}

// RecordFailure increments the failure count and may open the circuit.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.state == StateHalfOpen || cb.failures >= cb.cfg.MaxFailures {
		cb.state = StateOpen
		cb.openedAt = time.Now()
	}
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) CurrentState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
