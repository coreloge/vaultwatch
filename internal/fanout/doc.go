// Package fanout provides concurrent broadcast delivery of lease events to
// multiple downstream handlers.
//
// A Fanout is constructed with one or more Handler implementations and
// dispatches each incoming lease.Info to all of them in parallel, collecting
// any errors for the caller to inspect or log.
//
// Usage:
//
//	f := fanout.New(handlerA, handlerB, handlerC)
//	errs := f.Send(ctx, info)
package fanout
