package middleware

import (
	"encoding/json"
	"net/http"

	"yojana-portal/backend/internal/db"
)

// AdminMiddleware checks if the authenticated user has is_admin = true in the database
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserIDFromContext(r.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized access: session details missing",
			})
			return
		}

		// Look up is_admin status in the database
		var isAdmin bool
		query := "SELECT is_admin FROM users WHERE id = $1"
		err = db.DB.QueryRow(query, userID).Scan(&isAdmin)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to verify administrator privileges: database error",
			})
			return
		}

		if !isAdmin {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Forbidden: administrative privileges required",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
