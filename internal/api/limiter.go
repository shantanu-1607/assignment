package api

import (
	"sync"
	"time"
)

type bucket struct {
	count   int
	expires time.Time
}

// Fixed-window, per-process demonstration. Keys use socket peers, never X-Forwarded-For.
// Full key tables reject new peers rather than evicting active abuse counters.
type limiter struct {
	mu             sync.Mutex
	entries        map[string]bucket
	limit, maxKeys int
	now            func() time.Time
}

func (l *limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	b, ok := l.entries[key]
	if !ok || !now.Before(b.expires) {
		if !ok && len(l.entries) >= l.maxKeys {
			for k, v := range l.entries {
				if !now.Before(v.expires) {
					delete(l.entries, k)
				}
			}
			if len(l.entries) >= l.maxKeys {
				return false
			}
		}
		b = bucket{expires: now.Add(time.Minute)}
	}
	if b.count >= l.limit {
		return false
	}
	b.count++
	l.entries[key] = b
	return true
}
