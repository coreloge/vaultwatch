// Package observe provides a lightweight event observer that fans out
// lease lifecycle events to multiple registered listeners.
package observe

import (
	"sync"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Handler is a function that receives a lease.Info event.
type Handler func(info lease.Info)

// Observer multiplexes lease events to a set of registered handlers.
type Observer struct {
	mu       sync.RWMutex
	handlers []Handler
}

// New returns an initialised Observer with no handlers.
func New() *Observer {
	return &Observer{}
}

// Register adds h to the set of handlers that will be called on Emit.
// Handlers are called in registration order.
func (o *Observer) Register(h Handler) {
	if h == nil {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	o.handlers = append(o.handlers, h)
}

// Emit delivers info to every registered handler.
// Each handler is invoked synchronously; callers that need isolation
// should wrap their handler in a goroutine themselves.
func (o *Observer) Emit(info lease.Info) {
	o.mu.RLock()
	handlers := make([]Handler, len(o.handlers))
	copy(handlers, o.handlers)
	o.mu.RUnlock()

	for _, h := range handlers {
		h(info)
	}
}

// Len returns the number of registered handlers.
func (o *Observer) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.handlers)
}

// Reset removes all registered handlers.
func (o *Observer) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.handlers = nil
}
