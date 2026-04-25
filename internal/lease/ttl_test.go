package lease

import (
	"testing"
	"time"
)

func TestNewTTLFromSeconds(t *testing.T) {
	ttl := NewTTLFromSeconds(3600)
	if ttl.Seconds() != 3600 {
		t.Errorf("expected 3600 seconds, got %d", ttl.Seconds())
	}
	if ttl.Duration() != time.Hour {
		t.Errorf("expected 1h duration, got %s", ttl.Duration())
	}
}

func TestTTL_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		ttl      TTL
		expected bool
	}{
		{"positive", NewTTLFromSeconds(60), false},
		{"zero", NewTTLFromSeconds(0), true},
		{"negative", NewTTL(-time.Second), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.ttl.IsZero() != tc.expected {
				t.Errorf("IsZero() = %v, want %v", tc.ttl.IsZero(), tc.expected)
			}
		})
	}
}

func TestTTL_ExpiresAt(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	ttl := NewTTLFromSeconds(3600)
	expected := now.Add(time.Hour)
	if !ttl.ExpiresAt(now).Equal(expected) {
		t.Errorf("ExpiresAt() = %v, want %v", ttl.ExpiresAt(now), expected)
	}
}

func TestTTL_String(t *testing.T) {
	tests := []struct {
		ttl      TTL
		expected string
	}{
		{NewTTLFromSeconds(0), "expired"},
		{NewTTL(-time.Minute), "expired"},
		{NewTTLFromSeconds(45), "45s"},
		{NewTTLFromSeconds(125), "2m5s"},
		{NewTTLFromSeconds(3661), "1h1m1s"},
	}
	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			if got := tc.ttl.String(); got != tc.expected {
				t.Errorf("String() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestTTL_RemainingFrom_Clamped(t *testing.T) {
	// A TTL that expired in the past should return zero remaining.
	past := time.Now().Add(-2 * time.Hour)
	ttl := NewTTLFromSeconds(3600) // 1h TTL set from 2h ago
	remaining := ttl.RemainingFrom(past)
	if !remaining.IsZero() {
		t.Errorf("expected zero remaining for expired lease, got %s", remaining)
	}
}
