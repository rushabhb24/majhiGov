package middleware

import (
	"net/http"
	"sync"
	"time"
)

type ipRequestTracker struct {
	timestamps []time.Time
}

type RateLimiter struct {
	mu           sync.Mutex
	trackers     map[string]*ipRequestTracker
	maxRequests  int
	windowPeriod time.Duration
}

func NewRateLimiter(maxRequests int, windowPeriod time.Duration) *RateLimiter {
	return &RateLimiter{
		trackers:     make(map[string]*ipRequestTracker),
		maxRequests:  maxRequests,
		windowPeriod: windowPeriod,
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Strip port if present in RemoteAddr (e.g. "127.0.0.1:56889")
		if idx := len(ip) - 1; idx >= 0 {
			for idx >= 0 && ip[idx] != ':' {
				idx--
			}
			if idx >= 0 {
				ip = ip[:idx]
			}
		}

		// Handle forward headers if hosted behind a proxy (like Nginx/Cloudflare)
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		} else if xri := r.Header.Get("X-Real-IP"); xri != "" {
			ip = xri
		}

		rl.mu.Lock()
		tracker, exists := rl.trackers[ip]
		if !exists {
			tracker = &ipRequestTracker{timestamps: []time.Time{}}
			rl.trackers[ip] = tracker
		}

		now := time.Now()
		// Filter out timestamps outside the window period
		validTimestamps := []time.Time{}
		for _, t := range tracker.timestamps {
			if now.Sub(t) <= rl.windowPeriod {
				validTimestamps = append(validTimestamps, t)
			}
		}

		if len(validTimestamps) >= rl.maxRequests {
			rl.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"Too many requests. Rate limit exceeded. Please try again later."}`))
			return
		}

		validTimestamps = append(validTimestamps, now)
		tracker.timestamps = validTimestamps
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
