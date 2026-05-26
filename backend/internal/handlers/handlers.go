package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/models"
)

// EnableCors helper to add CORS headers
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// GetSchemesHandler fetches schemes, supporting category filtering, search, and sorting
func GetSchemesHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
	EnableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
	EnableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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
	EnableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
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
	EnableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
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
		       p.education_level, p.is_disabled
		FROM users u
		JOIN user_profiles p ON u.id = p.user_id
		WHERE u.email = $1`

	err = db.DB.QueryRow(query, req.Email).Scan(
		&userID, &email, &phone, &passwordHash, &isAdmin,
		&profile.ID, &profile.FullName, &profile.DateOfBirth, &profile.Gender, &profile.State, &profile.District,
		&profile.CasteCategory, &profile.AnnualIncome, &profile.Occupation, &profile.EmployeeType,
		&profile.EducationLevel, &profile.IsDisabled,
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

	// Issuer JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secure-32-char-jwt-secret-key-majhigov"
	}

	expiryHoursStr := os.Getenv("JWT_EXPIRY_HOURS")
	expiryHours := 24
	if eh, err := strconv.Atoi(expiryHoursStr); err == nil {
		expiryHours = eh
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
	resp := models.AuthResponse{
		Success: true,
		Message: "Login successful!",
		Token:   tokenString,
		Profile: &profile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetUserProfileHandler extracts, decodes, and parses Bearer token to return user profile
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid authorization header", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secure-32-char-jwt-secret-key-majhigov"
	}

	// Parse Token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired session token", http.StatusUnauthorized)
		return
	}

	// Extract Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token format claims", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Missing user session identification in token", http.StatusUnauthorized)
		return
	}
	userID := int(userIDFloat)

	// Fetch Profile
	var profile models.UserProfile
	query := `
		SELECT id, user_id, full_name, date_of_birth, gender, state, district,
		       caste_category, annual_income, occupation, employee_type,
		       education_level, is_disabled
		FROM user_profiles
		WHERE user_id = $1`

	err = db.DB.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.FullName, &profile.DateOfBirth, &profile.Gender, &profile.State, &profile.District,
		&profile.CasteCategory, &profile.AnnualIncome, &profile.Occupation, &profile.EmployeeType,
		&profile.EducationLevel, &profile.IsDisabled,
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
