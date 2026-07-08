package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/middleware"
	"yojana-portal/backend/internal/models"

	"github.com/lib/pq"
)

// GetPrivateJobsHandler returns paginated private/internship/walkin/hackathon jobs with filters
func GetPrivateJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	jobType := r.URL.Query().Get("job_type") // private, internship, walkin, hackathon
	search := r.URL.Query().Get("search")
	location := r.URL.Query().Get("location")
	workMode := r.URL.Query().Get("work_mode")

	page := 1
	limit := 10
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	offset := (page - 1) * limit

	baseQuery := `
		SELECT pj.id, pj.title, COALESCE(pj.title_hi,''), COALESCE(pj.title_mr,''),
		       pj.company_id, COALESCE(c.name,''), COALESCE(c.logo_url,''),
		       COALESCE(pj.description,''), COALESCE(pj.requirements,''), COALESCE(pj.responsibilities,''),
		       COALESCE(pj.salary_min,0), COALESCE(pj.salary_max,0), COALESCE(pj.salary_currency,'INR'),
		       pj.job_type, COALESCE(pj.work_mode,'onsite'), COALESCE(pj.location,''),
		       COALESCE(pj.experience_min,0), COALESCE(pj.experience_max,5),
		       COALESCE(pj.education_qualification,''), COALESCE(pj.skills,'{}'),
		       COALESCE(pj.employment_type,'full-time'),
		       COALESCE(pj.application_start_date::text,''), COALESCE(pj.application_end_date::text,''),
		       COALESCE(pj.official_website,''), COALESCE(pj.apply_link,''),
		       COALESCE(pj.stipend,''), COALESCE(pj.prize_pool,''),
		       pj.is_active, pj.created_at, pj.updated_at
		FROM private_jobs pj
		LEFT JOIN companies c ON pj.company_id = c.id
		WHERE pj.is_active = true`

	var args []interface{}
	argIdx := 1

	if jobType != "" {
		baseQuery += fmt.Sprintf(" AND pj.job_type = $%d", argIdx)
		args = append(args, jobType)
		argIdx++
	}
	if search != "" {
		baseQuery += fmt.Sprintf(" AND (pj.title ILIKE $%d OR pj.description ILIKE $%d OR c.name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if location != "" {
		baseQuery += fmt.Sprintf(" AND pj.location ILIKE $%d", argIdx)
		args = append(args, "%"+location+"%")
		argIdx++
	}
	if workMode != "" {
		baseQuery += fmt.Sprintf(" AND pj.work_mode = $%d", argIdx)
		args = append(args, workMode)
		argIdx++
	}

	// Count query
	countQuery := strings.Replace(baseQuery, `SELECT pj.id, pj.title, COALESCE(pj.title_hi,''), COALESCE(pj.title_mr,''),
		       pj.company_id, COALESCE(c.name,''), COALESCE(c.logo_url,''),
		       COALESCE(pj.description,''), COALESCE(pj.requirements,''), COALESCE(pj.responsibilities,''),
		       COALESCE(pj.salary_min,0), COALESCE(pj.salary_max,0), COALESCE(pj.salary_currency,'INR'),
		       pj.job_type, COALESCE(pj.work_mode,'onsite'), COALESCE(pj.location,''),
		       COALESCE(pj.experience_min,0), COALESCE(pj.experience_max,5),
		       COALESCE(pj.education_qualification,''), COALESCE(pj.skills,'{}'),
		       COALESCE(pj.employment_type,'full-time'),
		       COALESCE(pj.application_start_date::text,''), COALESCE(pj.application_end_date::text,''),
		       COALESCE(pj.official_website,''), COALESCE(pj.apply_link,''),
		       COALESCE(pj.stipend,''), COALESCE(pj.prize_pool,''),
		       pj.is_active, pj.created_at, pj.updated_at`, "SELECT COUNT(*)", 1)

	var total int
	db.DB.QueryRow(countQuery, args...).Scan(&total)

	baseQuery += fmt.Sprintf(" ORDER BY pj.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(baseQuery, args...)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var jobs []models.PrivateJob = []models.PrivateJob{}
	for rows.Next() {
		var j models.PrivateJob
		err := rows.Scan(
			&j.ID, &j.Title, &j.TitleHi, &j.TitleMr,
			&j.CompanyID, &j.CompanyName, &j.CompanyLogoURL,
			&j.Description, &j.Requirements, &j.Responsibilities,
			&j.SalaryMin, &j.SalaryMax, &j.SalaryCurrency,
			&j.JobType, &j.WorkMode, &j.Location,
			&j.ExperienceMin, &j.ExperienceMax,
			&j.EducationQualification, pq.Array(&j.Skills),
			&j.EmploymentType,
			&j.ApplicationStartDate, &j.ApplicationEndDate,
			&j.OfficialWebsite, &j.ApplyLink,
			&j.Stipend, &j.PrizePool,
			&j.IsActive, &j.CreatedAt, &j.UpdatedAt,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to scan job row: "+err.Error())
			return
		}
		if j.Skills == nil {
			j.Skills = []string{}
		}
		jobs = append(jobs, j)
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"data": jobs,
		"meta": map[string]interface{}{
			"page":    page,
			"limit":   limit,
			"total":   total,
			"hasNext": (page * limit) < total,
		},
	})
}

