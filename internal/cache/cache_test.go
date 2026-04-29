package cache_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/cache"
)

func newCache(ttl time.Duration) *cache.Cache[string, string] {
	return cache.New[string, string](ttl)
}

func TestSet_AndGet_ReturnsValue(t *testing.T) {
	c := newCache(time.Minute)
	c.Set("key", "value")
	got, ok := c.Get("key")
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if got != "value" {
		t.Fatalf("expected 'value', got %q", got)
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	c := newCache(time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestGet_ExpiredEntry_ReturnsFalse(t *testing.T) {
	c := newCache(time.Millisecond)
	c.Set("key", "value")
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected expired entry to be absent")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c := newCache(time.Minute)
	c.Set("key", "value")
	c.Delete("key")
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	c := newCache(time.Millisecond)
	c.Set("a", "1")
	c.Set("b", "2")
	time.Sleep(5 * time.Millisecond)
	c.Set("c", "3") // fresh entry with a new cache instance TTL
	removed := c.Purge()
	if removed != 2 {
		t.Fatalf("expected 2 removed, got %d", removed)
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1 remaining, got %d", c.Len())
	}
}

func TestLen_ReflectsTotalEntries(t *testing.T) {
	c := newCache(time.Minute)
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
	c.Set("x", "1")
	c.Set("y", "2")
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestNew_ZeroTTL_DefaultsToOneMinute(t *testing.T) {
	c := cache.New[string, int](0)
	c.Set("k", 42)
	v, ok := c.Get("k")
	if !ok || v != 42 {
		t.Fatal("expected value to be retrievable with default TTL")
	}
}
