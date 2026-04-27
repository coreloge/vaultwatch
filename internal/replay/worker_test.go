package replay_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
	"github.com/your-org/vaultwatch/internal/replay"
)

type mockDispatcher struct {
	calls  atomic.Int32
	failOn int32 // fail if calls <= failOn
}

func (m *mockDispatcher) Dispatch(_ context.Context, _ lease.Info) error {
	n := m.calls.Add(1)
	if n <= m.failOn {
		return errors.New("dispatch error")
	}
	return nil
}

func TestWorker_ReplaySucceeds(t *testing.T) {
	s := replay.New(time.Hour)
	s.Add(sampleInfo("lease-1"))

	d := &mockDispatcher{}
	w := replay.NewWorker(s, d, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go w.Run(ctx)
	<-ctx.Done()

	if d.calls.Load() == 0 {
		t.Error("expected at least one dispatch call")
	}
	if s.Len() != 0 {
		t.Errorf("expected store to be empty after successful replay, got %d", s.Len())
	}
}

func TestWorker_RequeuesOnFailure(t *testing.T) {
	s := replay.New(time.Hour)
	s.Add(sampleInfo("lease-fail"))

	// always fail
	d := &mockDispatcher{failOn: 99}
	w := replay.NewWorker(s, d, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	go w.Run(ctx)
	<-ctx.Done()

	if d.calls.Load() == 0 {
		t.Error("expected dispatch to have been attempted")
	}
	// entry should have been re-added each time
	if s.Len() == 0 {
		t.Error("expected entry to be re-queued after failure")
	}
}

func TestWorker_DefaultInterval(t *testing.T) {
	s := replay.New(time.Hour)
	d := &mockDispatcher{}
	// zero interval should not panic
	w := replay.NewWorker(s, d, 0)
	if w == nil {
		t.Fatal("expected non-nil worker")
	}
}
