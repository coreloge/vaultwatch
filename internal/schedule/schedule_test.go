package schedule_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/schedule"
)

func TestNew_DefaultsInterval(t *testing.T) {
	s := schedule.New(0)
	if s.Interval() != 30*time.Second {
		t.Fatalf("expected 30s default, got %s", s.Interval())
	}
}

func TestNew_CustomInterval(t *testing.T) {
	s := schedule.New(5 * time.Second)
	if s.Interval() != 5*time.Second {
		t.Fatalf("expected 5s, got %s", s.Interval())
	}
}

func TestRun_CallsImmediately(t *testing.T) {
	var calls int64
	s := schedule.New(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.Run(ctx, func(_ context.Context) {
			atomic.AddInt64(&calls, 1)
		})
		close(done)
	}()

	<-done
	if atomic.LoadInt64(&calls) < 1 {
		t.Fatal("expected at least one immediate call")
	}
}

func TestRun_TicksAtInterval(t *testing.T) {
	var calls int64
	s := schedule.New(50 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.Run(ctx, func(_ context.Context) {
			atomic.AddInt64(&calls, 1)
		})
		close(done)
	}()

	<-done
	// immediate call + ~3 ticks within 180 ms
	if c := atomic.LoadInt64(&calls); c < 3 {
		t.Fatalf("expected >=3 calls, got %d", c)
	}
}

func TestRun_StopHaltsLoop(t *testing.T) {
	var calls int64
	s := schedule.New(20 * time.Millisecond)
	ctx := context.Background()

	done := make(chan struct{})
	go func() {
		s.Run(ctx, func(_ context.Context) {
			atomic.AddInt64(&calls, 1)
		})
		close(done)
	}()

	time.Sleep(55 * time.Millisecond)
	s.Stop()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("scheduler did not stop within timeout")
	}

	snap := atomic.LoadInt64(&calls)
	time.Sleep(60 * time.Millisecond)
	if atomic.LoadInt64(&calls) != snap {
		t.Fatal("scheduler continued ticking after Stop")
	}
}

func TestRun_ContextPassedToCallback(t *testing.T) {
	var received context.Context
	s := schedule.New(10 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		s.Run(ctx, func(c context.Context) {
			received = c
			cancel()
		})
		close(done)
	}()

	<-done
	if received == nil {
		t.Fatal("context was not passed to callback")
	}
}
