// Package envelope wraps a lease alert payload with delivery metadata
// such as attempt count, origin, and a unique message ID.
package envelope

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/youorg/vaultwatch/internal/lease"
)

// Envelope wraps an alert payload with delivery metadata.
type Envelope struct {
	ID        string
	CreatedAt time.Time
	Attempts  int
	Origin    string
	Lease     lease.Info
}

// New creates a new Envelope for the given lease.Info and origin label.
func New(info lease.Info, origin string) *Envelope {
	return &Envelope{
		ID:        newID(),
		CreatedAt: time.Now().UTC(),
		Attempts:  0,
		Origin:    origin,
		Lease:     info,
	}
}

// Increment records an additional delivery attempt.
func (e *Envelope) Increment() {
	e.Attempts++
}

// Age returns how long ago the envelope was created.
func (e *Envelope) Age() time.Duration {
	return time.Since(e.CreatedAt)
}

// String returns a short human-readable summary.
func (e *Envelope) String() string {
	return fmt.Sprintf("envelope id=%s lease=%s attempts=%d origin=%s",
		e.ID, e.Lease.LeaseID, e.Attempts, e.Origin)
}

func newID() string {
	return uuid.NewString()
}
