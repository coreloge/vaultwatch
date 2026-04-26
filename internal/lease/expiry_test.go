package lease

import (
	"testing"
	"time"
)

func newChecker(criticalHours, warningHours float64) *ExpiryChecker {
	w := ExpiryWindow{
		Critical: time.Duration(criticalHours * float64(time.Hour)),
		Warning:  time.Duration(warningHours * float64(time.Hour)),
	}
	return NewExpiryChecker(w)
}

func TestIsCritical_WithinThreshold(t *testing.T) {
	c := newChecker(1, 6)
	ttl := NewTTLFromSeconds(1800) // 30 minutes
	if !c.IsCritical(ttl) {
		t.Error("expected IsCritical to be true for 30-minute TTL")
	}
}

func TestIsCritical_OutsideThreshold(t *testing.T) {
	c := newChecker(1, 6)
	ttl := NewTTLFromSeconds(7200) // 2 hours
	if c.IsCritical(ttl) {
		t.Error("expected IsCritical to be false for 2-hour TTL")
	}
}

func TestIsWarning_WithinThreshold(t *testing.T) {
	c := newChecker(1, 6)
	ttl := NewTTLFromSeconds(10800) // 3 hours
	if !c.IsWarning(ttl) {
		t.Error("expected IsWarning to be true for 3-hour TTL")
	}
}

func TestIsWarning_OutsideThreshold(t *testing.T) {
	c := newChecker(1, 6)
	ttl := NewTTLFromSeconds(25200) // 7 hours
	if c.IsWarning(ttl) {
		t.Error("expected IsWarning to be false for 7-hour TTL")
	}
}

func TestIsExpired_Negative(t *testing.T) {
	c := newChecker(1, 6)
	ttl := NewTTLFromSeconds(-1)
	if !c.IsExpired(ttl) {
		t.Error("expected IsExpired to be true for negative TTL")
	}
}

func TestStatusFor_ReturnsCorrectStatus(t *testing.T) {
	c := newChecker(1, 6)
	cases := []struct {
		secs     int
		want     Status
	}{
		{-1, StatusExpired},
		{1800, StatusCritical},
		{10800, StatusWarning},
		{86400, StatusOK},
	}
	for _, tc := range cases {
		ttl := NewTTLFromSeconds(tc.secs)
		got := c.StatusFor(ttl)
		if got != tc.want {
			t.Errorf("StatusFor(%ds): got %v, want %v", tc.secs, got, tc.want)
		}
	}
}

func TestDefaultExpiryWindow(t *testing.T) {
	w := DefaultExpiryWindow()
	if w.Critical != time.Hour {
		t.Errorf("expected Critical=1h, got %v", w.Critical)
	}
	if w.Warning != 6*time.Hour {
		t.Errorf("expected Warning=6h, got %v", w.Warning)
	}
}
