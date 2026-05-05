// Package ratelimit provides a small, in-memory fixed-window rate limiter
// keyed on an arbitrary string (typically client IP). It's intended for
// throttling abuse on low-traffic, single-process endpoints — the contact
// form being the only consumer today.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter caps the number of allowed events per key inside a fixed
// time window. The window resets the next time Allow is called after
// the previous window expires.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	limit   int
	window  time.Duration
}

type bucket struct {
	count   int
	resetAt time.Time
}

// New returns a Limiter that allows at most limit calls per window per key.
func New(limit int, window time.Duration) *Limiter {
	return &Limiter{
		buckets: make(map[string]*bucket),
		limit:   limit,
		window:  window,
	}
}

// Allow records an attempt for key and returns true if the key is still
// under its quota for the current window; false once the limit is hit.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[key]
	if !ok || now.After(b.resetAt) {
		l.buckets[key] = &bucket{count: 1, resetAt: now.Add(l.window)}
		return true
	}
	if b.count >= l.limit {
		return false
	}
	b.count++
	return true
}
