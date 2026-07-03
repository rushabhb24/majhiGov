package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/models"
)

// Helper to write JSON error responses
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Helper to write JSON success responses
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// GetAdminAnalyticsHandler aggregates operational metrics for dashboard overview
func GetAdminAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var totalSchemes, totalUsers, totalCategories, expiringSchemes int
	var totalApplications, pendingApplications, approvedApplications, rejectedApplications int

	// 1. Total Schemes
	err := db.DB.QueryRow("SELECT COUNT(*) FROM schemes").Scan(&totalSchemes)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// 2. Total Registered Users
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = false").Scan(&totalUsers)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// 3. Total Categories
	err = db.DB.QueryRow("SELECT COUNT(*) FROM scheme_categories").Scan(&totalCategories)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// 4. Schemes Expiring within 30 days
	err = db.DB.QueryRow("SELECT COUNT(*) FROM schemes WHERE application_end_date <= CURRENT_DATE + INTERVAL '30 days' AND application_end_date >= CURRENT_DATE AND is_active = true").Scan(&expiringSchemes)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// 4a. Total Applications
	err = db.DB.QueryRow("SELECT COUNT(*) FROM user_applied_schemes").Scan(&totalApplications)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// 4b. Pending Applications
	err = db.DB.QueryRow("SELECT COUNT(*) FROM user_applied_schemes WHERE status = 'pending'").Scan(&pendingApplications)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// 4c. Approved Applications
	err = db.DB.QueryRow("SELECT COUNT(*) FROM user_applied_schemes WHERE status = 'approved'").Scan(&approvedApplications)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// 4d. Rejected Applications
	err = db.DB.QueryRow("SELECT COUNT(*) FROM user_applied_schemes WHERE status = 'rejected'").Scan(&rejectedApplications)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	// 5. Schemes by Category count
	type CatCount struct {
		Name   string `json:"name"`
		Icon   string `json:"icon"`
		Count  int    `json:"count"`
		Percent int   `json:"percent"`
	}
	catRows, err := db.DB.Query(`
		SELECT c.name, c.icon, COUNT(s.id) 
		FROM scheme_categories c 
		LEFT JOIN schemes s ON c.id = s.category_id 
		GROUP BY c.id, c.name, c.icon 
		ORDER BY COUNT(s.id) DESC`)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer catRows.Close()

	var schemesByCategory []CatCount = []CatCount{}
	for catRows.Next() {
		var cc CatCount
		if err := catRows.Scan(&cc.Name, &cc.Icon, &cc.Count); err == nil {
			if totalSchemes > 0 {
				cc.Percent = (cc.Count * 100) / totalSchemes
			}
			schemesByCategory = append(schemesByCategory, cc)
		}
	}

	// 6. Top Applied Schemes
	type TopScheme struct {
		Title string `json:"title"`
		Count int    `json:"count"`
	}
	topRows, err := db.DB.Query(`
		SELECT s.title, COUNT(a.id) 
		FROM schemes s 
		JOIN user_applied_schemes a ON s.id = a.scheme_id 
		GROUP BY s.id, s.title 
		ORDER BY COUNT(a.id) DESC 
		LIMIT 5`)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer topRows.Close()

	var topAppliedSchemes []TopScheme = []TopScheme{}
	for topRows.Next() {
		var ts TopScheme
		if err := topRows.Scan(&ts.Title, &ts.Count); err == nil {
			topAppliedSchemes = append(topAppliedSchemes, ts)
		}
	}

	// 6a. Recent Applications
	type RecentApp struct {
		ID         int       `json:"id"`
		FullName   string    `json:"full_name"`
		SchemeName string    `json:"scheme_name"`
		Status     string    `json:"status"`
		AppliedAt  time.Time `json:"applied_at"`
		TimeAgo    string    `json:"time_ago"`
	}
	appRows, err := db.DB.Query(`
		SELECT a.id, p.full_name, s.title, a.status, a.applied_at 
		FROM user_applied_schemes a 
		JOIN users u ON a.user_id = u.id 
		JOIN user_profiles p ON u.id = p.user_id 
		JOIN schemes s ON a.scheme_id = s.id 
		ORDER BY a.applied_at DESC 
		LIMIT 5`)
	
	var recentApps []RecentApp = []RecentApp{}
	if err == nil {
		defer appRows.Close()
		now := time.Now()
		for appRows.Next() {
			var ra RecentApp
			if err := appRows.Scan(&ra.ID, &ra.FullName, &ra.SchemeName, &ra.Status, &ra.AppliedAt); err == nil {
				diff := now.Sub(ra.AppliedAt)
				if diff.Hours() < 1 {
					ra.TimeAgo = fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
				} else if diff.Hours() < 24 {
					ra.TimeAgo = fmt.Sprintf("%d hours ago", int(diff.Hours()))
				} else {
					ra.TimeAgo = fmt.Sprintf("%d days ago", int(diff.Hours()/24))
				}
				recentApps = append(recentApps, ra)
			}
		}
	}

	// recentApps stays as empty slice if no real applications exist

	// 7. Recent Activity logs compiled from database states
	type ActivityItem struct {
		Type    string    `json:"type"`
		Text    string    `json:"text"`
		Created time.Time `json:"created_at"`
		TimeAgo string    `json:"time_ago"`
	}
	
	var activities []ActivityItem = []ActivityItem{}

	// Query new schemes added
	newSchemeRows, err := db.DB.Query(`
		SELECT s.title, s.created_at, c.name 
		FROM schemes s 
		JOIN scheme_categories c ON s.category_id = c.id 
		ORDER BY s.created_at DESC LIMIT 4`)
	if err == nil {
		defer newSchemeRows.Close()
		for newSchemeRows.Next() {
			var title, catName string
			var created time.Time
			if err := newSchemeRows.Scan(&title, &created, &catName); err == nil {
				activities = append(activities, ActivityItem{
					Type:    "scheme",
					Text:    fmt.Sprintf("New scheme '%s' added to %s category", title, catName),
					Created: created,
				})
			}
		}
	}

	// Query user registrations
	userRegRows, err := db.DB.Query(`
		SELECT p.full_name, p.state, u.created_at 
		FROM users u 
		JOIN user_profiles p ON u.id = p.user_id 
		WHERE u.is_admin = false 
		ORDER BY u.created_at DESC LIMIT 4`)
	if err == nil {
		defer userRegRows.Close()
		for userRegRows.Next() {
			var fullName, state string
			var created time.Time
			if err := userRegRows.Scan(&fullName, &state, &created); err == nil {
				activities = append(activities, ActivityItem{
					Type:    "user",
					Text:    fmt.Sprintf("User '%s' registered from %s", fullName, state),
					Created: created,
				})
			}
		}
	}

	// Format activity time_ago in a simple format
	now := time.Now()
	for i := range activities {
		diff := now.Sub(activities[i].Created)
		if diff.Hours() < 1 {
			activities[i].TimeAgo = fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
		} else if diff.Hours() < 24 {
			activities[i].TimeAgo = fmt.Sprintf("%d hours ago", int(diff.Hours()))
		} else {
			activities[i].TimeAgo = fmt.Sprintf("%d days ago", int(diff.Hours()/24))
		}
	}

	// activities stays as empty slice if no real activity exists

	response := map[string]interface{}{
		"total_schemes":         totalSchemes,
		"total_users":           totalUsers,
		"total_categories":      totalCategories,
		"expiring_schemes":      expiringSchemes,
		"total_applications":    totalApplications,
		"pending_applications":   pendingApplications,
		"approved_applications":  approvedApplications,
		"rejected_applications":  rejectedApplications,
		"schemes_by_category":   schemesByCategory,
		"top_schemes":           topAppliedSchemes,
		"recent_activity":       activities,
		"recent_applications":   recentApps,
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// AdminSchemesHandler handles GET (list all) and POST (create scheme)
func AdminSchemesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getAdminAllSchemes(w, r)
	case "POST":
		createAdminScheme(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func getAdminAllSchemes(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT s.id, s.title, s.title_hi, s.title_mr, s.description, s.description_hi, s.description_mr,
		       s.category_id, c.name as category_name, s.government_level, s.state, s.benefits,
		       s.application_start_date, s.application_end_date, s.official_website, s.apply_link, s.is_active, s.created_at
		FROM schemes s
		JOIN scheme_categories c ON s.category_id = c.id
		ORDER BY s.created_at DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to query database: "+err.Error())
		return
	}
	defer rows.Close()

	var schemes []models.Scheme = []models.Scheme{}
	for rows.Next() {
		var s models.Scheme
		var startDate, endDate time.Time
		err := rows.Scan(
			&s.ID, &s.Title, &s.TitleHi, &s.TitleMr, &s.Description, &s.DescriptionHi, &s.DescriptionMr,
			&s.CategoryID, &s.CategoryName, &s.GovernmentLevel, &s.State, &s.Benefits,
			&startDate, &endDate, &s.OfficialWebsite, &s.ApplyLink, &s.IsActive, &s.CreatedAt,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to scan scheme: "+err.Error())
			return
		}
		s.ApplicationStartDate = startDate.Format("2006-01-02")
		s.ApplicationEndDate = endDate.Format("2006-01-02")
		schemes = append(schemes, s)
	}

	writeJSONResponse(w, http.StatusOK, schemes)
}

type SchemeCreatePayload struct {
	Title                string                     `json:"title"`
	TitleHi              string                     `json:"title_hi"`
	TitleMr              string                     `json:"title_mr"`
	Description          string                     `json:"description"`
	DescriptionHi        string                     `json:"description_hi"`
	DescriptionMr        string                     `json:"description_mr"`
	CategoryID           int                        `json:"category_id"`
	GovernmentLevel      string                     `json:"government_level"`
	State                *string                    `json:"state"`
	Benefits             string                     `json:"benefits"`
	ApplicationStartDate string                     `json:"application_start_date"`
	ApplicationEndDate   string                     `json:"application_end_date"`
	OfficialWebsite      string                     `json:"official_website"`
	ApplyLink            string                     `json:"apply_link"`
	IsActive             bool                       `json:"is_active"`
	Eligibility          models.EligibilityCriteria `json:"eligibility"`
	Documents            []models.SchemeDocument    `json:"documents"`
	FAQs                 []models.SchemeFAQ         `json:"faqs"`
}

func createAdminScheme(w http.ResponseWriter, r *http.Request) {
	var req SchemeCreatePayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.Title == "" || req.CategoryID == 0 || req.GovernmentLevel == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required parameters (title, category_id, government_level)")
		return
	}

	// Start Transaction
	tx, err := db.DB.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Transaction failed: "+err.Error())
		return
	}
	defer tx.Rollback()

	// Insert Scheme
	var schemeID int
	queryScheme := `
		INSERT INTO schemes (title, title_hi, title_mr, description, description_hi, description_mr,
		                     category_id, government_level, state, benefits, application_start_date,
		                     application_end_date, official_website, apply_link, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW()) RETURNING id`

	err = tx.QueryRow(queryScheme,
		req.Title, req.TitleHi, req.TitleMr,
		req.Description, req.DescriptionHi, req.DescriptionMr,
		req.CategoryID, req.GovernmentLevel, req.State, req.Benefits,
		req.ApplicationStartDate, req.ApplicationEndDate,
		req.OfficialWebsite, req.ApplyLink, req.IsActive,
	).Scan(&schemeID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to insert scheme: "+err.Error())
		return
	}

	// Insert Eligibility Rules
	queryElig := `
		INSERT INTO eligibility_criteria (scheme_id, min_age, max_age, gender, caste_categories,
		                                 min_income, max_income, states, occupations, employee_types,
		                                 education_levels, disability_required, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())`

	_, err = tx.Exec(queryElig,
		schemeID, req.Eligibility.MinAge, req.Eligibility.MaxAge, req.Eligibility.Gender,
		pq.Array(req.Eligibility.CasteCategories), req.Eligibility.MinIncome, req.Eligibility.MaxIncome,
		pq.Array(req.Eligibility.States), pq.Array(req.Eligibility.Occupations), pq.Array(req.Eligibility.EmployeeTypes),
		pq.Array(req.Eligibility.EducationLevels), req.Eligibility.DisabilityRequired,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to insert eligibility: "+err.Error())
		return
	}

	// Insert Documents
	for _, doc := range req.Documents {
		if doc.DocumentName == "" {
			continue
		}
		queryDoc := `
			INSERT INTO scheme_documents (scheme_id, document_name, document_name_hi, document_name_mr, is_mandatory, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())`
		_, err = tx.Exec(queryDoc, schemeID, doc.DocumentName, doc.DocumentNameHi, doc.DocumentNameMr, doc.IsMandatory)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to insert document: "+err.Error())
			return
		}
	}

	// Insert FAQs
	for _, faq := range req.FAQs {
		if faq.Question == "" {
			continue
		}
		queryFAQ := `
			INSERT INTO scheme_faqs (scheme_id, question, answer, question_hi, answer_hi, question_mr, answer_mr, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
		_, err = tx.Exec(queryFAQ, schemeID, faq.Question, faq.Answer, faq.QuestionHi, faq.AnswerHi, faq.QuestionMr, faq.AnswerMr)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to insert FAQ: "+err.Error())
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to commit scheme: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success":   true,
		"scheme_id": schemeID,
		"message":   "Scheme and eligibility parameters successfully created!",
	})
}

// AdminSchemeDetailsHandler handles PUT (update scheme) and DELETE (deactivate/toggle)
func AdminSchemeDetailsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		writeJSONError(w, http.StatusBadRequest, "Missing scheme ID parameter")
		return
	}
	idStr := parts[4]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid scheme ID format")
		return
	}

	switch r.Method {
	case "PUT":
		updateAdminScheme(w, r, id)
	case "DELETE":
		deleteAdminScheme(w, r, id)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func updateAdminScheme(w http.ResponseWriter, r *http.Request, schemeID int) {
	var req SchemeCreatePayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Transaction initialization failed: "+err.Error())
		return
	}
	defer tx.Rollback()

	// Update Scheme main record
	queryScheme := `
		UPDATE schemes SET title=$1, title_hi=$2, title_mr=$3, description=$4, description_hi=$5, description_mr=$6,
		                   category_id=$7, government_level=$8, state=$9, benefits=$10, application_start_date=$11,
		                   application_end_date=$12, official_website=$13, apply_link=$14, is_active=$15, updated_at=NOW()
		WHERE id=$16`

	_, err = tx.Exec(queryScheme,
		req.Title, req.TitleHi, req.TitleMr,
		req.Description, req.DescriptionHi, req.DescriptionMr,
		req.CategoryID, req.GovernmentLevel, req.State, req.Benefits,
		req.ApplicationStartDate, req.ApplicationEndDate,
		req.OfficialWebsite, req.ApplyLink, req.IsActive, schemeID,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update scheme record: "+err.Error())
		return
	}

	// Update Eligibility (UPSERT style)
	queryElig := `
		INSERT INTO eligibility_criteria (scheme_id, min_age, max_age, gender, caste_categories,
		                                 min_income, max_income, states, occupations, employee_types,
		                                 education_levels, disability_required, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
		ON CONFLICT (scheme_id) DO UPDATE SET
			min_age=EXCLUDED.min_age, max_age=EXCLUDED.max_age, gender=EXCLUDED.gender,
			caste_categories=EXCLUDED.caste_categories, min_income=EXCLUDED.min_income, max_income=EXCLUDED.max_income,
			states=EXCLUDED.states, occupations=EXCLUDED.occupations, employee_types=EXCLUDED.employee_types,
			education_levels=EXCLUDED.education_levels, disability_required=EXCLUDED.disability_required`

	_, err = tx.Exec(queryElig,
		schemeID, req.Eligibility.MinAge, req.Eligibility.MaxAge, req.Eligibility.Gender,
		pq.Array(req.Eligibility.CasteCategories), req.Eligibility.MinIncome, req.Eligibility.MaxIncome,
		pq.Array(req.Eligibility.States), pq.Array(req.Eligibility.Occupations), pq.Array(req.Eligibility.EmployeeTypes),
		pq.Array(req.Eligibility.EducationLevels), req.Eligibility.DisabilityRequired,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update scheme eligibility rules: "+err.Error())
		return
	}

	// Clear previous documents and FAQs to write fresh arrays
	_, _ = tx.Exec("DELETE FROM scheme_documents WHERE scheme_id=$1", schemeID)
	_, _ = tx.Exec("DELETE FROM scheme_faqs WHERE scheme_id=$1", schemeID)

	// Write Documents
	for _, doc := range req.Documents {
		if doc.DocumentName == "" {
			continue
		}
		queryDoc := `
			INSERT INTO scheme_documents (scheme_id, document_name, document_name_hi, document_name_mr, is_mandatory, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())`
		_, err = tx.Exec(queryDoc, schemeID, doc.DocumentName, doc.DocumentNameHi, doc.DocumentNameMr, doc.IsMandatory)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to update documents: "+err.Error())
			return
		}
	}

	// Write FAQs
	for _, faq := range req.FAQs {
		if faq.Question == "" {
			continue
		}
		queryFAQ := `
			INSERT INTO scheme_faqs (scheme_id, question, answer, question_hi, answer_hi, question_mr, answer_mr, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`
		_, err = tx.Exec(queryFAQ, schemeID, faq.Question, faq.Answer, faq.QuestionHi, faq.AnswerHi, faq.QuestionMr, faq.AnswerMr)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to update FAQs: "+err.Error())
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to commit transaction: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Scheme, FAQ, documents, and eligibility criteria updated successfully!",
	})
}

func deleteAdminScheme(w http.ResponseWriter, r *http.Request, schemeID int) {
	// Instead of hard deleting, we toggle the is_active status of the scheme.
	// This prevents breaking foreign keys if user_applied_schemes rows exist in DB.
	var isActive bool
	err := db.DB.QueryRow("SELECT is_active FROM schemes WHERE id = $1", schemeID).Scan(&isActive)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "Scheme not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	newStatus := !isActive
	_, err = db.DB.Exec("UPDATE schemes SET is_active = $1, updated_at = NOW() WHERE id = $2", newStatus, schemeID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update scheme status: "+err.Error())
		return
	}

	statusText := "deactivated"
	if newStatus {
		statusText = "activated"
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Scheme successfully %s!", statusText),
		"active":  newStatus,
	})
}

// AdminCategoriesHandler handles category fetching and creation
func AdminCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getAdminCategories(w, r)
	case "POST":
		createAdminCategory(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func getAdminCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, name_hi, name_mr, icon, description FROM scheme_categories ORDER BY id ASC")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var categories []models.SchemeCategory = []models.SchemeCategory{}
	for rows.Next() {
		var c models.SchemeCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.NameHi, &c.NameMr, &c.Icon, &c.Description); err == nil {
			categories = append(categories, c)
		}
	}

	writeJSONResponse(w, http.StatusOK, categories)
}

func createAdminCategory(w http.ResponseWriter, r *http.Request) {
	var req models.SchemeCategory
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.Name == "" || req.Icon == "" || req.Description == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required parameters (name, icon, description)")
		return
	}

	var categoryID int
	query := `
		INSERT INTO scheme_categories (name, name_hi, name_mr, icon, description, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`

	err = db.DB.QueryRow(query, req.Name, req.NameHi, req.NameMr, req.Icon, req.Description).Scan(&categoryID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create category: "+err.Error())
		return
	}

	req.ID = categoryID
	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success":  true,
		"category": req,
		"message":  "Category successfully created!",
	})
}

// AdminCategoryDetailsHandler handles category deletes
func AdminCategoryDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		writeJSONError(w, http.StatusBadRequest, "Missing category ID parameter")
		return
	}
	idStr := parts[4]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	// Verify no schemes are bound to category first
	var schemesCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM schemes WHERE category_id = $1", id).Scan(&schemesCount)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	if schemesCount > 0 {
		writeJSONError(w, http.StatusBadRequest, "Cannot delete category: schemes are currently bound to this category. Delete/reassign schemes first.")
		return
	}

	_, err = db.DB.Exec("DELETE FROM scheme_categories WHERE id = $1", id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete category: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Category deleted successfully!",
	})
}

// AdminUsersHandler fetches user listings and statuses
func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	type UserSummary struct {
		ID            int       `json:"id"`
		Email         string    `json:"email"`
		Phone         string    `json:"phone"`
		IsVerified    bool      `json:"is_verified"`
		IsAdmin       bool      `json:"is_admin"`
		CreatedAt     time.Time `json:"created_at"`
		FullName      string    `json:"full_name"`
		State         string    `json:"state"`
		District      string    `json:"district"`
		Occupation    string    `json:"occupation"`
		CasteCategory string    `json:"caste_category"`
	}

	query := `
		SELECT u.id, u.email, u.phone, u.is_verified, u.is_admin, u.created_at,
		       p.full_name, p.state, p.district, p.occupation, p.caste_category
		FROM users u
		JOIN user_profiles p ON u.id = p.user_id
		ORDER BY u.created_at DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var users []UserSummary = []UserSummary{}
	for rows.Next() {
		var us UserSummary
		err := rows.Scan(
			&us.ID, &us.Email, &us.Phone, &us.IsVerified, &us.IsAdmin, &us.CreatedAt,
			&us.FullName, &us.State, &us.District, &us.Occupation, &us.CasteCategory,
		)
		if err == nil {
			users = append(users, us)
		}
	}

	writeJSONResponse(w, http.StatusOK, users)
}

// AdminUserToggleHandler toggles user verification (active/inactive status)
func AdminUserToggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.UserID == 0 {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body payload")
		return
	}

	var isVerified bool
	err = db.DB.QueryRow("SELECT is_verified FROM users WHERE id = $1", req.UserID).Scan(&isVerified)
	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	newStatus := !isVerified
	_, err = db.DB.Exec("UPDATE users SET is_verified = $1 WHERE id = $2", newStatus, req.UserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to toggle user verification: "+err.Error())
		return
	}

	statusText := "banned / inactivated"
	if newStatus {
		statusText = "restored / verified"
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("User status successfully changed to %s!", statusText),
		"verified": newStatus,
	})
}

// AdminCreateHandler creates a new administrative user account
func AdminCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if req.Email == "" || req.Phone == "" || req.Password == "" || req.FullName == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required administrative fields")
		return
	}

	// Uniqueness check
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR phone = $2)", req.Email, req.Phone).Scan(&exists)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if exists {
		writeJSONError(w, http.StatusBadRequest, "Email or phone number is already registered")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to secure password: "+err.Error())
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Transaction failed: "+err.Error())
		return
	}
	defer tx.Rollback()

	var adminID int
	queryUser := `
		INSERT INTO users (email, phone, password_hash, is_verified, is_admin, created_at)
		VALUES ($1, $2, $3, true, true, NOW()) RETURNING id`

	err = tx.QueryRow(queryUser, req.Email, req.Phone, string(hash)).Scan(&adminID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to insert user: "+err.Error())
		return
	}

	queryProfile := `
		INSERT INTO user_profiles (user_id, full_name, date_of_birth, gender, state, district,
		                          caste_category, annual_income, occupation, employee_type,
		                          education_level, is_disabled, updated_at)
		VALUES ($1, $2, '1990-01-01', 'Male', 'Maharashtra', 'Mumbai',
		        'General', 0.00, 'Other', 'Government Employee', 'Graduate', false, NOW())`

	_, err = tx.Exec(queryProfile, adminID, req.FullName)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to insert profile: "+err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to commit transaction: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "New Administrator successfully registered!",
	})
}

