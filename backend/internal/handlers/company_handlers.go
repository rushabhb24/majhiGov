package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/models"
)

// GetCompaniesHandler returns all active companies with job counts
func GetCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	search := r.URL.Query().Get("search")
	industry := r.URL.Query().Get("industry")

	query := `
		SELECT c.id, c.name, COALESCE(c.slug,''), COALESCE(c.description,''),
		       COALESCE(c.logo_url,''), COALESCE(c.website,''), COALESCE(c.industry,''),
		       COALESCE(c.company_size,''), COALESCE(c.location,''),
		       COALESCE(c.founded_year,0), c.is_verified, c.is_active,
		       c.created_at, c.updated_at,
		       COUNT(pj.id) as job_count
		FROM companies c
		LEFT JOIN private_jobs pj ON pj.company_id = c.id AND pj.is_active = true
		WHERE c.is_active = true`

	var args []interface{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(" AND (c.name ILIKE $%d OR c.industry ILIKE $%d OR c.location ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if industry != "" {
		query += fmt.Sprintf(" AND c.industry ILIKE $%d", argIdx)
		args = append(args, "%"+industry+"%")
		argIdx++
	}

	query += " GROUP BY c.id ORDER BY c.name ASC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var companies []models.Company = []models.Company{}
	for rows.Next() {
		var c models.Company
		err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.Description,
			&c.LogoURL, &c.Website, &c.Industry,
			&c.CompanySize, &c.Location,
			&c.FoundedYear, &c.IsVerified, &c.IsActive,
			&c.CreatedAt, &c.UpdatedAt, &c.JobCount,
		)
		if err != nil {
			continue
		}
		companies = append(companies, c)
	}

	writeJSONResponse(w, http.StatusOK, companies)
}

// GetCompanyDetailsHandler returns a single company along with its active job listings
func GetCompanyDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeJSONError(w, http.StatusBadRequest, "Missing company ID")
		return
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	var c models.Company
	err = db.DB.QueryRow(`
		SELECT id, name, COALESCE(slug,''), COALESCE(description,''),
		       COALESCE(logo_url,''), COALESCE(website,''), COALESCE(industry,''),
		       COALESCE(company_size,''), COALESCE(location,''),
		       COALESCE(founded_year,0), is_verified, is_active, created_at, updated_at
		FROM companies WHERE id = $1 AND is_active = true`, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Description,
		&c.LogoURL, &c.Website, &c.Industry,
		&c.CompanySize, &c.Location,
		&c.FoundedYear, &c.IsVerified, &c.IsActive,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Company not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// Fetch active jobs for this company
	jobRows, err := db.DB.Query(`
		SELECT id, title, job_type, COALESCE(work_mode,''), COALESCE(location,''),
		       COALESCE(salary_min,0), COALESCE(salary_max,0),
		       COALESCE(employment_type,''), COALESCE(experience_min,0), COALESCE(experience_max,5),
		       COALESCE(apply_link,''), created_at
		FROM private_jobs WHERE company_id = $1 AND is_active = true ORDER BY created_at DESC`, id)

	type JobSummary struct {
		ID             int     `json:"id"`
		Title          string  `json:"title"`
		JobType        string  `json:"job_type"`
		WorkMode       string  `json:"work_mode"`
		Location       string  `json:"location"`
		SalaryMin      float64 `json:"salary_min"`
		SalaryMax      float64 `json:"salary_max"`
		EmploymentType string  `json:"employment_type"`
		ExperienceMin  int     `json:"experience_min"`
		ExperienceMax  int     `json:"experience_max"`
		ApplyLink      string  `json:"apply_link"`
		CreatedAt      string  `json:"created_at"`
	}

	var jobs []JobSummary = []JobSummary{}
	if err == nil {
		defer jobRows.Close()
		for jobRows.Next() {
			var j JobSummary
			var createdAt interface{}
			if err := jobRows.Scan(&j.ID, &j.Title, &j.JobType, &j.WorkMode, &j.Location,
				&j.SalaryMin, &j.SalaryMax, &j.EmploymentType, &j.ExperienceMin, &j.ExperienceMax,
				&j.ApplyLink, &createdAt); err == nil {
				j.CreatedAt = fmt.Sprintf("%v", createdAt)
				jobs = append(jobs, j)
			}
		}
		c.JobCount = len(jobs)
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"company": c,
		"jobs":    jobs,
	})
}

