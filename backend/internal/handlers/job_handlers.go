package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/models"
)

// GetJobsHandler returns active government jobs, supporting qualification, keyword, and organization filters
func GetJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	qualification := r.URL.Query().Get("qualification")
	search := r.URL.Query().Get("search")
	org := r.URL.Query().Get("organization")

	query := `
		SELECT id, title, title_hi, title_mr, organization, department, vacancies,
		       education_qualification, experience_required, required_documents,
		       application_start_date, application_end_date, official_website, apply_link,
		       application_fee, is_active, created_at, updated_at
		FROM govt_jobs
		WHERE is_active = true`

	var args []interface{}
	argCount := 1

	if qualification != "" {
		// Custom matching: if user education level is provided, show jobs requiring that level or lower
		// e.g. "Graduate" matches "Graduate", "12th Pass", "10th Pass", "None"
		var quals []string
		qLower := strings.ToLower(qualification)
		if qLower == "graduate" || qLower == "post graduate" {
			quals = []string{"None", "Primary", "10th Pass", "12th Pass", "Graduate"}
		} else if qLower == "12th pass" {
			quals = []string{"None", "Primary", "10th Pass", "12th Pass"}
		} else if qLower == "10th pass" {
			quals = []string{"None", "Primary", "10th Pass"}
		} else {
			quals = []string{"None", "Primary"}
		}
		query += fmt.Sprintf(" AND education_qualification = ANY($%d)", argCount)
		args = append(args, pq.Array(quals))
		argCount++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (LOWER(title) LIKE LOWER($%d) OR LOWER(organization) LIKE LOWER($%d) OR LOWER(department) LIKE LOWER($%d))", argCount, argCount, argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	if org != "" {
		query += fmt.Sprintf(" AND LOWER(organization) = LOWER($%d)", argCount)
		args = append(args, org)
		argCount++
	}

	query += " ORDER BY application_end_date ASC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var jobs []models.GovtJob = []models.GovtJob{}
	for rows.Next() {
		var j models.GovtJob
		var startDate, endDate time.Time
		err := rows.Scan(
			&j.ID, &j.Title, &j.TitleHi, &j.TitleMr, &j.Organization, &j.Department, &j.Vacancies,
			&j.EducationQualification, &j.ExperienceRequired, pq.Array(&j.RequiredDocuments),
			&startDate, &endDate, &j.OfficialWebsite, &j.ApplyLink,
			&j.ApplicationFee, &j.IsActive, &j.CreatedAt, &j.UpdatedAt,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to scan job: "+err.Error())
			return
		}
		j.ApplicationStartDate = startDate.Format("2006-01-02")
		j.ApplicationEndDate = endDate.Format("2006-01-02")
		jobs = append(jobs, j)
	}

	writeJSONResponse(w, http.StatusOK, jobs)
}

// GetJobDetailsHandler fetches detailed view parameters for a single Government Job
func GetJobDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeJSONError(w, http.StatusBadRequest, "Missing job ID parameter")
		return
	}
	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid job ID format")
		return
	}

	var j models.GovtJob
	var startDate, endDate time.Time
	query := `
		SELECT id, title, title_hi, title_mr, organization, department, vacancies,
		       education_qualification, experience_required, required_documents,
		       application_start_date, application_end_date, official_website, apply_link,
		       application_fee, is_active, created_at, updated_at
		FROM govt_jobs
		WHERE id = $1`

	err = db.DB.QueryRow(query, id).Scan(
		&j.ID, &j.Title, &j.TitleHi, &j.TitleMr, &j.Organization, &j.Department, &j.Vacancies,
		&j.EducationQualification, &j.ExperienceRequired, pq.Array(&j.RequiredDocuments),
		&startDate, &endDate, &j.OfficialWebsite, &j.ApplyLink,
		&j.ApplicationFee, &j.IsActive, &j.CreatedAt, &j.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Government job advertisement not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	j.ApplicationStartDate = startDate.Format("2006-01-02")
	j.ApplicationEndDate = endDate.Format("2006-01-02")

	writeJSONResponse(w, http.StatusOK, j)
}

// AdminJobsHandler handles listing active/inactive jobs (GET) and creating new jobs (POST)
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

func getAdminAllJobs(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, title, title_hi, title_mr, organization, department, vacancies,
		       education_qualification, experience_required, required_documents,
		       application_start_date, application_end_date, official_website, apply_link,
		       application_fee, is_active, created_at, updated_at
		FROM govt_jobs
		ORDER BY created_at DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to query database: "+err.Error())
		return
	}
	defer rows.Close()

	var jobs []models.GovtJob = []models.GovtJob{}
	for rows.Next() {
		var j models.GovtJob
		var startDate, endDate time.Time
		err := rows.Scan(
			&j.ID, &j.Title, &j.TitleHi, &j.TitleMr, &j.Organization, &j.Department, &j.Vacancies,
			&j.EducationQualification, &j.ExperienceRequired, pq.Array(&j.RequiredDocuments),
			&startDate, &endDate, &j.OfficialWebsite, &j.ApplyLink,
			&j.ApplicationFee, &j.IsActive, &j.CreatedAt, &j.UpdatedAt,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to scan admin job: "+err.Error())
			return
		}
		j.ApplicationStartDate = startDate.Format("2006-01-02")
		j.ApplicationEndDate = endDate.Format("2006-01-02")
		jobs = append(jobs, j)
	}

	writeJSONResponse(w, http.StatusOK, jobs)
}