// AdminNotificationsHandler broadcasts notifications to dynamic user cohorts
func AdminNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getAdminNotifications(w, r)
		return
	}

	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		SendTo  string `json:"send_to"` // "All Users", "All Farmers", "All Students", "Specific State"
		State   string `json:"state"`   // Used if send_to is "Specific State"
		Title   string `json:"title"`
		Message string `json:"message"`
		Type    string `json:"type"` // "Deadline Reminder", "New Scheme Alert", "System Update"
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if req.SendTo == "" || req.Title == "" || req.Message == "" || req.Type == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required parameters (send_to, title, message, type)")
		return
	}

	// 1. Identify Target Users
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT u.id FROM users u JOIN user_profiles p ON u.id = p.user_id WHERE u.is_admin = false")

	var args []interface{}
	switch req.SendTo {
	case "All Farmers":
		queryBuilder.WriteString(" AND p.occupation = 'Farmer'")
	case "All Students":
		queryBuilder.WriteString(" AND p.occupation = 'Student'")
	case "Specific State":
		queryBuilder.WriteString(" AND p.state = $1")
		args = append(args, req.State)
	}

	rows, err := db.DB.Query(queryBuilder.String(), args...)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch targets: "+err.Error())
		return
	}
	defer rows.Close()

	var targetIDs []int = []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			targetIDs = append(targetIDs, id)
		}
	}

	if len(targetIDs) == 0 {
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Notification broadcast completed! Target list was empty, no records inserted.",
			"count":   0,
		})
		return
	}

	// 2. Broadcast in Transaction
	tx, err := db.DB.Begin()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Transaction failed: "+err.Error())
		return
	}
	defer tx.Rollback()

	queryNotify := `
		INSERT INTO notifications (user_id, title, message, type, is_read, created_at)
		VALUES ($1, $2, $3, $4, false, NOW())`

	for _, uid := range targetIDs {
		_, err = tx.Exec(queryNotify, uid, req.Title, req.Message, req.Type)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to insert notification: "+err.Error())
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to broadcast: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Notification broadcast successfully sent to %d users!", len(targetIDs)),
		"count":   len(targetIDs),
	})
}

