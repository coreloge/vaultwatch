package alert

import (
	"testing"
	"time"
)

func newTestBuilder() *Builder {
	return NewBuilder(30*time.Minute, 10*time.Minute)
}

func TestBuild_ContainsExpectedFields(t *testing.T) {
	b := newTestBuilder()
	expires := time.Now().Add(5 * time.Minute)
	p := b.Build("lease-123", "secret/db", expires, 300)

	if p.LeaseID != "lease-123" {
		t.Errorf("expected lease-123, got %s", p.LeaseID)
	}
	if p.Secret != "secret/db" {
		t.Errorf("expected secret/db, got %s", p.Secret)
	}
	if p.TTL != 300 {
		t.Errorf("expected TTL 300, got %d", p.TTL)
	}
	if p.Message == "" {
		t.Error("expected non-empty message")
	}
	if p.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestClassify_Critical(t *testing.T) {
	b := newTestBuilder()
	expires := time.Now().Add(5 * time.Minute) // under 10m critical threshold
	p := b.Build("lease-1", "secret/x", expires, 300)

	if p.Severity != SeverityCritical {
		t.Errorf("expected critical, got %s", p.Severity)
	}
}

func TestClassify_Warning(t *testing.T) {
	b := newTestBuilder()
	expires := time.Now().Add(20 * time.Minute) // between 10m and 30m
	p := b.Build("lease-2", "secret/y", expires, 1200)

	if p.Severity != SeverityWarning {
		t.Errorf("expected warning, got %s", p.Severity)
	}
}

func TestClassify_Info(t *testing.T) {
	b := newTestBuilder()
	expires := time.Now().Add(60 * time.Minute) // above 30m warning threshold
	p := b.Build("lease-3", "secret/z", expires, 3600)

	if p.Severity != SeverityInfo {
		t.Errorf("expected info, got %s", p.Severity)
	}
}

func TestNewBuilder_Thresholds(t *testing.T) {
	b := NewBuilder(1*time.Hour, 15*time.Minute)
	if b.warningThreshold != 1*time.Hour {
		t.Errorf("unexpected warning threshold: %v", b.warningThreshold)
	}
	if b.criticalThreshold != 15*time.Minute {
		t.Errorf("unexpected critical threshold: %v", b.criticalThreshold)
	}
}
