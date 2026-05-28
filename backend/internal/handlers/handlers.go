package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/middleware"
	"yojana-portal/backend/internal/models"
)

// GetSchemesHandler fetches schemes, supporting category filtering, search, and sorting
func GetSchemesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort_by") // "date_desc", "title_asc"

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT s.id, s.title, s.title_hi, s.title_mr, s.description, s.description_hi, s.description_mr,
		       s.category_id, c.name as category_name, c.name_hi as category_name_hi, c.name_mr as category_name_mr,
		       s.government_level, s.state, s.benefits, s.application_start_date, s.application_end_date,
		       s.official_website, s.apply_link, s.is_active, s.created_at, s.updated_at
		FROM schemes s
		JOIN scheme_categories c ON s.category_id = c.id
		WHERE s.is_active = true
	`)
	
	var args []interface{}
	argCount := 1

	if category != "" && category != "All" {
		queryBuilder.WriteString(" AND (c.name = $")
		queryBuilder.WriteString(strconv.Itoa(argCount))
		queryBuilder.WriteString(" OR c.name_hi = $")
		queryBuilder.WriteString(strconv.Itoa(argCount))
		queryBuilder.WriteString(" OR c.name_mr = $")
		queryBuilder.WriteString(strconv.Itoa(argCount))
		queryBuilder.WriteString(")")
		args = append(args, category)
		argCount++
	}

	if search != "" {
		queryBuilder.WriteString(" AND (s.title ILIKE $")
		queryBuilder.WriteString(strconv.Itoa(argCount))
		queryBuilder.WriteString(" OR s.description ILIKE $")
		queryBuilder.WriteString(strconv.Itoa(argCount))
		queryBuilder.WriteString(" OR s.title_hi ILIKE $")
		queryBuilder.WriteString(strconv.Itoa(argCount))
		queryBuilder.WriteString(" OR s.title_mr ILIKE $")
		queryBuilder.WriteString(strconv.Itoa(argCount))
		queryBuilder.WriteString(")")
		args = append(args, "%"+search+"%")
		argCount++
	}

	if sortBy == "title_asc" {
		queryBuilder.WriteString(" ORDER BY s.title ASC")
	} else {
		queryBuilder.WriteString(" ORDER BY s.created_at DESC")
	}

	rows, err := db.DB.Query(queryBuilder.String(), args...)
	if err != nil {
		http.Error(w, "Failed to query database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var schemes []models.Scheme = []models.Scheme{}
	for rows.Next() {
		var s models.Scheme
		err := rows.Scan(
			&s.ID, &s.Title, &s.TitleHi, &s.TitleMr, &s.Description, &s.DescriptionHi, &s.DescriptionMr,
			&s.CategoryID, &s.CategoryName, &s.CategoryNameHi, &s.CategoryNameMr,
			&s.GovernmentLevel, &s.State, &s.Benefits, &s.ApplicationStartDate, &s.ApplicationEndDate,
			&s.OfficialWebsite, &s.ApplyLink, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to scan scheme row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		schemes = append(schemes, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schemes)
}

// GetSchemeDetailsHandler retrieves details for a single scheme along with linked documents, FAQs, and eligibility
func GetSchemeDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid scheme ID in path", http.StatusBadRequest)
		return
	}
	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid scheme ID format", http.StatusBadRequest)
		return
	}

	// Fetch Scheme main record
	var s models.Scheme
	queryScheme := `
		SELECT s.id, s.title, s.title_hi, s.title_mr, s.description, s.description_hi, s.description_mr,
		       s.category_id, c.name as category_name, c.name_hi as category_name_hi, c.name_mr as category_name_mr,
		       s.government_level, s.state, s.benefits, s.application_start_date, s.application_end_date,
		       s.official_website, s.apply_link, s.is_active, s.created_at, s.updated_at
		FROM schemes s
		JOIN scheme_categories c ON s.category_id = c.id
		WHERE s.id = $1 AND s.is_active = true`
	
	err = db.DB.QueryRow(queryScheme, id).Scan(
		&s.ID, &s.Title, &s.TitleHi, &s.TitleMr, &s.Description, &s.DescriptionHi, &s.DescriptionMr,
		&s.CategoryID, &s.CategoryName, &s.CategoryNameHi, &s.CategoryNameMr,
		&s.GovernmentLevel, &s.State, &s.Benefits, &s.ApplicationStartDate, &s.ApplicationEndDate,
		&s.OfficialWebsite, &s.ApplyLink, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Scheme not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch Documents
	rowsDocs, err := db.DB.Query(`
		SELECT id, scheme_id, document_name, document_name_hi, document_name_mr, is_mandatory, created_at 
		FROM scheme_documents 
		WHERE scheme_id = $1`, id)
	if err != nil {
		http.Error(w, "Failed to query documents: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rowsDocs.Close()

	s.Documents = []models.SchemeDocument{}
	for rowsDocs.Next() {
		var doc models.SchemeDocument
		if err := rowsDocs.Scan(&doc.ID, &doc.SchemeID, &doc.DocumentName, &doc.DocumentNameHi, &doc.DocumentNameMr, &doc.IsMandatory, &doc.CreatedAt); err != nil {
			http.Error(w, "Failed to scan document: "+err.Error(), http.StatusInternalServerError)
			return
		}
		s.Documents = append(s.Documents, doc)
	}

	// Fetch FAQs
	rowsFAQs, err := db.DB.Query(`
		SELECT id, scheme_id, question, answer, question_hi, answer_hi, question_mr, answer_mr, created_at 
		FROM scheme_faqs 
		WHERE scheme_id = $1`, id)
	if err != nil {
		http.Error(w, "Failed to query FAQs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rowsFAQs.Close()

	s.FAQs = []models.SchemeFAQ{}
	for rowsFAQs.Next() {
		var faq models.SchemeFAQ
		if err := rowsFAQs.Scan(&faq.ID, &faq.SchemeID, &faq.Question, &faq.Answer, &faq.QuestionHi, &faq.AnswerHi, &faq.QuestionMr, &faq.AnswerMr, &faq.CreatedAt); err != nil {
			http.Error(w, "Failed to scan FAQ: "+err.Error(), http.StatusInternalServerError)
			return
		}
		s.FAQs = append(s.FAQs, faq)
	}

	// Fetch Eligibility Criteria
	var e models.EligibilityCriteria
	queryCriteria := `
		SELECT id, scheme_id, min_age, max_age, gender, caste_categories, min_income, max_income, states, occupations, employee_types, education_levels, disability_required, created_at 
		FROM eligibility_criteria 
		WHERE scheme_id = $1`
	err = db.DB.QueryRow(queryCriteria, id).Scan(
		&e.ID, &e.SchemeID, &e.MinAge, &e.MaxAge, &e.Gender, &e.CasteCategories,
		&e.MinIncome, &e.MaxIncome, &e.States, &e.Occupations, &e.EmployeeTypes,
		&e.EducationLevels, &e.DisabilityRequired, &e.CreatedAt,
	)
	if err == sql.ErrNoRows {
		s.Eligibility = nil
	} else if err != nil {
		http.Error(w, "Failed to query eligibility: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		s.Eligibility = &e
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

// CheckEligibilityHandler evaluates user profile properties against relational eligibility rules using string array checking
func CheckEligibilityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.UserProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch schemes along with eligibility criteria in a single JOIN query
	query := `
		SELECT s.id, s.title, s.title_hi, s.title_mr, s.description, s.description_hi, s.description_mr,
		       s.category_id, c.name as category_name, c.name_hi as category_name_hi, c.name_mr as category_name_mr,
		       s.government_level, s.state, s.benefits, s.application_start_date, s.application_end_date,
		       s.official_website, s.apply_link, s.is_active, s.created_at, s.updated_at,
		       e.id as elig_id, e.min_age, e.max_age, e.gender, e.caste_categories, e.min_income, e.max_income, 
		       e.states, e.occupations, e.employee_types, e.education_levels, e.disability_required
		FROM schemes s
		JOIN scheme_categories c ON s.category_id = c.id
		LEFT JOIN eligibility_criteria e ON s.id = e.scheme_id
		WHERE s.is_active = true`

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch schemes for checking: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var response models.EligibilityResponse
	response.Eligible = []models.EligibilityStatus{}
	response.NotEligible = []models.EligibilityStatus{}

	for rows.Next() {
		var s models.Scheme
		var eligNullID sql.NullInt64
		var e models.EligibilityCriteria

		err := rows.Scan(
			&s.ID, &s.Title, &s.TitleHi, &s.TitleMr, &s.Description, &s.DescriptionHi, &s.DescriptionMr,
			&s.CategoryID, &s.CategoryName, &s.CategoryNameHi, &s.CategoryNameMr,
			&s.GovernmentLevel, &s.State, &s.Benefits, &s.ApplicationStartDate, &s.ApplicationEndDate,
			&s.OfficialWebsite, &s.ApplyLink, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
			&eligNullID, &e.MinAge, &e.MaxAge, &e.Gender, &e.CasteCategories, &e.MinIncome, &e.MaxIncome,
			&e.States, &e.Occupations, &e.EmployeeTypes, &e.EducationLevels, &e.DisabilityRequired,
		)
		if err != nil {
			http.Error(w, "Scanning error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if eligNullID.Valid {
			e.ID = int(eligNullID.Int64)
			e.SchemeID = s.ID
			s.Eligibility = &e
		} else {
			s.Eligibility = nil
		}

		var isEligible bool
		var reasons []string

		if s.Eligibility == nil {
			isEligible = true
			reasons = []string{"This scheme has no standard restrictions. Everyone is eligible!"}
		} else {
			isEligible, reasons = evaluateEligibilityCriteria(req, *s.Eligibility)
		}

		status := models.EligibilityStatus{
			Scheme:     s,
			IsEligible: isEligible,
			Reasons:    reasons,
		}

		if isEligible {
			response.Eligible = append(response.Eligible, status)
		} else {
			response.NotEligible = append(response.NotEligible, status)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func evaluateEligibilityCriteria(profile models.UserProfileRequest, rules models.EligibilityCriteria) (bool, []string) {
	var reasons []string
	eligible := true

	// 1. Age Check
	if profile.Age < rules.MinAge || profile.Age > rules.MaxAge {
		eligible = false
		reasons = append(reasons, "Your age ("+strconv.Itoa(profile.Age)+") is not in the required range of "+strconv.Itoa(rules.MinAge)+" to "+strconv.Itoa(rules.MaxAge)+" years.")
	} else {
		reasons = append(reasons, "Age meets the requirement of "+strconv.Itoa(rules.MinAge)+"-"+strconv.Itoa(rules.MaxAge)+" years.")
	}

	// 2. Income Check
	if profile.AnnualIncome < rules.MinIncome || profile.AnnualIncome > rules.MaxIncome {
		eligible = false
		reasons = append(reasons, "Your family income (₹"+strconv.FormatFloat(profile.AnnualIncome, 'f', 0, 64)+") is outside the eligible bracket of ₹"+strconv.FormatFloat(rules.MinIncome, 'f', 0, 64)+" to ₹"+strconv.FormatFloat(rules.MaxIncome, 'f', 0, 64)+".")
	} else {
		reasons = append(reasons, "Income falls within the allowed eligibility range.")
	}

	// 3. Gender Check
	if !strings.EqualFold(rules.Gender, "All") && !strings.EqualFold(rules.Gender, profile.Gender) {
		eligible = false
		reasons = append(reasons, "Gender ("+profile.Gender+") is not eligible. Required: "+rules.Gender+".")
	} else {
		reasons = append(reasons, "Gender meets the requirements.")
	}

	// 4. State Check (Array)
	if len(rules.States) > 0 && !containsString(rules.States, profile.State) {
		eligible = false
		reasons = append(reasons, "Scheme is not active in "+profile.State+". Eligible: "+strings.Join(rules.States, ", "))
	} else {
		reasons = append(reasons, "State is eligible.")
	}

	// 5. Caste Check (Array)
	if len(rules.CasteCategories) > 0 && !containsString(rules.CasteCategories, profile.Caste) {
		eligible = false
		reasons = append(reasons, "Your caste category ("+profile.Caste+") is not eligible. Required: "+strings.Join(rules.CasteCategories, ", "))
	} else {
		reasons = append(reasons, "Caste category is eligible.")
	}

	// 6. Occupation Check (Array)
	if len(rules.Occupations) > 0 && !containsString(rules.Occupations, profile.Occupation) {
		eligible = false
		reasons = append(reasons, "Your occupation ("+profile.Occupation+") is not eligible. Required: "+strings.Join(rules.Occupations, ", "))
	} else {
		reasons = append(reasons, "Occupation is eligible.")
	}

	// 7. Employee Type Check (Array)
	if len(rules.EmployeeTypes) > 0 && !containsString(rules.EmployeeTypes, profile.EmployeeType) {
		eligible = false
		reasons = append(reasons, "Employment category ("+profile.EmployeeType+") is not eligible. Required: "+strings.Join(rules.EmployeeTypes, ", "))
	} else {
		reasons = append(reasons, "Employment category matches.")
	}

	// 8. Education Level Check (Array)
	if len(rules.EducationLevels) > 0 && !containsString(rules.EducationLevels, profile.EducationLevel) {
		eligible = false
		reasons = append(reasons, "Education level ("+profile.EducationLevel+") is not eligible. Required: "+strings.Join(rules.EducationLevels, ", "))
	} else {
		reasons = append(reasons, "Education level is eligible.")
	}

	// 9. Disability Check
	if rules.DisabilityRequired && !profile.IsDisabled {
		eligible = false
		reasons = append(reasons, "This scheme is specifically reserved for differently-abled citizens.")
	} else if rules.DisabilityRequired && profile.IsDisabled {
		reasons = append(reasons, "Disability status satisfies the requirement.")
	}

	return eligible, reasons
}

func containsString(list []string, value string) bool {
	if len(list) == 0 {
		return true
	}
	for _, item := range list {
		if strings.EqualFold(item, "All") || strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}

// RegisterHandler inserts new user credentials and profile information atomically
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Phone == "" || req.Password == "" || req.FullName == "" {
		http.Error(w, "Missing required registration parameters", http.StatusBadRequest)
		return
	}

	// Check email/phone uniqueness
	var exists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR phone = $2)"
	err = db.DB.QueryRow(queryCheck, req.Email, req.Phone).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email or phone number is already registered", http.StatusBadRequest)
		return
	}

	// Hashing
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Begin Transaction
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction initialization failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert User
	var userID int
	queryUser := `
		INSERT INTO users (email, phone, password_hash, is_verified, is_admin)
		VALUES ($1, $2, $3, false, false) RETURNING id`
	err = tx.QueryRow(queryUser, req.Email, req.Phone, string(hash)).Scan(&userID)
	if err != nil {
		http.Error(w, "Failed to insert user account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert User Profile
	queryProfile := `
		INSERT INTO user_profiles (
			user_id, full_name, date_of_birth, gender, state, district,
			caste_category, annual_income, occupation, employee_type,
			education_level, is_disabled
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err = tx.Exec(
		queryProfile, userID, req.FullName, req.DateOfBirth, req.Gender, req.State, req.District,
		req.CasteCategory, req.AnnualIncome, req.Occupation, req.EmployeeType,
		req.EducationLevel, req.IsDisabled,
	)
	if err != nil {
		http.Error(w, "Failed to insert user profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit Transaction
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit user registration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User account registered successfully!",
	})
}

