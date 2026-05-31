package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/models"
)

// AdminJobsHandler handles GET (list all) and POST (create job ad) for admins
func AdminJobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getAdminAllJobs(w, r)
	case "POST":
		createAdminJob(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// AdminJobDetailsHandler handles PUT (edit) and DELETE (remove) for a specific job ID
func AdminJobDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL path (e.g. /api/admin/jobs/:id)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeJSONError(w, http.StatusBadRequest, "Invalid request path")
		return
	}
	jobIDStr := parts[3]
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	switch r.Method {
	case "PUT":
		updateAdminJob(w, r, jobID)
	case "DELETE":
		deleteAdminJob(w, r, jobID)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func getAdminAllJobs(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, title, title_hi, title_mr, organization, organization_hi, organization_mr,
		       description, description_hi, description_mr, education_qualification, documents_required,
		       min_age, max_age, last_date, apply_link, general_fee, obc_fee, sc_st_fee, women_fee,
		       is_active, clicks_count, created_at, updated_at
		FROM government_jobs
		ORDER BY created_at DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to query database: "+err.Error())
		return
	}
	defer rows.Close()

	var jobs []models.GovernmentJob = []models.GovernmentJob{}
	for rows.Next() {
		var j models.GovernmentJob
		var lastDate time.Time
		err := rows.Scan(
			&j.ID, &j.Title, &j.TitleHi, &j.TitleMr, &j.Organization, &j.OrganizationHi, &j.OrganizationMr,
			&j.Description, &j.DescriptionHi, &j.DescriptionMr, &j.EducationQualification, (*pq.StringArray)(&j.DocumentsRequired),
			&j.MinAge, &j.MaxAge, &lastDate, &j.ApplyLink, &j.GeneralFee, &j.ObcFee, &j.ScStFee, &j.WomenFee,
			&j.IsActive, &j.ClicksCount, &j.CreatedAt, &j.UpdatedAt,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to scan job row: "+err.Error())
			return
		}
		j.LastDate = lastDate.Format("2006-01-02")
		jobs = append(jobs, j)
	}

	writeJSONResponse(w, http.StatusOK, jobs)
}

type JobPayload struct {
	Title                  string   `json:"title"`
	TitleHi                string   `json:"title_hi"`
	TitleMr                string   `json:"title_mr"`
	Organization           string   `json:"organization"`
	OrganizationHi         string   `json:"organization_hi"`
	OrganizationMr         string   `json:"organization_mr"`
	Description            string   `json:"description"`
	DescriptionHi          string   `json:"description_hi"`
	DescriptionMr          string   `json:"description_mr"`
	EducationQualification string   `json:"education_qualification"`
	DocumentsRequired      []string `json:"documents_required"`
	MinAge                 int      `json:"min_age"`
	MaxAge                 int      `json:"max_age"`
	LastDate               string   `json:"last_date"` // YYYY-MM-DD
	ApplyLink              string   `json:"apply_link"`
	GeneralFee             float64  `json:"general_fee"`
	ObcFee                 float64  `json:"obc_fee"`
	ScStFee                float64  `json:"sc_st_fee"`
	WomenFee               float64  `json:"women_fee"`
	IsActive               bool     `json:"is_active"`
}

func createAdminJob(w http.ResponseWriter, r *http.Request) {
	var req JobPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if req.Title == "" || req.Organization == "" || req.LastDate == "" || req.ApplyLink == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required parameters (title, organization, last_date, apply_link)")
		return
	}

	// Insert
	query := `
		INSERT INTO government_jobs (
			title, title_hi, title_mr, organization, organization_hi, organization_mr,
			description, description_hi, description_mr, education_qualification, documents_required,
			min_age, max_age, last_date, apply_link, general_fee, obc_fee, sc_st_fee, women_fee, is_active,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, NOW(), NOW())
		RETURNING id`

	var jobID int
	err = db.DB.QueryRow(query,
		req.Title, req.TitleHi, req.TitleMr, req.Organization, req.OrganizationHi, req.OrganizationMr,
		req.Description, req.DescriptionHi, req.DescriptionMr, req.EducationQualification, pq.StringArray(req.DocumentsRequired),
		req.MinAge, req.MaxAge, req.LastDate, req.ApplyLink, req.GeneralFee, req.ObcFee, req.ScStFee, req.WomenFee, req.IsActive,
	).Scan(&jobID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{"success": true, "id": jobID, "message": "Job advertisement added successfully!"})
}

func updateAdminJob(w http.ResponseWriter, r *http.Request, jobID int) {
	var req JobPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if req.Title == "" || req.Organization == "" || req.LastDate == "" || req.ApplyLink == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	query := `
		UPDATE government_jobs
		SET title = $1, title_hi = $2, title_mr = $3,
		    organization = $4, organization_hi = $5, organization_mr = $6,
		    description = $7, description_hi = $8, description_mr = $9,
		    education_qualification = $10, documents_required = $11,
		    min_age = $12, max_age = $13, last_date = $14, apply_link = $15,
		    general_fee = $16, obc_fee = $17, sc_st_fee = $18, women_fee = $19,
		    is_active = $20, updated_at = NOW()
		WHERE id = $21`

	_, err = db.DB.Exec(query,
		req.Title, req.TitleHi, req.TitleMr, req.Organization, req.OrganizationHi, req.OrganizationMr,
		req.Description, req.DescriptionHi, req.DescriptionMr, req.EducationQualification, pq.StringArray(req.DocumentsRequired),
		req.MinAge, req.MaxAge, req.LastDate, req.ApplyLink, req.GeneralFee, req.ObcFee, req.ScStFee, req.WomenFee, req.IsActive,
		jobID,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Job advertisement updated successfully!"})
}

func deleteAdminJob(w http.ResponseWriter, r *http.Request, jobID int) {
	_, err := db.DB.Exec("DELETE FROM government_jobs WHERE id = $1", jobID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Job advertisement deleted successfully!"})
}