// GetPrivateJobDetailsHandler returns a single private job by ID
func GetPrivateJobDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeJSONError(w, http.StatusBadRequest, "Missing job ID")
		return
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	var j models.PrivateJob
	query := `
		SELECT pj.id, pj.title, COALESCE(pj.title_hi,''), COALESCE(pj.title_mr,''),
		       pj.company_id, COALESCE(c.name,''), COALESCE(c.logo_url,''),
		       COALESCE(pj.description,''), COALESCE(pj.requirements,''), COALESCE(pj.responsibilities,''),
		       COALESCE(pj.salary_min,0), COALESCE(pj.salary_max,0), COALESCE(pj.salary_currency,'INR'),
		       pj.job_type, COALESCE(pj.work_mode,'onsite'), COALESCE(pj.location,''),
		       COALESCE(pj.experience_min,0), COALESCE(pj.experience_max,5),
		       COALESCE(pj.education_qualification,''), COALESCE(pj.skills,'{}'),
		       COALESCE(pj.employment_type,'full-time'),
		       COALESCE(pj.application_start_date::text,''), COALESCE(pj.application_end_date::text,''),
		       COALESCE(pj.official_website,''), COALESCE(pj.apply_link,''),
		       COALESCE(pj.stipend,''), COALESCE(pj.prize_pool,''),
		       pj.is_active, pj.created_at, pj.updated_at
		FROM private_jobs pj
		LEFT JOIN companies c ON pj.company_id = c.id
		WHERE pj.id = $1`

	err = db.DB.QueryRow(query, id).Scan(
		&j.ID, &j.Title, &j.TitleHi, &j.TitleMr,
		&j.CompanyID, &j.CompanyName, &j.CompanyLogoURL,
		&j.Description, &j.Requirements, &j.Responsibilities,
		&j.SalaryMin, &j.SalaryMax, &j.SalaryCurrency,
		&j.JobType, &j.WorkMode, &j.Location,
		&j.ExperienceMin, &j.ExperienceMax,
		&j.EducationQualification, pq.Array(&j.Skills),
		&j.EmploymentType,
		&j.ApplicationStartDate, &j.ApplicationEndDate,
		&j.OfficialWebsite, &j.ApplyLink,
		&j.Stipend, &j.PrizePool,
		&j.IsActive, &j.CreatedAt, &j.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Private job not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if j.Skills == nil {
		j.Skills = []string{}
	}

	writeJSONResponse(w, http.StatusOK, j)
}

