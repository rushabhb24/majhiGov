package handlers

import (
    "net/http"
)

// HealthHandler provides a simple health check endpoint.
// It responds with HTTP 200 and a small JSON payload indicating the service is alive.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"status":"ok"}`))
}