// AdminCompaniesHandler handles admin GET (list all) and POST (create) for companies
func AdminCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		adminListCompanies(w, r)
	case "POST":
		adminCreateCompany(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func adminListCompanies(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT c.id, c.name, COALESCE(c.slug,''), COALESCE(c.description,''),
		       COALESCE(c.logo_url,''), COALESCE(c.website,''), COALESCE(c.industry,''),
		       COALESCE(c.company_size,''), COALESCE(c.location,''),
		       COALESCE(c.founded_year,0), c.is_verified, c.is_active,
		       c.created_at, c.updated_at,
		       COUNT(pj.id) as job_count
		FROM companies c
		LEFT JOIN private_jobs pj ON pj.company_id = c.id AND pj.is_active = true
		GROUP BY c.id
		ORDER BY c.created_at DESC`)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var companies []models.Company = []models.Company{}
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.Description,
			&c.LogoURL, &c.Website, &c.Industry,
			&c.CompanySize, &c.Location,
			&c.FoundedYear, &c.IsVerified, &c.IsActive,
			&c.CreatedAt, &c.UpdatedAt, &c.JobCount,
		); err == nil {
			companies = append(companies, c)
		}
	}
	writeJSONResponse(w, http.StatusOK, companies)
}

type CompanyPayload struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	Website     string `json:"website"`
	Industry    string `json:"industry"`
	CompanySize string `json:"company_size"`
	Location    string `json:"location"`
	FoundedYear int    `json:"founded_year"`
	IsVerified  bool   `json:"is_verified"`
	IsActive    bool   `json:"is_active"`
}

func adminCreateCompany(w http.ResponseWriter, r *http.Request) {
	var req CompanyPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "company name is required")
		return
	}

	// Generate slug from name if not provided
	if req.Slug == "" {
		req.Slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}

	var companyID int
	err := db.DB.QueryRow(`
		INSERT INTO companies (name, slug, description, logo_url, website, industry, company_size, location, founded_year, is_verified, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`,
		req.Name, req.Slug, req.Description, req.LogoURL, req.Website,
		req.Industry, req.CompanySize, req.Location, req.FoundedYear,
		req.IsVerified, req.IsActive,
	).Scan(&companyID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create company: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success":    true,
		"company_id": companyID,
		"message":    "Company created successfully",
	})
}

// AdminCompanyDetailsHandler handles PUT (update) and DELETE (toggle) for a single company
func AdminCompanyDetailsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		writeJSONError(w, http.StatusBadRequest, "Missing company ID")
		return
	}
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	switch r.Method {
	case "PUT":
		adminUpdateCompany(w, r, id)
	case "DELETE":
		adminToggleCompany(w, r, id)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func adminUpdateCompany(w http.ResponseWriter, r *http.Request, companyID int) {
	var req CompanyPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, err := db.DB.Exec(`
		UPDATE companies SET
			name=$1, slug=$2, description=$3, logo_url=$4, website=$5,
			industry=$6, company_size=$7, location=$8, founded_year=$9,
			is_verified=$10, is_active=$11, updated_at=NOW()
		WHERE id=$12`,
		req.Name, req.Slug, req.Description, req.LogoURL, req.Website,
		req.Industry, req.CompanySize, req.Location, req.FoundedYear,
		req.IsVerified, req.IsActive, companyID,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update company: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Company updated successfully",
	})
}

func adminToggleCompany(w http.ResponseWriter, r *http.Request, companyID int) {
	var isActive bool
	err := db.DB.QueryRow("SELECT is_active FROM companies WHERE id = $1", companyID).Scan(&isActive)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Company not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}

	newStatus := !isActive
	_, err = db.DB.Exec("UPDATE companies SET is_active = $1, updated_at = NOW() WHERE id = $2", newStatus, companyID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to toggle company status")
		return
	}

	statusText := "deactivated"
	if newStatus {
		statusText = "activated"
	}
	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Company %s successfully", statusText),
		"active":  newStatus,
	})
}

// NewsletterSubscribeHandler handles email newsletter subscription
func NewsletterSubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		writeJSONError(w, http.StatusBadRequest, "Valid email address is required")
		return
	}

	// Validate email format
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		writeJSONError(w, http.StatusBadRequest, "Invalid email address format")
		return
	}

	// Upsert subscriber
	var subID int
	var isNew bool
	err := db.DB.QueryRow(`
		INSERT INTO newsletter_subscribers (email, name, is_active, subscribed_at)
		VALUES ($1, $2, true, NOW())
		ON CONFLICT (email) DO UPDATE SET is_active=true, name=EXCLUDED.name
		RETURNING id, (xmax = 0)`, req.Email, req.Name).Scan(&subID, &isNew)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to process subscription: "+err.Error())
		return
	}

	if isNew {
		SendNewsletterWelcomeEmail(req.Email, req.Name)
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Successfully subscribed to the newsletter!",
	})
}
