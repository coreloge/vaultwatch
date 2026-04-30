package fanout_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/youorg/vaultwatch/internal/fanout"
	"github.com/youorg/vaultwatch/internal/lease"
)

// slowHandler simulates a handler that takes a small amount of work.
type slowHandler struct {
	total atomic.Int64
}

func (h *slowHandler) Handle(_ context.Context, _ lease.Info) error {
	h.total.Add(1)
	return nil
}

// TestConcurrentSend_AllEventsDelivered fires many goroutines each calling
// Send and verifies that every handler receives every event exactly once.
func TestConcurrentSend_AllEventsDelivered(t *testing.T) {
	const goroutines = 10

	h := &slowHandler{}
	f := fanout.New(h)

	info := lease.Info{LeaseID: "secret/data/concurrent", Status: lease.StatusCritical}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < goroutines; i++ {
			go func() {
				f.Send(context.Background(), info)
			}()
		}
	}()
	<-done
}

// TestSend_AllHandlersCalledEvenIfOneFails ensures partial failures do not
// short-circuit remaining handlers.
func TestSend_AllHandlersCalledEvenIfOneFails(t *testing.T) {
	import_errors := func() error { return context.DeadlineExceeded }()

	var called [3]atomic.Int64
	handlers := []fanout.Handler{
		&countHandler{},
		&countHandler{err: import_errors},
		&countHandler{},
	}
	f := fanout.New(handlers...)
	_ = called

	errs := f.Send(context.Background(), sampleInfo())
	if len(errs) != 1 {
		t.Errorf("expected exactly 1 error, got %d", len(errs))
	}
	for i, h := range handlers {
		if h.(*countHandler).calls.Load() != 1 {
			t.Errorf("handler %d was not called exactly once", i)
		}
	}
}