func getAdminNotifications(w http.ResponseWriter, r *http.Request) {
	// Returns a synthesized log list of recent broadcast activity for presentation purposes
	type SentNotification struct {
		Title   string    `json:"title"`
		Message string    `json:"message"`
		Type    string    `json:"type"`
		Target  string    `json:"target"`
		Created time.Time `json:"created_at"`
		TimeAgo string    `json:"time_ago"`
	}

	// We can group recent notifications to represent historic broadcast triggers
	rows, err := db.DB.Query(`
		SELECT title, message, type, COUNT(user_id) as count, MAX(created_at) as created
		FROM notifications 
		GROUP BY title, message, type 
		ORDER BY created DESC 
		LIMIT 6`)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	var list []SentNotification = []SentNotification{}
	now := time.Now()
	for rows.Next() {
		var sn SentNotification
		var count int
		var created time.Time
		if err := rows.Scan(&sn.Title, &sn.Message, &sn.Type, &count, &created); err == nil {
			sn.Created = created
			sn.Target = fmt.Sprintf("Sent to %d users", count)
			diff := now.Sub(created)
			if diff.Hours() < 1 {
				sn.TimeAgo = fmt.Sprintf("%d mins ago", int(diff.Minutes()))
			} else if diff.Hours() < 24 {
				sn.TimeAgo = fmt.Sprintf("%d hours ago", int(diff.Hours()))
			} else {
				sn.TimeAgo = fmt.Sprintf("%d days ago", int(diff.Hours()/24))
			}
			list = append(list, sn)
		}
	}

	// Mock list if fresh DB has zero values
	if len(list) == 0 {
		list = []SentNotification{
			{Title: "Deadline Reminder — Ladli Behna Yojana", Message: "Dear citizen, Ladli Behna Yojana deadline is approaching...", Type: "Deadline Reminder", Target: "Sent to 4,230 women users", TimeAgo: "2 hours ago"},
			{Title: "New Scheme Alert — PM Vishwakarma Yojana", Message: "New benefits under PM Vishwakarma are active...", Type: "New Scheme Alert", Target: "Sent to 18,430 users", TimeAgo: "1 day ago"},
			{Title: "Deadline Reminder — NSP Scholarship 2024", Message: "Verify your SC/ST documents for scholarship upload...", Type: "Deadline Reminder", Target: "Sent to 2,100 student users", TimeAgo: "3 days ago"},
		}
	}

	writeJSONResponse(w, http.StatusOK, list)
}

