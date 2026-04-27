// Package window implements a sliding time-window counter.
//
// A Counter records timestamped events and exposes a Count method
// that returns how many events fall within a configurable trailing
// duration.  Old events are lazily evicted on each Add or Count call,
// keeping memory usage proportional to the event rate rather than
// total uptime.
//
// Typical use:
//
//	c := window.New(time.Minute)
//	c.Add()          // record an event now
//	n := c.Count()   // events in the last minute
package window
