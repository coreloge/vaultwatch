// Package batch provides a collector that accumulates lease events and
// flushes them as a grouped payload to a downstream handler.
package batch

import (
	"context"
	"sync"
	"time"

	"github.com/yourusername/vaultwatch/internal/lease"
)

// Handler is called with a slice of lease infos when the batch is flushed.
type Handler func(ctx context.Context, items []lease.Info)

// Config controls flush behaviour.
type Config struct {
	MaxSize  int
	Window   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxSize: 50,
		Window:  10 * time.Second,
	}
}

// Collector accumulates lease.Info values and flushes them periodically or
// when the batch reaches MaxSize.
type Collector struct {
	cfg     Config
	handler Handler
	mu      sync.Mutex
	items   []lease.Info
	timer   *time.Timer
}

// New creates a Collector with the given config and flush handler.
func New(cfg Config, handler Handler) *Collector {
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = DefaultConfig().MaxSize
	}
	if cfg.Window <= 0 {
		cfg.Window = DefaultConfig().Window
	}
	return &Collector{
		cfg:     cfg,
		handler: handler,
	}
}

// Add appends a lease.Info to the batch. It flushes immediately if MaxSize is
// reached, or arms a timer for the next Window flush.
func (c *Collector) Add(ctx context.Context, info lease.Info) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = append(c.items, info)

	if len(c.items) >= c.cfg.MaxSize {
		if c.timer != nil {
			c.timer.Stop()
			c.timer = nil
		}
		c.flush(ctx)
		return
	}

	if c.timer == nil {
		c.timer = time.AfterFunc(c.cfg.Window, func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.flush(ctx)
			c.timer = nil
		})
	}
}

// Flush forces an immediate flush regardless of batch size or timer state.
func (c *Collector) Flush(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.flush(ctx)
}

// flush sends current items to the handler and resets the slice. Caller must
// hold c.mu.
func (c *Collector) flush(ctx context.Context) {
	if len(c.items) == 0 {
		return
	}
	snapshot := make([]lease.Info, len(c.items))
	copy(snapshot, c.items)
	c.items = c.items[:0]
	go c.handler(ctx, snapshot)
}
