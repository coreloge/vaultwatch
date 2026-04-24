package lease

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"
)

func newTestRenewer(store *Store, renewFn RenewFunc, threshold time.Duration) *Renewer {
	logger := log.New(os.Stdout, "", 0)
	return NewRenewer(store, renewFn, threshold, logger)
}

func TestRenewEligible_RenewsExpiring(t *testing.T) {
	store := NewStore()
	expiring := Info{
		LeaseID:    "secret/expiring",
		ExpireTime: time.Now().Add(30 * time.Second),
		TTL:        30 * time.Second,
	}
	store.Set(expiring)

	renewFn := func(_ context.Context, id string) (time.Duration, error) {
		return 3600 * time.Second, nil
	}

	r := newTestRenewer(store, renewFn, 5*time.Minute)
	count := r.RenewEligible(context.Background())

	if count != 1 {
		t.Fatalf("expected 1 renewal, got %d", count)
	}

	updated, ok := store.Get(expiring.LeaseID)
	if !ok {
		t.Fatal("expected lease to still exist in store")
	}
	if time.Until(updated.ExpireTime) < time.Hour-5*time.Second {
		t.Errorf("expected renewed TTL near 1h, got %s", time.Until(updated.ExpireTime))
	}
}

func TestRenewEligible_SkipsHealthyLeases(t *testing.T) {
	store := NewStore()
	healthy := Info{
		LeaseID:    "secret/healthy",
		ExpireTime: time.Now().Add(2 * time.Hour),
		TTL:        2 * time.Hour,
	}
	store.Set(healthy)

	called := false
	renewFn := func(_ context.Context, _ string) (time.Duration, error) {
		called = true
		return 0, nil
	}

	r := newTestRenewer(store, renewFn, 5*time.Minute)
	count := r.RenewEligible(context.Background())

	if count != 0 {
		t.Errorf("expected 0 renewals, got %d", count)
	}
	if called {
		t.Error("renewFn should not have been called for a healthy lease")
	}
}

func TestRenewEligible_HandlesRenewError(t *testing.T) {
	store := NewStore()
	expiring := Info{
		LeaseID:    "secret/broken",
		ExpireTime: time.Now().Add(10 * time.Second),
		TTL:        10 * time.Second,
	}
	store.Set(expiring)

	renewFn := func(_ context.Context, _ string) (time.Duration, error) {
		return 0, errors.New("vault unavailable")
	}

	r := newTestRenewer(store, renewFn, 5*time.Minute)
	count := r.RenewEligible(context.Background())

	if count != 0 {
		t.Errorf("expected 0 successful renewals, got %d", count)
	}
}