// ApplyPrivateJobHandler records user application to a private job
func ApplyPrivateJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		PrivateJobID int `json:"private_job_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PrivateJobID == 0 {
		writeJSONError(w, http.StatusBadRequest, "Invalid private_job_id")
		return
	}

	// Verify job exists and is active
	var applyLink string
	err = db.DB.QueryRow("SELECT COALESCE(apply_link,'') FROM private_jobs WHERE id = $1 AND is_active = true", req.PrivateJobID).Scan(&applyLink)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Private job not found or no longer active")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}

	_, err = db.DB.Exec(
		`INSERT INTO user_applied_private_jobs (user_id, private_job_id, applied_at) VALUES ($1, $2, NOW()) ON CONFLICT (user_id, private_job_id) DO NOTHING`,
		userID, req.PrivateJobID,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to record application")
		return
	}

	LogAuditEvent(userID, "apply_private_job", "private_job", req.PrivateJobID, fmt.Sprintf(`{"job_id":%d}`, req.PrivateJobID), r)

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"apply_link": applyLink,
		"message":    "Application recorded. Redirecting to company portal.",
	})
}

// GetUserPrivateJobApplicationsHandler returns all private job applications for logged-in user
func GetUserPrivateJobApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	type AppRow struct {
		ID          int       `json:"id"`
		PrivateJobID int      `json:"private_job_id"`
		Title       string    `json:"title"`
		CompanyName string    `json:"company_name"`
		JobType     string    `json:"job_type"`
		Location    string    `json:"location"`
		ApplyLink   string    `json:"apply_link"`
		AppliedAt   time.Time `json:"applied_at"`
	}

	rows, err := db.DB.Query(`
		SELECT a.id, a.private_job_id, pj.title, COALESCE(c.name,''), pj.job_type,
		       COALESCE(pj.location,''), COALESCE(pj.apply_link,''), a.applied_at
		FROM user_applied_private_jobs a
		JOIN private_jobs pj ON a.private_job_id = pj.id
		LEFT JOIN companies c ON pj.company_id = c.id
		WHERE a.user_id = $1
		ORDER BY a.applied_at DESC`, userID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var apps []AppRow = []AppRow{}
	for rows.Next() {
		var a AppRow
		if err := rows.Scan(&a.ID, &a.PrivateJobID, &a.Title, &a.CompanyName, &a.JobType, &a.Location, &a.ApplyLink, &a.AppliedAt); err == nil {
			apps = append(apps, a)
		}
	}

	writeJSONResponse(w, http.StatusOK, apps)
}

// AdminPrivateJobsHandler handles admin GET (list all) and POST (create) for private jobs
func AdminPrivateJobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		adminListPrivateJobs(w, r)
	case "POST":
		adminCreatePrivateJob(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func adminListPrivateJobs(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT pj.id, pj.title, COALESCE(pj.title_hi,''), COALESCE(pj.title_mr,''),
		       pj.company_id, COALESCE(c.name,''), COALESCE(c.logo_url,''),
		       COALESCE(pj.description,''), COALESCE(pj.requirements,''), COALESCE(pj.responsibilities,''),
		       COALESCE(pj.salary_min,0), COALESCE(pj.salary_max,0), COALESCE(pj.salary_currency,'INR'),
		       pj.job_type, COALESCE(pj.work_mode,'onsite'), COALESCE(pj.location,''),
		       COALESCE(pj.experience_min,0), COALESCE(pj.experience_max,5),
		       COALESCE(pj.education_qualification,''), COALESCE(pj.skills,'{}'),
		       COALESCE(pj.employment_type,'full-time'),
		       COALESCE(pj.application_start_date::text,''), COALESCE(pj.application_end_date::text,''),
		       COALESCE(pj.official_website,''), COALESCE(pj.apply_link,''),
		       COALESCE(pj.stipend,''), COALESCE(pj.prize_pool,''),
		       pj.is_active, pj.created_at, pj.updated_at
		FROM private_jobs pj
		LEFT JOIN companies c ON pj.company_id = c.id
		ORDER BY pj.created_at DESC`)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var jobs []models.PrivateJob = []models.PrivateJob{}
	for rows.Next() {
		var j models.PrivateJob
		err := rows.Scan(
			&j.ID, &j.Title, &j.TitleHi, &j.TitleMr,
			&j.CompanyID, &j.CompanyName, &j.CompanyLogoURL,
			&j.Description, &j.Requirements, &j.Responsibilities,
			&j.SalaryMin, &j.SalaryMax, &j.SalaryCurrency,
			&j.JobType, &j.WorkMode, &j.Location,
			&j.ExperienceMin, &j.ExperienceMax,
			&j.EducationQualification, pq.Array(&j.Skills),
			&j.EmploymentType,
			&j.ApplicationStartDate, &j.ApplicationEndDate,
			&j.OfficialWebsite, &j.ApplyLink,
			&j.Stipend, &j.PrizePool,
			&j.IsActive, &j.CreatedAt, &j.UpdatedAt,
		)
		if err != nil {
			continue
		}
		if j.Skills == nil {
			j.Skills = []string{}
		}
		jobs = append(jobs, j)
	}
	writeJSONResponse(w, http.StatusOK, jobs)
}

type PrivateJobPayload struct {
	Title                  string   `json:"title"`
	TitleHi                string   `json:"title_hi"`
	TitleMr                string   `json:"title_mr"`
	CompanyID              *int     `json:"company_id"`
	Description            string   `json:"description"`
	Requirements           string   `json:"requirements"`
	Responsibilities       string   `json:"responsibilities"`
	SalaryMin              float64  `json:"salary_min"`
	SalaryMax              float64  `json:"salary_max"`
	SalaryCurrency         string   `json:"salary_currency"`
	JobType                string   `json:"job_type"`
	WorkMode               string   `json:"work_mode"`
	Location               string   `json:"location"`
	ExperienceMin          int      `json:"experience_min"`
	ExperienceMax          int      `json:"experience_max"`
	EducationQualification string   `json:"education_qualification"`
	Skills                 []string `json:"skills"`
	EmploymentType         string   `json:"employment_type"`
	ApplicationStartDate   string   `json:"application_start_date"`
	ApplicationEndDate     string   `json:"application_end_date"`
	OfficialWebsite        string   `json:"official_website"`
	ApplyLink              string   `json:"apply_link"`
	Stipend                string   `json:"stipend"`
	PrizePool              string   `json:"prize_pool"`
	IsActive               bool     `json:"is_active"`
}

