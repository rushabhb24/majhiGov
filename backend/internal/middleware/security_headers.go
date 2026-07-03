package middleware

import "net/http"

// SecurityHeadersMiddleware adds common security‑related HTTP headers to each response.
// It is lightweight and can be extended with additional headers as needed.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Prevent MIME sniffing
        w.Header().Set("X-Content-Type-Options", "nosniff")
        // Disallow the page from being displayed in an iframe
        w.Header().Set("X-Frame-Options", "DENY")
        // Enable basic XSS protection in older browsers
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        // Content Security Policy – restrict resources to same origin
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        // Continue to the next handler
        next.ServeHTTP(w, r)
    })
}
