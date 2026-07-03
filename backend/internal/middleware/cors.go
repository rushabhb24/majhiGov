package middleware

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// regexp to match localhost or 127.0.0.1 with any port (e.g. http://localhost:5173, http://localhost:5174)
var localhostRegex = regexp.MustCompile(`^https?://(localhost|127\.0\.0\.1)(:\d+)?$`)

// CorsMiddleware restricts origins to the configured allowlist (ALLOWED_ORIGINS env var).
// Falls back to localhost dev origins if env var is not set.
// Adds Access-Control-Allow-Credentials: true for httpOnly cookie support.
func CorsMiddleware(next http.Handler) http.Handler {
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		for _, o := range strings.Split(allowedOriginsStr, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}
	// Default dev origins when env var is not set
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{
			"http://localhost:5173",
			"http://localhost:4173",
			"http://localhost:3000",
			"http://127.0.0.1:5173",
		}
	}

	isAllowed := func(origin string) bool {
		// Automatically allow any localhost or 127.0.0.1 port for development ease
		if localhostRegex.MatchString(origin) {
			return true
		}
		for _, allowed := range allowedOrigins {
			if strings.EqualFold(allowed, origin) {
				return true
			}
		}
		return false
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Debug CORS issues
		log.Printf("CORS Check: origin=%q allowedOrigins=%v matchedLocalhost=%t", origin, allowedOrigins, localhostRegex.MatchString(origin))

		if origin != "" && isAllowed(origin) {
			// Allowed origin — reflect it back and allow credentials (required for httpOnly cookies)
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin != "" {
			// Origin not in allowlist — block state-changing methods
			if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"Origin not allowed"}`))
				return
			}
		}
		// No Origin header = same-origin request, always allowed

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
