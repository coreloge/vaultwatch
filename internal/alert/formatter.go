package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Format defines the serialisation format for alert payloads.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Formatter serialises a Payload into bytes for delivery.
type Formatter struct {
	format Format
}

// NewFormatter returns a Formatter for the given format string.
// Defaults to JSON if the format is unrecognised.
func NewFormatter(format string) *Formatter {
	f := Format(format)
	if f != FormatText {
		f = FormatJSON
	}
	return &Formatter{format: f}
}

// Encode serialises the payload according to the configured format.
func (f *Formatter) Encode(p Payload) ([]byte, error) {
	switch f.format {
	case FormatText:
		return f.encodeText(p), nil
	default:
		return f.encodeJSON(p)
	}
}

func (f *Formatter) encodeJSON(p Payload) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(p); err != nil {
		return nil, fmt.Errorf("alert: json encode: %w", err)
	}
	return buf.Bytes(), nil
}

func (f *Formatter) encodeText(p Payload) []byte {
	s := fmt.Sprintf(
		"[%s] %s | lease=%s secret=%s ttl=%ds expires=%s",
		p.Severity,
		p.Message,
		p.LeaseID,
		p.Secret,
		p.TTL,
		p.ExpiresAt.Format("2006-01-02T15:04:05Z"),
	)
	return []byte(s)
}