// LoginHandler matches raw password with hash and issues signed JWT bearer tokens
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid login format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Join User and Profile to fetch all in one query
	var userID int
	var email, phone, passwordHash string
	var isAdmin bool
	var profile models.UserProfile

	query := `
		SELECT u.id, u.email, u.phone, u.password_hash, u.is_admin,
		       p.id, p.full_name, p.date_of_birth, p.gender, p.state, p.district,
		       p.caste_category, p.annual_income, p.occupation, p.employee_type,
		       p.education_level, p.is_disabled, COALESCE(p.avatar_url, '')
		FROM users u
		JOIN user_profiles p ON u.id = p.user_id
		WHERE u.email = $1`

	err = db.DB.QueryRow(query, req.Email).Scan(
		&userID, &email, &phone, &passwordHash, &isAdmin,
		&profile.ID, &profile.FullName, &profile.DateOfBirth, &profile.Gender, &profile.State, &profile.District,
		&profile.CasteCategory, &profile.AnnualIncome, &profile.Occupation, &profile.EmployeeType,
		&profile.EducationLevel, &profile.IsDisabled, &profile.AvatarURL,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database lookup error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify Hash
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Issue JWT using middleware's centralized secret
	secret := middleware.GetJWTSecret()

	expiryHours := 24
	if ehStr := os.Getenv("JWT_EXPIRY_HOURS"); ehStr != "" {
		if eh, err := strconv.Atoi(ehStr); err == nil {
			expiryHours = eh
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Failed to sign authentication token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	profile.UserID = userID
	profile.Email = email
	profile.Phone = phone
	resp := models.AuthResponse{
		Success: true,
		Message: "Login successful!",
		Token:   tokenString,
		Profile: &profile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UserProfileHandler routes GET and PUT requests for user profile
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetUserProfileHandler(w, r)
	case "PUT":
		UpdateUserProfileHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUserProfileHandler returns the authenticated user's profile
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch Profile joined with credentials
	var profile models.UserProfile
	query := `
		SELECT p.id, p.user_id, p.full_name, p.date_of_birth, p.gender, p.state, p.district,
		       p.caste_category, p.annual_income, p.occupation, p.employee_type,
		       p.education_level, p.is_disabled, COALESCE(p.avatar_url, ''), u.email, u.phone
		FROM user_profiles p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1`

	err = db.DB.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.FullName, &profile.DateOfBirth, &profile.Gender, &profile.State, &profile.District,
		&profile.CasteCategory, &profile.AnnualIncome, &profile.Occupation, &profile.EmployeeType,
		&profile.EducationLevel, &profile.IsDisabled, &profile.AvatarURL, &profile.Email, &profile.Phone,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "User profile not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database retrieval error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"profile": profile,
	})
}

// UpdateUserProfileHandler updates the authenticated user's demographic details
func UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		FullName       string  `json:"full_name"`
		DateOfBirth    string  `json:"date_of_birth"`
		Gender         string  `json:"gender"`
		State          string  `json:"state"`
		District       string  `json:"district"`
		CasteCategory  string  `json:"caste_category"`
		AnnualIncome   float64 `json:"annual_income"`
		Occupation     string  `json:"occupation"`
		EmployeeType   string  `json:"employee_type"`
		EducationLevel string  `json:"education_level"`
		IsDisabled     bool    `json:"is_disabled"`
		AvatarURL      string  `json:"avatar_url"`
		Email          string  `json:"email"`
		Phone          string  `json:"phone"`
		Password       string  `json:"password"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.FullName == "" {
		http.Error(w, "Full name is required", http.StatusBadRequest)
		return
	}

	// Begin SQL transaction
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Verify and update User email, phone, and password if provided
	if req.Email != "" {
		var exists bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)", req.Email, userID).Scan(&exists)
		if err != nil {
			http.Error(w, "Failed to check email uniqueness: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Email is already registered by another user", http.StatusBadRequest)
			return
		}

		_, err = tx.Exec("UPDATE users SET email = $1 WHERE id = $2", req.Email, userID)
		if err != nil {
			http.Error(w, "Failed to update email: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if req.Phone != "" {
		var exists bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1 AND id != $2)", req.Phone, userID).Scan(&exists)
		if err != nil {
			http.Error(w, "Failed to check phone uniqueness: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Phone number is already registered by another user", http.StatusBadRequest)
			return
		}

		_, err = tx.Exec("UPDATE users SET phone = $1 WHERE id = $2", req.Phone, userID)
		if err != nil {
			http.Error(w, "Failed to update phone number: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", string(hash), userID)
		if err != nil {
			http.Error(w, "Failed to update password: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 2. Update user profile details
	queryUpdateProfile := `
		UPDATE user_profiles SET
			full_name=$1, date_of_birth=$2, gender=$3, state=$4, district=$5,
			caste_category=$6, annual_income=$7, occupation=$8, employee_type=$9,
			education_level=$10, is_disabled=$11, avatar_url=$12, updated_at=NOW()
		WHERE user_id=$13`

	_, err = tx.Exec(queryUpdateProfile,
		req.FullName, req.DateOfBirth, req.Gender, req.State, req.District,
		req.CasteCategory, req.AnnualIncome, req.Occupation, req.EmployeeType,
		req.EducationLevel, req.IsDisabled, req.AvatarURL, userID,
	)
	if err != nil {
		http.Error(w, "Failed to update profile details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit database changes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch updated profile
	var profile models.UserProfile
	queryFetch := `
		SELECT p.id, p.user_id, p.full_name, p.date_of_birth, p.gender, p.state, p.district,
		       p.caste_category, p.annual_income, p.occupation, p.employee_type,
		       p.education_level, p.is_disabled, COALESCE(p.avatar_url, ''), u.email, u.phone
		FROM user_profiles p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1`

	err = db.DB.QueryRow(queryFetch, userID).Scan(
		&profile.ID, &profile.UserID, &profile.FullName, &profile.DateOfBirth, &profile.Gender, &profile.State, &profile.District,
		&profile.CasteCategory, &profile.AnnualIncome, &profile.Occupation, &profile.EmployeeType,
		&profile.EducationLevel, &profile.IsDisabled, &profile.AvatarURL, &profile.Email, &profile.Phone,
	)
	if err != nil {
		http.Error(w, "Failed to fetch updated profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
		"profile": profile,
	})
}

// ToggleSavedSchemeHandler toggles the bookmark status of a scheme in PostgreSQL
func ToggleSavedSchemeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req struct {
		SchemeID int `json:"scheme_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.SchemeID == 0 {
		http.Error(w, "Invalid scheme ID payload", http.StatusBadRequest)
		return
	}

	// Check if already saved
	var exists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM user_saved_schemes WHERE user_id = $1 AND scheme_id = $2)"
	err = db.DB.QueryRow(queryCheck, userID, req.SchemeID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var saved bool
	if exists {
		// Delete bookmark
		queryDel := "DELETE FROM user_saved_schemes WHERE user_id = $1 AND scheme_id = $2"
		_, err = db.DB.Exec(queryDel, userID, req.SchemeID)
		if err != nil {
			http.Error(w, "Failed to remove bookmark: "+err.Error(), http.StatusInternalServerError)
			return
		}
		saved = false
	} else {
		// Insert bookmark
		queryIns := "INSERT INTO user_saved_schemes (user_id, scheme_id, saved_at) VALUES ($1, $2, NOW())"
		_, err = db.DB.Exec(queryIns, userID, req.SchemeID)
		if err != nil {
			http.Error(w, "Failed to save bookmark: "+err.Error(), http.StatusInternalServerError)
			return
		}
		saved = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"saved":   saved,
		"message": "Bookmark synced successfully!",
	})
}