// AdminApplicationsHandler lists all submitted citizen applications
func AdminApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := `
		SELECT a.id, a.user_id, p.full_name, u.email, u.phone, 
		       a.scheme_id, s.title, s.government_level, 
		       a.status, a.applied_at, COALESCE(a.notes, ''), a.updated_at, COALESCE(p.aadhaar_encrypted, '')
		FROM user_applied_schemes a
		JOIN users u ON a.user_id = u.id
		JOIN user_profiles p ON u.id = p.user_id
		JOIN schemes s ON a.scheme_id = s.id
		ORDER BY a.applied_at DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	type AdminApplicationResponse struct {
		ID              int       `json:"id"`
		UserID          int       `json:"user_id"`
		FullName        string    `json:"full_name"`
		Email           string    `json:"email"`
		Phone           string    `json:"phone"`
		SchemeID        int       `json:"scheme_id"`
		SchemeTitle     string    `json:"scheme_title"`
		GovernmentLevel string    `json:"government_level"`
		Status          string    `json:"status"`
		AppliedAt       time.Time `json:"applied_at"`
		Notes           string    `json:"notes"`
		UpdatedAt       time.Time `json:"updated_at"`
		Aadhaar         string    `json:"aadhaar"`
	}

	var apps []AdminApplicationResponse = []AdminApplicationResponse{}
	for rows.Next() {
		var a AdminApplicationResponse
		var aadhaarEncrypted string
		err := rows.Scan(
			&a.ID, &a.UserID, &a.FullName, &a.Email, &a.Phone,
			&a.SchemeID, &a.SchemeTitle, &a.GovernmentLevel,
			&a.Status, &a.AppliedAt, &a.Notes, &a.UpdatedAt, &aadhaarEncrypted,
		)
		if err == nil {
			decrypted, _ := db.Decrypt(aadhaarEncrypted)
			a.Aadhaar = decrypted
			apps = append(apps, a)
		}
	}

	writeJSONResponse(w, http.StatusOK, apps)
}

// AdminApplicationStatusHandler approves or rejects an application
func AdminApplicationStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		ApplicationID int    `json:"application_id"`
		Status        string `json:"status"` // "approved" or "rejected"
		Notes         string `json:"notes"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ApplicationID == 0 || (req.Status != "approved" && req.Status != "rejected") {
		writeJSONError(w, http.StatusBadRequest, "Invalid application payload")
		return
	}

	// Verify application exists
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_applied_schemes WHERE id = $1)", req.ApplicationID).Scan(&exists)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	if !exists {
		writeJSONError(w, http.StatusNotFound, "Application not found")
		return
	}

	// Update status
	queryUpdate := `
		UPDATE user_applied_schemes 
		SET status = $1, notes = $2, updated_at = NOW() 
		WHERE id = $3`
	_, err = db.DB.Exec(queryUpdate, req.Status, req.Notes, req.ApplicationID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update application status: "+err.Error())
		return
	}

	// Fetch user_id, phone, full_name, and scheme title to insert a notification and print SMS log
	var userID int
	var schemeTitle, phone, fullName string
	err = db.DB.QueryRow(`
		SELECT a.user_id, s.title, u.phone, p.full_name
		FROM user_applied_schemes a 
		JOIN schemes s ON a.scheme_id = s.id 
		JOIN users u ON a.user_id = u.id
		JOIN user_profiles p ON u.id = p.user_id
		WHERE a.id = $1`, req.ApplicationID).Scan(&userID, &schemeTitle, &phone, &fullName)
	if err == nil {
		notifyTitle := "Application Update"
		notifyMsg := fmt.Sprintf("Your application for '%s' has been %s by the administrator.", schemeTitle, req.Status)
		if req.Notes != "" {
			notifyMsg += fmt.Sprintf(" Remarks: %s", req.Notes)
		}
		// Insert user notification
		_, _ = db.DB.Exec(`
			INSERT INTO notifications (user_id, title, message, type, is_read, created_at)
			VALUES ($1, $2, $3, 'Application Update', false, NOW())`,
			userID, notifyTitle, notifyMsg)

		// Simulated production SMS dispatch log
		log.Printf("[SMS GATEWAY] Sending SMS notification to phone +91 %s: \"Dear %s, your application for '%s' has been %s. Remarks: %s\"", 
			phone, fullName, schemeTitle, strings.ToUpper(req.Status), req.Notes)
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Application successfully %s!", req.Status),
	})
}

