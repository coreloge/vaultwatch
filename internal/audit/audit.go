// Package audit provides structured event logging for lease check and alert activity.
package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// EventType classifies the kind of audit event.
type EventType string

const (
	EventLeaseChecked EventType = "lease.checked"
	EventAlertSent    EventType = "alert.sent"
	EventAlertFailed  EventType = "alert.failed"
	EventLeaseRenewed EventType = "lease.renewed"
)

// Event represents a single auditable action.
type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	Type      EventType         `json:"type"`
	LeaseID   string            `json:"lease_id,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// New returns an audit Logger writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// Log encodes and writes a single audit event.
func (l *Logger) Log(eventType EventType, leaseID string, meta map[string]string) error {
	e := Event{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		LeaseID:   leaseID,
		Meta:      meta,
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = l.w.Write(data)
	return err
}
