// Package notify provides the Dispatcher type which coordinates building,
// formatting, and delivering lease expiration alerts to configured webhook
// endpoints.
//
// A Dispatcher is constructed with a webhook target URL, an optional HMAC
// signing secret, and an output format ("json" or "text"). It composes the
// alert, webhook, and lease packages to produce a single entry point for
// sending notifications from the monitor loop.
//
// Example usage:
//
//	d, err := notify.New(notify.Config{
//		WebhookURL: "https://hooks.example.com/vault",
//		Format:     "json",
//	}, logger)
//	if err != nil { ... }
//	if err := d.Dispatch(ctx, leaseInfo); err != nil { ... }
package notify