// GetSavedSchemesHandler retrieves all scheme IDs saved by the authenticated user
func GetSavedSchemesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query("SELECT scheme_id FROM user_saved_schemes WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to query database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var schemeIDs []int = []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			http.Error(w, "Error scanning database: "+err.Error(), http.StatusInternalServerError)
			return
		}
		schemeIDs = append(schemeIDs, id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schemeIDs)
}

// ApplySchemeHandler submits a new application to the database
func ApplySchemeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req struct {
		SchemeID int    `json:"scheme_id"`
		Notes    string `json:"notes"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.SchemeID == 0 {
		http.Error(w, "Invalid application payload", http.StatusBadRequest)
		return
	}

	// Prevent double application if a pending one already exists
	var exists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM user_applied_schemes WHERE user_id = $1 AND scheme_id = $2 AND status = 'pending')"
	err = db.DB.QueryRow(queryCheck, userID, req.SchemeID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database verification error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "You already have a pending application under review for this scheme.", http.StatusBadRequest)
		return
	}

	// Insert Application
	queryApply := `
		INSERT INTO user_applied_schemes (user_id, scheme_id, status, applied_at, notes, updated_at)
		VALUES ($1, $2, 'pending', NOW(), $3, NOW())`
	_, err = db.DB.Exec(queryApply, userID, req.SchemeID, req.Notes)
	if err != nil {
		http.Error(w, "Failed to submit application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Your application has been submitted successfully!",
	})
}

// GetUserApplicationsHandler retrieves all applications along with scheme names and statuses
func GetUserApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	query := `
		SELECT a.id, a.scheme_id, s.title, s.title_hi, s.title_mr, s.government_level, 
		       a.status, a.applied_at, a.notes, s.apply_link, s.official_website
		FROM user_applied_schemes a
		JOIN schemes s ON a.scheme_id = s.id
		WHERE a.user_id = $1
		ORDER BY a.applied_at DESC`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to query applications: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ApplicationResponse struct {
		ID              int       `json:"id"`
		SchemeID        int       `json:"scheme_id"`
		Title           string    `json:"title"`
		TitleHi         string    `json:"title_hi"`
		TitleMr         string    `json:"title_mr"`
		GovernmentLevel string    `json:"government_level"`
		Status          string    `json:"status"`
		AppliedAt       time.Time `json:"applied_at"`
		Notes           string    `json:"notes"`
		ApplyLink       string    `json:"apply_link"`
		OfficialWebsite string    `json:"official_website"`
	}

	var apps []ApplicationResponse = []ApplicationResponse{}
	for rows.Next() {
		var a ApplicationResponse
		err := rows.Scan(&a.ID, &a.SchemeID, &a.Title, &a.TitleHi, &a.TitleMr, &a.GovernmentLevel, &a.Status, &a.AppliedAt, &a.Notes, &a.ApplyLink, &a.OfficialWebsite)
		if err != nil {
			http.Error(w, "Error scanning database rows: "+err.Error(), http.StatusInternalServerError)
			return
		}
		apps = append(apps, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

// SavedSchemesHandler routes GET and POST requests for user bookmarks
func SavedSchemesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetSavedSchemesHandler(w, r)
	} else if r.Method == "POST" {
		ToggleSavedSchemeHandler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// TranslateHandler handles client GET translation proxy requests
func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.URL.Query().Get("q")
	target := r.URL.Query().Get("target")

	if text == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	if target == "" {
		target = "hi" // Default target Hindi
	}

	translated, err := translateTextViaGoogle(text, target)
	if err != nil {
		http.Error(w, "Translation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"translatedText": translated,
	})
}

func translateTextViaGoogle(text string, target string) (string, error) {
	apiURL := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=%s&dt=t&q=%s",
		target, url.QueryEscape(text))

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google translate api returned status %d", resp.StatusCode)
	}

	var result []interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result) == 0 || result[0] == nil {
		return "", fmt.Errorf("invalid translation array structure")
	}

	parts, ok := result[0].([]interface{})
	if !ok {
		return "", fmt.Errorf("invalid translation parts array")
	}

	var builder strings.Builder
	for _, p := range parts {
		inner, ok := p.([]interface{})
		if ok && len(inner) > 0 {
			translatedStr, ok := inner[0].(string)
			if ok {
				builder.WriteString(translatedStr)
			}
		}
	}

	return builder.String(), nil
}
