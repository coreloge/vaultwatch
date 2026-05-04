// Package tag provides lease tagging and label enrichment.
//
// A Tagger maps lease identifiers to a set of string key/value labels that can
// be attached to outgoing alert payloads. Labels are composed from two sources:
//
//  1. Static tags – applied unconditionally to every lease.
//  2. Prefix tags – applied only when the lease ID begins with a registered
//     prefix, allowing environment- or path-scoped metadata.
//
// Prefix tags override static tags when keys collide, giving operators fine-
// grained control over per-path annotations without duplicating configuration.
package tag
