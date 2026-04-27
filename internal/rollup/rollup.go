// Package rollup aggregates multiple lease alerts within a time window
// into a single batched notification, reducing webhook noise.
package rollup

import (
	"sync"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
)

// Config holds configuration for the rollup window.
type Config struct {
	// Window is the duration to collect alerts before flushing.
	Window time.Duration
	// MaxSize is the maximum number of events before an early flush.
	MaxSize int
}

// DefaultConfig returns sensible defaults for rollup.
func DefaultConfig() Config {
	return Config{
		Window:  30 * time.Second,
		MaxSize: 50,
	}
}

// Batch is a collection of lease infos accumulated within a window.
type Batch struct {
	Events    []lease.Info
	WindowEnd time.Time
}

// Rollup buffers lease events and flushes them as batches.
type Rollup struct {
	mu      sync.Mutex
	cfg     Config
	buf     []lease.Info
	flushCh chan Batch
	timer   *time.Timer
}

// New creates a new Rollup with the given config.
func New(cfg Config) *Rollup {
	if cfg.Window <= 0 {
		cfg.Window = DefaultConfig().Window
	}
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = DefaultConfig().MaxSize
	}
	return &Rollup{
		cfg:     cfg,
		flushCh: make(chan Batch, 16),
	}
}

// Add enqueues a lease event. If the buffer reaches MaxSize, it flushes immediately.
func (r *Rollup) Add(info lease.Info) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.buf) == 0 {
		r.timer = time.AfterFunc(r.cfg.Window, r.flushLocked)
	}
	r.buf = append(r.buf, info)

	if len(r.buf) >= r.cfg.MaxSize {
		r.timer.Stop()
		r.flush()
	}
}

// Batches returns a channel on which flushed batches are delivered.
func (r *Rollup) Batches() <-chan Batch {
	return r.flushCh
}

// flushLocked is called by the timer (no lock held by caller).
func (r *Rollup) flushLocked() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flush()
}

// flush drains the buffer into a Batch and sends it on flushCh.
// Caller must hold r.mu.
func (r *Rollup) flush() {
	if len(r.buf) == 0 {
		return
	}
	batch := Batch{
		Events:    make([]lease.Info, len(r.buf)),
		WindowEnd: time.Now(),
	}
	copy(batch.Events, r.buf)
	r.buf = r.buf[:0]
	select {
	case r.flushCh <- batch:
	default:
	}
}
