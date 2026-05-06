package clock_test

import (
	"testing"
	"time"

	"github.com/warden/vaultwatch/internal/clock"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestReal_Now_IsRecent(t *testing.T) {
	c := clock.New()
	before := time.Now()
	now := c.Now()
	after := time.Now()

	if now.Before(before) || now.After(after) {
		t.Errorf("Now() = %v, want between %v and %v", now, before, after)
	}
}

func TestReal_Since_IsNonNegative(t *testing.T) {
	c := clock.New()
	past := time.Now().Add(-time.Second)
	if c.Since(past) < 0 {
		t.Error("Since() returned negative duration for a past time")
	}
}

func TestReal_Until_IsPastForExpired(t *testing.T) {
	c := clock.New()
	past := time.Now().Add(-time.Second)
	if c.Until(past) >= 0 {
		t.Error("Until() should be negative for a time in the past")
	}
}

func TestMock_Now_ReturnsInitial(t *testing.T) {
	m := clock.NewMock(epoch)
	if got := m.Now(); !got.Equal(epoch) {
		t.Errorf("Now() = %v, want %v", got, epoch)
	}
}

func TestMock_Advance_MovesTime(t *testing.T) {
	m := clock.NewMock(epoch)
	m.Advance(5 * time.Minute)
	want := epoch.Add(5 * time.Minute)
	if got := m.Now(); !got.Equal(want) {
		t.Errorf("Now() after Advance = %v, want %v", got, want)
	}
}

func TestMock_Set_SetsAbsoluteTime(t *testing.T) {
	m := clock.NewMock(epoch)
	target := epoch.Add(24 * time.Hour)
	m.Set(target)
	if got := m.Now(); !got.Equal(target) {
		t.Errorf("Now() after Set = %v, want %v", got, target)
	}
}

func TestMock_Since_ReflectsMockTime(t *testing.T) {
	m := clock.NewMock(epoch)
	past := epoch.Add(-10 * time.Second)
	if got := m.Since(past); got != 10*time.Second {
		t.Errorf("Since() = %v, want %v", got, 10*time.Second)
	}
}

func TestMock_Until_ReflectsMockTime(t *testing.T) {
	m := clock.NewMock(epoch)
	future := epoch.Add(30 * time.Second)
	if got := m.Until(future); got != 30*time.Second {
		t.Errorf("Until() = %v, want %v", got, 30*time.Second)
	}
}

func TestMock_AdvanceMultipleTimes(t *testing.T) {
	m := clock.NewMock(epoch)
	m.Advance(time.Minute)
	m.Advance(time.Minute)
	want := epoch.Add(2 * time.Minute)
	if got := m.Now(); !got.Equal(want) {
		t.Errorf("Now() = %v, want %v", got, want)
	}
}
