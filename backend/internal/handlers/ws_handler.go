package handlers

import (
	"net/http"

	"yojana-portal/backend/internal/middleware"
	"yojana-portal/backend/internal/ws"
)

// WebSocketHandler upgrades the HTTP connection to a WebSocket for real-time notifications.
// Requires AuthMiddleware to have run first (sets user_id in context via cookie-validated JWT).
func WebSocketHandler(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "Authentication required for WebSocket connection")
			return
		}
		hub.ServeWS(userID, w, r)
	}
}
