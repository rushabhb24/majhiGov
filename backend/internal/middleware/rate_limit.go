package middleware

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type ipRequestTracker struct {
	timestamps []time.Time
}

type RateLimiter struct {
	mu             sync.Mutex
	trackers        map[string]*ipRequestTracker
	maxRequests    int
	windowPeriod   time.Duration
	trustedProxies []string
}

func NewRateLimiter(maxRequests int, windowPeriod time.Duration) *RateLimiter {
	// Load trusted proxy CIDRs from env: TRUSTED_PROXIES=127.0.0.1,10.0.0.0/8
	var trustedProxies []string
	if tp := os.Getenv("TRUSTED_PROXIES"); tp != "" {
		for _, p := range strings.Split(tp, ",") {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				trustedProxies = append(trustedProxies, trimmed)
			}
		}
	}
	// Always trust localhost
	trustedProxies = append(trustedProxies, "127.0.0.1", "::1")

	return &RateLimiter{
		trackers:       make(map[string]*ipRequestTracker),
		maxRequests:   maxRequests,
		windowPeriod:  windowPeriod,
		trustedProxies: trustedProxies,
	}
}

// isTrustedProxy checks if a given IP is in the trusted proxy list
func (rl *RateLimiter) isTrustedProxy(ip string) bool {
	parsedIP := net.ParseIP(ip)
	for _, trustedStr := range rl.trustedProxies {
		if trustedStr == ip {
			return true
		}
		_, cidr, err := net.ParseCIDR(trustedStr)
		if err == nil && parsedIP != nil && cidr.Contains(parsedIP) {
			return true
		}
	}
	return false
}

// extractClientIP determines the real client IP.
// Only trusts X-Forwarded-For / X-Real-IP if the direct connection is from a trusted proxy.
// This prevents IP spoofing attacks.
func (rl *RateLimiter) extractClientIP(r *http.Request) string {
	directIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		directIP = r.RemoteAddr
	}

	if rl.isTrustedProxy(directIP) {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			// X-Forwarded-For: client, proxy1, proxy2 — leftmost is the original client
			parts := strings.Split(xff, ",")
			clientIP := strings.TrimSpace(parts[0])
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			trimmed := strings.TrimSpace(xri)
			if net.ParseIP(trimmed) != nil {
				return trimmed
			}
		}
	}
	// Not a trusted proxy or no forwarded header — use direct connection IP
	return directIP
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := rl.extractClientIP(r)

		rl.mu.Lock()
		tracker, exists := rl.trackers[ip]
		if !exists {
			tracker = &ipRequestTracker{timestamps: []time.Time{}}
			rl.trackers[ip] = tracker
		}

		now := time.Now()
		validTimestamps := []time.Time{}
		for _, t := range tracker.timestamps {
			if now.Sub(t) <= rl.windowPeriod {
				validTimestamps = append(validTimestamps, t)
			}
		}

		if len(validTimestamps) >= rl.maxRequests {
			rl.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
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
