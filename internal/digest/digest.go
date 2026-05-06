// Package digest provides content-based fingerprinting for lease alert payloads.
// A digest is a short, stable hash derived from the fields that define
// the identity of an alert, allowing downstream components to detect
// duplicate or equivalent notifications without comparing full payloads.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/your-org/vaultwatch/internal/lease"
)

// Digester computes fingerprints for lease events.
type Digester struct {
	// truncate controls how many hex characters to retain (0 = full 64).
	truncate int
}

// New returns a Digester. truncate specifies the desired output length
// (number of hex characters). Pass 0 to keep the full SHA-256 hex string.
func New(truncate int) *Digester {
	if truncate < 0 {
		truncate = 0
	}
	return &Digester{truncate: truncate}
}

// Compute returns a deterministic fingerprint for the given lease info.
// The digest is derived from the lease ID, its status, and the expiry
// minute (truncated to the minute boundary) so that repeated checks
// within the same minute produce the same digest.
func (d *Digester) Compute(info lease.Info) string {
	minute := info.ExpiresAt.Truncate(time.Minute).Unix()
	raw := fmt.Sprintf("%s|%s|%d", info.LeaseID, info.Status, minute)
	sum := sha256.Sum256([]byte(raw))
	hex := hex.EncodeToString(sum[:])
	if d.truncate > 0 && d.truncate < len(hex) {
		return hex[:d.truncate]
	}
	return hex
}

// Equal reports whether two digests are identical.
func Equal(a, b string) bool {
	return a == b
}
