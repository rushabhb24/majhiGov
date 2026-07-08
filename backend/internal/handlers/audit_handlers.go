package handlers

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"strconv"

	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/models"
)

// LogAuditEvent records a security or admin event in the database audit log.
// Runs asynchronously to avoid blocking the main request cycle.
func LogAuditEvent(userID int, action, resourceType string, resourceID int, details string, r *http.Request) {
	ipAddress := "0.0.0.0"
	if r != nil {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ipAddress = ip
		} else {
			ipAddress = r.RemoteAddr
		}
	}

	go func() {
		var err error
		if userID > 0 {
			if resourceID > 0 {
				_, err = db.DB.Exec(`
					INSERT INTO audit_logs (user_id, action, resource_type, resource_id, details, ip_address, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
					userID, action, resourceType, resourceID, details, ipAddress,
				)
			} else {
				_, err = db.DB.Exec(`
					INSERT INTO audit_logs (user_id, action, resource_type, details, ip_address, created_at)
					VALUES ($1, $2, $3, $4, $5, NOW())`,
					userID, action, resourceType, details, ipAddress,
				)
			}
		} else {
			if resourceID > 0 {
				_, err = db.DB.Exec(`
					INSERT INTO audit_logs (action, resource_type, resource_id, details, ip_address, created_at)
					VALUES ($1, $2, $3, $4, $5, NOW())`,
					action, resourceType, resourceID, details, ipAddress,
				)
			} else {
				_, err = db.DB.Exec(`
					INSERT INTO audit_logs (action, resource_type, details, ip_address, created_at)
					VALUES ($1, $2, $3, $4, NOW())`,
					action, resourceType, details, ipAddress,
				)
			}
		}

		if err != nil {
			log.Printf("[AuditLog] Failed to insert audit log: %v", err)
		}
	}()
}

// GetAdminAuditLogsHandler fetches all audit logs with pagination and optional filters
func GetAdminAuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	actionFilter := r.URL.Query().Get("action")
	resourceTypeFilter := r.URL.Query().Get("resource_type")
	userIDStr := r.URL.Query().Get("user_id")

	page := 1
	limit := 20
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	offset := (page - 1) * limit

	query := `
		SELECT al.id, al.user_id, al.action, al.resource_type, al.resource_id, 
		       COALESCE(al.details::text, '{}'), al.ip_address, al.created_at,
		       COALESCE(u.email, 'System/Guest') as user_email
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE 1=1`

	var args []interface{}
	argIdx := 1

	if actionFilter != "" {
		query += " AND al.action = $" + strconv.Itoa(argIdx)
		args = append(args, actionFilter)
		argIdx++
	}
	if resourceTypeFilter != "" {
		query += " AND al.resource_type = $" + strconv.Itoa(argIdx)
		args = append(args, resourceTypeFilter)
		argIdx++
	}
	if userIDStr != "" {
		if uid, err := strconv.Atoi(userIDStr); err == nil {
			query += " AND al.user_id = $" + strconv.Itoa(argIdx)
			args = append(args, uid)
			argIdx++
		}
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM (" + query + ") as count_table"
	var total int
	err := db.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	query += " ORDER BY al.created_at DESC LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	type AuditLogRow struct {
		models.AuditLog
		UserEmail string `json:"user_email"`
	}

	var logs []AuditLogRow = []AuditLogRow{}
	for rows.Next() {
		var row AuditLogRow
		var userIDNull sql.NullInt64
		var resourceIDNull sql.NullInt64
		var detailsStr string

		err := rows.Scan(
			&row.ID, &userIDNull, &row.Action, &row.ResourceType, &resourceIDNull,
			&detailsStr, &row.IPAddress, &row.CreatedAt, &row.UserEmail,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Scan error: "+err.Error())
			return
		}

		if userIDNull.Valid {
			uid := int(userIDNull.Int64)
			row.UserID = &uid
		}
		if resourceIDNull.Valid {
			rid := int(resourceIDNull.Int64)
			row.ResourceID = &rid
		}
		row.Details = detailsStr
		logs = append(logs, row)
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data": logs,
		"meta": map[string]interface{}{
			"page":    page,
			"limit":   limit,
			"total":   total,
			"hasNext": (page * limit) < total,
		},
	})
}
