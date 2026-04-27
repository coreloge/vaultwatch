package envelope_test

import (
	"sync"
	"testing"

	"github.com/youorg/vaultwatch/internal/envelope"
)

func TestQueue_PushAndPop(t *testing.T) {
	q := envelope.NewQueue()
	e := envelope.New(sampleInfo(), "test")
	q.Push(e)
	got := q.Pop()
	if got == nil {
		t.Fatal("expected non-nil envelope")
	}
	if got.ID != e.ID {
		t.Fatalf("expected ID %q, got %q", e.ID, got.ID)
	}
}

func TestQueue_PopEmpty(t *testing.T) {
	q := envelope.NewQueue()
	if got := q.Pop(); got != nil {
		t.Fatalf("expected nil from empty queue, got %v", got)
	}
}

func TestQueue_Len(t *testing.T) {
	q := envelope.NewQueue()
	if q.Len() != 0 {
		t.Fatal("expected empty queue")
	}
	q.Push(envelope.New(sampleInfo(), "a"))
	q.Push(envelope.New(sampleInfo(), "b"))
	if q.Len() != 2 {
		t.Fatalf("expected len 2, got %d", q.Len())
	}
}

func TestQueue_Drain(t *testing.T) {
	q := envelope.NewQueue()
	for i := 0; i < 3; i++ {
		q.Push(envelope.New(sampleInfo(), "origin"))
	}
	out := q.Drain()
	if len(out) != 3 {
		t.Fatalf("expected 3 envelopes, got %d", len(out))
	}
	if q.Len() != 0 {
		t.Fatal("expected empty queue after drain")
	}
}

func TestQueue_FIFO_Order(t *testing.T) {
	q := envelope.NewQueue()
	a := envelope.New(sampleInfo(), "first")
	b := envelope.New(sampleInfo(), "second")
	q.Push(a)
	q.Push(b)
	if q.Pop().ID != a.ID {
		t.Fatal("expected first-pushed envelope to be popped first")
	}
}

func TestQueue_ConcurrentPush(t *testing.T) {
	q := envelope.NewQueue()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q.Push(envelope.New(sampleInfo(), "concurrent"))
		}()
	}
	wg.Wait()
	if q.Len() != 50 {
		t.Fatalf("expected 50 envelopes, got %d", q.Len())
	}
}
