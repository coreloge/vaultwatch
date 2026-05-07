package observe_test

import (
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
	"github.com/your-org/vaultwatch/internal/observe"
)

func sampleInfo(id string) lease.Info {
	return lease.Info{
		LeaseID:   id,
		MountPath: "secret/",
		TTL:       lease.NewTTLFromSeconds(300),
	}
}

func TestNew_ReturnsEmptyObserver(t *testing.T) {
	o := observe.New()
	if o.Len() != 0 {
		t.Fatalf("expected 0 handlers, got %d", o.Len())
	}
}

func TestRegister_NilHandlerIgnored(t *testing.T) {
	o := observe.New()
	o.Register(nil)
	if o.Len() != 0 {
		t.Fatalf("expected 0 handlers after nil register, got %d", o.Len())
	}
}

func TestEmit_CallsAllHandlers(t *testing.T) {
	o := observe.New()

	var mu sync.Mutex
	var received []string

	o.Register(func(info lease.Info) {
		mu.Lock()
		received = append(received, "h1:"+info.LeaseID)
		mu.Unlock()
	})
	o.Register(func(info lease.Info) {
		mu.Lock()
		received = append(received, "h2:"+info.LeaseID)
		mu.Unlock()
	})

	o.Emit(sampleInfo("lease-abc"))

	if len(received) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(received))
	}
	if received[0] != "h1:lease-abc" {
		t.Errorf("unexpected first call: %s", received[0])
	}
	if received[1] != "h2:lease-abc" {
		t.Errorf("unexpected second call: %s", received[1])
	}
}

func TestEmit_NoHandlers_IsNoop(t *testing.T) {
	o := observe.New()
	// Should not panic.
	o.Emit(sampleInfo("lease-xyz"))
}

func TestReset_ClearsHandlers(t *testing.T) {
	o := observe.New()
	o.Register(func(lease.Info) {})
	o.Register(func(lease.Info) {})
	o.Reset()
	if o.Len() != 0 {
		t.Fatalf("expected 0 handlers after reset, got %d", o.Len())
	}
}

func TestEmit_ConcurrentSafe(t *testing.T) {
	o := observe.New()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		o.Register(func(lease.Info) { time.Sleep(time.Microsecond) })
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			o.Emit(sampleInfo("concurrent"))
		}(i)
	}
	wg.Wait()
}
