package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/middleware"
	"yojana-portal/backend/internal/models"
)

// GetJobsHandler handles GET /api/jobs for citizens, calculating match scores and concessions
func GetJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 1. Fetch all active jobs from database
	query := `
		SELECT id, title, title_hi, title_mr, organization, organization_hi, organization_mr,
		       description, description_hi, description_mr, education_qualification, documents_required,
		       min_age, max_age, last_date, apply_link, general_fee, obc_fee, sc_st_fee, women_fee,
		       is_active, clicks_count, created_at, updated_at
		FROM government_jobs
		WHERE is_active = true
		ORDER BY last_date ASC`

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

	// 2. Check if a user is logged in via Auth header (optional decoding)
	authHeader := r.Header.Get("Authorization")
	var profile *models.UserProfile
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		jwtSecret := middleware.GetJWTSecret()

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userIDFloat, ok := claims["user_id"].(float64); ok {
					userID := int(userIDFloat)
					// Retrieve citizen profile
					profile = getCitizenProfile(userID)
				}
			}
		}
	}

	// 3. Compute AI profile compatibility match scores and dynamic fee categories
	for i := range jobs {
		if profile != nil {
			// Compute Match Score
			score, reasons := computeJobMatchScore(profile, &jobs[i])
			jobs[i].MatchScore = score
			jobs[i].MatchReasons = reasons

			// Compute dynamic fee categories
			fee, label := calculateDynamicFee(profile, &jobs[i])
			jobs[i].CalculatedFee = fee
			jobs[i].FeeConcessionLabel = label
		} else {
			// Guest defaults
			jobs[i].MatchScore = 0
			jobs[i].MatchReasons = []string{"Log in to view compatibility match details."}
			jobs[i].CalculatedFee = jobs[i].GeneralFee
			jobs[i].FeeConcessionLabel = "Sign in to check fee concessions"
		}
	}

	writeJSONResponse(w, http.StatusOK, jobs)
}

// TrackJobClickHandler increments click analytics counters
func TrackJobClickHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse job ID from URL path (e.g. /api/jobs/:id/click)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeJSONError(w, http.StatusBadRequest, "Invalid request path")
		return
	}
	jobIDStr := parts[3]
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid job ID specification")
		return
	}

	// Increment Click counter in database
	_, err = db.DB.Exec("UPDATE government_jobs SET clicks_count = clicks_count + 1 WHERE id = $1", jobID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Redirection tracked"})
}

// Helper: retrieve citizen profile from database
func getCitizenProfile(userID int) *models.UserProfile {
	var p models.UserProfile
	query := `
		SELECT id, user_id, full_name, date_of_birth, gender, state, district,
		       caste_category, annual_income, occupation, employee_type,
		       education_level, is_disabled
		FROM user_profiles
		WHERE user_id = $1`

	err := db.DB.QueryRow(query, userID).Scan(
		&p.ID, &p.UserID, &p.FullName, &p.DateOfBirth, &p.Gender, &p.State, &p.District,
		&p.CasteCategory, &p.AnnualIncome, &p.Occupation, &p.EmployeeType,
		&p.EducationLevel, &p.IsDisabled,
	)
	if err != nil {
		return nil
	}
	return &p
}

// Helper: Calculate age from DOB string YYYY-MM-DD
func calculateAge(dobStr string) int {
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		// Fallback parse if timestamp format
		dob, err = time.Parse(time.RFC3339, dobStr)
		if err != nil {
			return 25 // default fallback
		}
	}
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

// Helper: Check education levels hierarchy
func checkEducationMatch(userEdu, jobEdu string) bool {
	eduRank := map[string]int{
		"10th Pass":     1,
		"12th Pass":     2,
		"Undergraduate": 3,
		"Graduate":      4,
		"Post Graduate": 5,
	}

	userRank, ok1 := eduRank[userEdu]
	jobRank, ok2 := eduRank[jobEdu]

	if !ok1 {
		userRank = 1
	}
	if !ok2 {
		return true // No requirement
	}

	return userRank >= jobRank
}

// Helper: AI profile matching algorithm
func computeJobMatchScore(profile *models.UserProfile, job *models.GovernmentJob) (int, []string) {
	score := 100
	var reasons []string

	// 1. Age match
	age := calculateAge(profile.DateOfBirth)
	if age < job.MinAge || age > job.MaxAge {
		score -= 30
		reasons = append(reasons, fmt.Sprintf("Your age (%d) is outside the permitted range of %d-%d years.", age, job.MinAge, job.MaxAge))
	} else {
		reasons = append(reasons, "Age matches job specification.")
	}

	// 2. Qualification match
	if !checkEducationMatch(profile.EducationLevel, job.EducationQualification) {
		score -= 40
		reasons = append(reasons, fmt.Sprintf("Minimum qualification required: '%s'. You have: '%s'.", job.EducationQualification, profile.EducationLevel))
	} else {
		reasons = append(reasons, "Educational qualification meets requirement.")
	}

	// 3. Documents readiness match (Check simulated availability)
	missingDocs := 0
	for _, doc := range job.DocumentsRequired {
		hasDoc := false
		if strings.EqualFold(doc, "Aadhaar Card") {
			hasDoc = true // seeded default
		} else if strings.Contains(strings.ToLower(doc), "degree") || strings.Contains(strings.ToLower(doc), "graduation") {
			if profile.EducationLevel == "Graduate" || profile.EducationLevel == "Post Graduate" {
				hasDoc = true
			}
		} else if strings.Contains(strings.ToLower(doc), "caste") {
			if profile.CasteCategory != "General" {
				hasDoc = true
			}
		} else if strings.Contains(strings.ToLower(doc), "domicile") {
			hasDoc = true // Assume available
		} else if strings.Contains(strings.ToLower(doc), "10th") || strings.Contains(strings.ToLower(doc), "12th") {
			hasDoc = true // Assume available
		}

		if !hasDoc {
			missingDocs++
			score -= 10
			reasons = append(reasons, fmt.Sprintf("Missing document: '%s'", doc))
		}
	}

	if missingDocs == 0 {
		reasons = append(reasons, "All required documents are ready in your profile.")
	}

	if score < 0 {
		score = 0
	}
	return score, reasons
}

// Helper: Calculate concessions based on demographic profile categories
func calculateDynamicFee(profile *models.UserProfile, job *models.GovernmentJob) (float64, string) {
	// Priority 1: Women concessions
	if strings.EqualFold(profile.Gender, "Female") && job.WomenFee < job.GeneralFee {
		if job.WomenFee == 0 {
			return 0, "Exempted (Women Concession)"
		}
		return job.WomenFee, fmt.Sprintf("₹%.0f (Women Concession Applied)", job.WomenFee)
	}

	// Priority 2: Caste concessions
	if (strings.EqualFold(profile.CasteCategory, "SC") || strings.EqualFold(profile.CasteCategory, "ST")) && job.ScStFee < job.GeneralFee {
		if job.ScStFee == 0 {
			return 0, "Exempted (SC/ST Concession)"
		}
		return job.ScStFee, fmt.Sprintf("₹%.0f (SC/ST Concession Applied)", job.ScStFee)
	}

	// Priority 3: OBC concessions
	if strings.EqualFold(profile.CasteCategory, "OBC") && job.ObcFee < job.GeneralFee {
		if job.ObcFee == 0 {
			return 0, "Exempted (OBC Concession)"
		}
		return job.ObcFee, fmt.Sprintf("₹%.0f (OBC Concession Applied)", job.ObcFee)
	}

	// Default
	return job.GeneralFee, "Standard General Category Fee"
}
