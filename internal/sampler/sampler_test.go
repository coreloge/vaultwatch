package sampler_test

import (
	"math/rand"
	"testing"

	"github.com/yourusername/vaultwatch/internal/lease"
	"github.com/yourusername/vaultwatch/internal/sampler"
)

func newInfo(status lease.Status) lease.Info {
	return lease.Info{
		LeaseID: "secret/data/test#abc",
		Status:  status,
	}
}

func deterministicSrc(seed int64) rand.Source {
	return rand.NewSource(seed)
}

func TestAllow_RateOne_AlwaysAllows(t *testing.T) {
	s := sampler.New(sampler.Config{Rate: 1.0}, deterministicSrc(1))
	for i := 0; i < 100; i++ {
		if !s.Allow(newInfo(lease.StatusWarning)) {
			t.Fatal("expected all events to be allowed at rate 1.0")
		}
	}
}

func TestAllow_RateZero_AlwaysDrops(t *testing.T) {
	s := sampler.New(sampler.Config{Rate: 0.0}, deterministicSrc(1))
	for i := 0; i < 100; i++ {
		if s.Allow(newInfo(lease.StatusWarning)) {
			t.Fatal("expected all events to be dropped at rate 0.0")
		}
	}
}

func TestAllow_StatusNotTargeted_AlwaysAllows(t *testing.T) {
	cfg := sampler.Config{
		Rate:     0.0,
		Statuses: []lease.Status{lease.StatusCritical},
	}
	s := sampler.New(cfg, deterministicSrc(1))
	// Warning is not in the targeted set, so it should always pass.
	for i := 0; i < 20; i++ {
		if !s.Allow(newInfo(lease.StatusWarning)) {
			t.Fatal("expected non-targeted status to always be allowed")
		}
	}
}

func TestAllow_StatusTargeted_AppliesSampling(t *testing.T) {
	cfg := sampler.Config{
		Rate:     0.0,
		Statuses: []lease.Status{lease.StatusCritical},
	}
	s := sampler.New(cfg, deterministicSrc(1))
	for i := 0; i < 20; i++ {
		if s.Allow(newInfo(lease.StatusCritical)) {
			t.Fatal("expected targeted status to be dropped at rate 0.0")
		}
	}
}

func TestAllow_RateClamped_Above(t *testing.T) {
	s := sampler.New(sampler.Config{Rate: 5.0}, deterministicSrc(1))
	if s.Rate() != 1.0 {
		t.Fatalf("expected rate clamped to 1.0, got %v", s.Rate())
	}
}

func TestAllow_RateClamped_Below(t *testing.T) {
	s := sampler.New(sampler.Config{Rate: -1.0}, deterministicSrc(1))
	if s.Rate() != 0.0 {
		t.Fatalf("expected rate clamped to 0.0, got %v", s.Rate())
	}
}

func TestDefaultConfig_RateIsOne(t *testing.T) {
	cfg := sampler.DefaultConfig()
	if cfg.Rate != 1.0 {
		t.Fatalf("expected default rate 1.0, got %v", cfg.Rate)
	}
}

func TestAllow_NilSource_DoesNotPanic(t *testing.T) {
	s := sampler.New(sampler.Config{Rate: 0.5}, nil)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	s.Allow(newInfo(lease.StatusWarning))
}