func adminCreatePrivateJob(w http.ResponseWriter, r *http.Request) {
	var req PrivateJobPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Title == "" || req.JobType == "" {
		writeJSONError(w, http.StatusBadRequest, "title and job_type are required")
		return
	}
	if req.SalaryCurrency == "" {
		req.SalaryCurrency = "INR"
	}

	var jobID int
	err := db.DB.QueryRow(`
		INSERT INTO private_jobs (
			title, title_hi, title_mr, company_id, description, requirements, responsibilities,
			salary_min, salary_max, salary_currency, job_type, work_mode, location,
			experience_min, experience_max, education_qualification, skills, employment_type,
			application_start_date, application_end_date, official_website, apply_link,
			stipend, prize_pool, is_active
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,
		          NULLIF($19,'')::date, NULLIF($20,'')::date, $21,$22,$23,$24,$25)
		RETURNING id`,
		req.Title, req.TitleHi, req.TitleMr, req.CompanyID,
		req.Description, req.Requirements, req.Responsibilities,
		req.SalaryMin, req.SalaryMax, req.SalaryCurrency,
		req.JobType, req.WorkMode, req.Location,
		req.ExperienceMin, req.ExperienceMax,
		req.EducationQualification, pq.Array(req.Skills), req.EmploymentType,
		req.ApplicationStartDate, req.ApplicationEndDate,
		req.OfficialWebsite, req.ApplyLink,
		req.Stipend, req.PrizePool, req.IsActive,
	).Scan(&jobID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create private job: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success":        true,
		"private_job_id": jobID,
		"message":        "Private job created successfully",
	})
}

// AdminPrivateJobDetailsHandler handles PUT (update) and DELETE (toggle) for a single private job
func AdminPrivateJobDetailsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		writeJSONError(w, http.StatusBadRequest, "Missing job ID")
		return
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	switch r.Method {
	case "PUT":
		adminUpdatePrivateJob(w, r, id)
	case "DELETE":
		adminTogglePrivateJob(w, r, id)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func adminUpdatePrivateJob(w http.ResponseWriter, r *http.Request, jobID int) {
	var req PrivateJobPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, err := db.DB.Exec(`
		UPDATE private_jobs SET
			title=$1, title_hi=$2, title_mr=$3, company_id=$4,
			description=$5, requirements=$6, responsibilities=$7,
			salary_min=$8, salary_max=$9, salary_currency=$10,
			job_type=$11, work_mode=$12, location=$13,
			experience_min=$14, experience_max=$15,
			education_qualification=$16, skills=$17, employment_type=$18,
			application_start_date=NULLIF($19,'')::date,
			application_end_date=NULLIF($20,'')::date,
			official_website=$21, apply_link=$22,
			stipend=$23, prize_pool=$24, is_active=$25, updated_at=NOW()
		WHERE id=$26`,
		req.Title, req.TitleHi, req.TitleMr, req.CompanyID,
		req.Description, req.Requirements, req.Responsibilities,
		req.SalaryMin, req.SalaryMax, req.SalaryCurrency,
		req.JobType, req.WorkMode, req.Location,
		req.ExperienceMin, req.ExperienceMax,
		req.EducationQualification, pq.Array(req.Skills), req.EmploymentType,
		req.ApplicationStartDate, req.ApplicationEndDate,
		req.OfficialWebsite, req.ApplyLink,
		req.Stipend, req.PrizePool, req.IsActive, jobID,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update private job: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Private job updated successfully",
	})
}

func adminTogglePrivateJob(w http.ResponseWriter, r *http.Request, jobID int) {
	var isActive bool
	err := db.DB.QueryRow("SELECT is_active FROM private_jobs WHERE id = $1", jobID).Scan(&isActive)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Private job not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}

	newStatus := !isActive
	_, err = db.DB.Exec("UPDATE private_jobs SET is_active = $1, updated_at = NOW() WHERE id = $2", newStatus, jobID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to toggle private job")
		return
	}

	statusText := "deactivated"
	if newStatus {
		statusText = "activated"
	}
	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Private job %s successfully", statusText),
		"active":  newStatus,
	})
}