type JobCreatePayload struct {
	Title                  string   `json:"title"`
	TitleHi                string   `json:"title_hi"`
	TitleMr                string   `json:"title_mr"`
	Organization           string   `json:"organization"`
	Department             string   `json:"department"`
	Vacancies              int      `json:"vacancies"`
	EducationQualification string   `json:"education_qualification"`
	ExperienceRequired     string   `json:"experience_required"`
	RequiredDocuments      []string `json:"required_documents"`
	ApplicationStartDate   string   `json:"application_start_date"`
	ApplicationEndDate     string   `json:"application_end_date"`
	OfficialWebsite        string   `json:"official_website"`
	ApplyLink              string   `json:"apply_link"`
	ApplicationFee         string   `json:"application_fee"`
	IsActive               bool     `json:"is_active"`
}

func createAdminJob(w http.ResponseWriter, r *http.Request) {
	var req JobCreatePayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body payload")
		return
	}

	if req.Title == "" || req.Organization == "" || req.ApplyLink == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing mandatory job criteria parameters")
		return
	}

	// Auto translation fallback triggers:
	// If titles/details for Hindi or Marathi are empty, auto-translate them using our package helper
	if req.TitleHi == "" {
		if hiText, err := translateTextViaGoogle(req.Title, "hi"); err == nil {
			req.TitleHi = hiText
		}
	}
	if req.TitleMr == "" {
		if mrText, err := translateTextViaGoogle(req.Title, "mr"); err == nil {
			req.TitleMr = mrText
		}
	}

	query := `
		INSERT INTO govt_jobs (
			title, title_hi, title_mr, organization, department, vacancies,
			education_qualification, experience_required, required_documents,
			application_start_date, application_end_date, official_website, apply_link,
			application_fee, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`

	var jobID int
	err = db.DB.QueryRow(query,
		req.Title, req.TitleHi, req.TitleMr, req.Organization, req.Department, req.Vacancies,
		req.EducationQualification, req.ExperienceRequired, pq.Array(req.RequiredDocuments),
		req.ApplicationStartDate, req.ApplicationEndDate, req.OfficialWebsite, req.ApplyLink,
		req.ApplicationFee, req.IsActive,
	).Scan(&jobID)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to insert government job: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"job_id":  jobID,
		"message": "Government Job posting successfully created!",
	})
}

// AdminJobDetailsHandler handles job updates (PUT) and toggles/deletes (DELETE)
func AdminJobDetailsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		writeJSONError(w, http.StatusBadRequest, "Missing job ID parameter")
		return
	}
	idStr := parts[4]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid job ID format")
		return
	}

	switch r.Method {
	case "PUT":
		updateAdminJob(w, r, id)
	case "DELETE":
		deleteAdminJob(w, r, id)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func updateAdminJob(w http.ResponseWriter, r *http.Request, jobID int) {
	var req JobCreatePayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	query := `
		UPDATE govt_jobs SET
			title=$1, title_hi=$2, title_mr=$3, organization=$4, department=$5, vacancies=$6,
			education_qualification=$7, experience_required=$8, required_documents=$9,
			application_start_date=$10, application_end_date=$11, official_website=$12, apply_link=$13,
			application_fee=$14, is_active=$15, updated_at=NOW()
		WHERE id=$16`

	_, err = db.DB.Exec(query,
		req.Title, req.TitleHi, req.TitleMr, req.Organization, req.Department, req.Vacancies,
		req.EducationQualification, req.ExperienceRequired, pq.Array(req.RequiredDocuments),
		req.ApplicationStartDate, req.ApplicationEndDate, req.OfficialWebsite, req.ApplyLink,
		req.ApplicationFee, req.IsActive, jobID,
	)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update government job record: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Government Job advertisement successfully updated!",
	})
}

func deleteAdminJob(w http.ResponseWriter, r *http.Request, jobID int) {
	var isActive bool
	err := db.DB.QueryRow("SELECT is_active FROM govt_jobs WHERE id = $1", jobID).Scan(&isActive)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Job posting not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	newStatus := !isActive
	_, err = db.DB.Exec("UPDATE govt_jobs SET is_active = $1, updated_at = NOW() WHERE id = $2", newStatus, jobID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to toggle job deactivation: "+err.Error())
		return
	}

	statusText := "deactivated"
	if newStatus {
		statusText = "activated"
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Government Job successfully %s!", statusText),
		"active":  newStatus,
	})
}
