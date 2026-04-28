// Package sink implements a fan-out delivery layer for VaultWatch alert
// payloads. A Sink holds a list of named Target implementations and
// delivers each alert to all of them, collecting per-target errors
// without short-circuiting so that a failing target does not prevent
// delivery to healthy ones.
//
// Typical usage:
//
//	s := sink.New(webhookTarget, logTarget)
//	if err := s.SendAll(ctx, payload); err != nil {
//		log.Println("one or more targets failed:", err)
//	}
package sink
