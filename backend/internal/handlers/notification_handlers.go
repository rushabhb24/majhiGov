package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/middleware"
)

// GetUserNotificationsHandler fetches all notifications for the logged-in user
func GetUserNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	type NotificationRow struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		Message   string    `json:"message"`
		Type      string    `json:"type"`
		IsRead    bool      `json:"is_read"`
		CreatedAt time.Time `json:"created_at"`
	}

	rows, err := db.DB.Query(
		`SELECT id, title, message, type, is_read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`,
		userID,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}
	defer rows.Close()

	var notifications []NotificationRow
	for rows.Next() {
		var n NotificationRow
		if err := rows.Scan(&n.ID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err == nil {
			notifications = append(notifications, n)
		}
	}
	if notifications == nil {
		notifications = []NotificationRow{}
	}
	writeJSONResponse(w, http.StatusOK, notifications)
}

// MarkNotificationsReadHandler marks all notifications as read for the logged-in user
func MarkNotificationsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	_, err = db.DB.Exec(`UPDATE notifications SET is_read = true WHERE user_id = $1`, userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to mark notifications as read")
		return
	}
	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"success": true})
}

// AdminSendNotificationHandler handles POST requests to send notifications to users or broadcast
func AdminSendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		UserID  int    `json:"user_id"` // 0 = broadcast to all non-admin users
		Title   string `json:"title"`
		Message string `json:"message"`
		Type    string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Title == "" || req.Message == "" {
		writeJSONError(w, http.StatusBadRequest, "title and message are required")
		return
	}
	if req.Type == "" {
		req.Type = "info"
	}

	if req.UserID == 0 {
		// Broadcast: insert notification record for all non-admin users
		_, err := db.DB.Exec(
			`INSERT INTO notifications (user_id, title, message, type)
			 SELECT id, $1, $2, $3 FROM users WHERE is_admin = false AND is_active = true`,
			req.Title, req.Message, req.Type,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to broadcast notification")
			return
		}
	} else {
		// Single user notification
		_, err := db.DB.Exec(
			`INSERT INTO notifications (user_id, title, message, type) VALUES ($1, $2, $3, $4)`,
			req.UserID, req.Title, req.Message, req.Type,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to send notification")
			return
		}
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Notification sent successfully",
	})
}
