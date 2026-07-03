package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/middleware"
	"yojana-portal/backend/internal/models"
)

// Recommendation holds a matched scheme or job with AI-generated explanation
type Recommendation struct {
	Type        string          `json:"type"` // "scheme" or "job"
	Scheme      *models.Scheme  `json:"scheme,omitempty"`
	Job         *models.GovtJob `json:"job,omitempty"`
	MatchScore  int             `json:"match_score"`
	Explanation string          `json:"explanation"`
}

// Simple in-memory recommendation cache (per userID, TTL = 1 hour)
var (
	recCache   = make(map[int]*recCacheEntry)
	recCacheMu sync.RWMutex
)

type recCacheEntry struct {
	data      []Recommendation
	fetchedAt time.Time
}

const recCacheTTL = 1 * time.Hour

// GetRecommendationsHandler returns AI-matched schemes and jobs for the logged-in user.
// Results are cached per-user for 1 hour to avoid repeated Gemini API calls.
func GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Serve from cache if fresh
	recCacheMu.RLock()
	entry, cached := recCache[userID]
	recCacheMu.RUnlock()
	if cached && time.Since(entry.fetchedAt) < recCacheTTL {
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{
			"recommendations": entry.data,
			"cached":          true,
		})
		return
	}

	// Fetch user profile
	var profile models.UserProfile
	profileQuery := `
		SELECT p.full_name, p.date_of_birth, p.gender, p.state, p.caste_category,
		       p.annual_income, p.occupation, p.employee_type, p.education_level, p.is_disabled
		FROM user_profiles p WHERE p.user_id = $1`
	err = db.DB.QueryRow(profileQuery, userID).Scan(
		&profile.FullName, &profile.DateOfBirth, &profile.Gender, &profile.State,
		&profile.CasteCategory, &profile.AnnualIncome, &profile.Occupation,
		&profile.EmployeeType, &profile.EducationLevel, &profile.IsDisabled,
	)
	if err == sql.ErrNoRows {
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{"recommendations": []interface{}{}})
		return
	} else if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to load user profile")
		return
	}

	userAge := calculateAge(profile.DateOfBirth)

	// Build UserProfileRequest for eligibility matching
	eligReq := models.UserProfileRequest{
		Age:            userAge,
		Gender:         profile.Gender,
		State:          profile.State,
		Caste:          profile.CasteCategory,
		AnnualIncome:   profile.AnnualIncome,
		Occupation:     profile.Occupation,
		EmployeeType:   profile.EmployeeType,
		EducationLevel: profile.EducationLevel,
		IsDisabled:     profile.IsDisabled,
	}

	// Fetch active schemes with eligibility criteria
	schemeQuery := `
		SELECT s.id, s.title, s.title_hi, s.title_mr, s.description, s.description_hi, s.description_mr,
		       s.category_id, c.name, c.name_hi, c.name_mr,
		       s.government_level, s.state, s.benefits, s.application_start_date, s.application_end_date,
		       s.official_website, s.apply_link, s.is_active, s.created_at, s.updated_at,
		       e.id, e.min_age, e.max_age, e.gender, e.caste_categories, e.min_income,
		       e.max_income, e.states, e.occupations, e.employee_types, e.education_levels, e.disability_required
		FROM schemes s
		JOIN scheme_categories c ON s.category_id = c.id
		LEFT JOIN eligibility_criteria e ON s.id = e.scheme_id
		WHERE s.is_active = true
		LIMIT 100`

	schemeRows, err := db.DB.Query(schemeQuery)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch schemes")
		return
	}
	defer schemeRows.Close()

	type scoredScheme struct {
		scheme models.Scheme
		score  int
	}
	var scoredSchemes []scoredScheme

	for schemeRows.Next() {
		var s models.Scheme
		var eligNullID sql.NullInt64
		var e models.EligibilityCriteria
		err := schemeRows.Scan(
			&s.ID, &s.Title, &s.TitleHi, &s.TitleMr, &s.Description, &s.DescriptionHi, &s.DescriptionMr,
			&s.CategoryID, &s.CategoryName, &s.CategoryNameHi, &s.CategoryNameMr,
			&s.GovernmentLevel, &s.State, &s.Benefits, &s.ApplicationStartDate, &s.ApplicationEndDate,
			&s.OfficialWebsite, &s.ApplyLink, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
			&eligNullID, &e.MinAge, &e.MaxAge, &e.Gender, &e.CasteCategories, &e.MinIncome,
			&e.MaxIncome, &e.States, &e.Occupations, &e.EmployeeTypes, &e.EducationLevels, &e.DisabilityRequired,
		)
		if err != nil {
			continue
		}
		if eligNullID.Valid {
			e.ID = int(eligNullID.Int64)
			s.Eligibility = &e
		}

		var score int
		if s.Eligibility != nil {
			isElig, _ := evaluateEligibilityCriteria(eligReq, *s.Eligibility)
			if isElig {
				score = countMatchingCriteria(eligReq, *s.Eligibility)
			}
		} else {
			score = 5 // No restrictions = broad match
		}
		if score > 0 {
			scoredSchemes = append(scoredSchemes, scoredScheme{scheme: s, score: score})
		}
	}

	// Sort by score descending, take top 5
	sort.Slice(scoredSchemes, func(i, j int) bool {
		return scoredSchemes[i].score > scoredSchemes[j].score
	})
	if len(scoredSchemes) > 5 {
		scoredSchemes = scoredSchemes[:5]
	}

	// Fetch matching jobs (education-based match)
	jobQuery := `
		SELECT id, title, title_hi, title_mr, organization, department, vacancies,
		       education_qualification, experience_required, required_documents,
		       application_start_date, application_end_date, official_website, apply_link,
		       application_fee, is_active, created_at, updated_at
		FROM govt_jobs WHERE is_active = true LIMIT 50`

	var matchedJobs []models.GovtJob
	jobRows, err := db.DB.Query(jobQuery)
	if err == nil {
		defer jobRows.Close()
		for jobRows.Next() {
			var j models.GovtJob
			var startDate, endDate time.Time
			err := jobRows.Scan(
				&j.ID, &j.Title, &j.TitleHi, &j.TitleMr, &j.Organization, &j.Department, &j.Vacancies,
				&j.EducationQualification, &j.ExperienceRequired, pq.Array(&j.RequiredDocuments),
				&startDate, &endDate, &j.OfficialWebsite, &j.ApplyLink,
				&j.ApplicationFee, &j.IsActive, &j.CreatedAt, &j.UpdatedAt,
			)
			if err != nil {
				continue
			}
			j.ApplicationStartDate = startDate.Format("2006-01-02")
			j.ApplicationEndDate = endDate.Format("2006-01-02")
			if educationSufficient(profile.EducationLevel, j.EducationQualification) {
				matchedJobs = append(matchedJobs, j)
			}
		}
	}
	if len(matchedJobs) > 5 {
		matchedJobs = matchedJobs[:5]
	}

	// Build recommendations list with Gemini explanations
	var recommendations []Recommendation
	for _, ss := range scoredSchemes {
		schemeCopy := ss.scheme
		desc := schemeCopy.Description
		if len(desc) > 200 {
			desc = desc[:200]
		}
		rec := Recommendation{
			Type:        "scheme",
			Scheme:      &schemeCopy,
			MatchScore:  ss.score,
			Explanation: generateGeminiExplanation(profile, schemeCopy.Title, desc, schemeCopy.Benefits, "scheme"),
		}
		recommendations = append(recommendations, rec)
	}
	for i := range matchedJobs {
		jobCopy := matchedJobs[i]
		rec := Recommendation{
			Type:        "job",
			Job:         &jobCopy,
			MatchScore:  7,
			Explanation: generateGeminiExplanation(profile, jobCopy.Title, "Government job opening", jobCopy.Organization+" - "+jobCopy.Department, "job"),
		}
		recommendations = append(recommendations, rec)
	}

	if recommendations == nil {
		recommendations = []Recommendation{}
	}

	// Store in cache
	recCacheMu.Lock()
	recCache[userID] = &recCacheEntry{data: recommendations, fetchedAt: time.Now()}
	recCacheMu.Unlock()

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"recommendations": recommendations,
		"cached":          false,
	})
}

// countMatchingCriteria returns the number of eligibility criteria the user meets
func countMatchingCriteria(profile models.UserProfileRequest, rules models.EligibilityCriteria) int {
	score := 0
	if profile.Age >= rules.MinAge && profile.Age <= rules.MaxAge {
		score++
	}
	if profile.AnnualIncome >= rules.MinIncome && profile.AnnualIncome <= rules.MaxIncome {
		score++
	}
	if strings.EqualFold(rules.Gender, "all") || strings.EqualFold(rules.Gender, profile.Gender) {
		score++
	}
	if len(rules.States) == 0 || containsString(rules.States, profile.State) {
		score++
	}
	if len(rules.CasteCategories) == 0 || containsString(rules.CasteCategories, profile.Caste) {
		score++
	}
	if len(rules.Occupations) == 0 || containsString(rules.Occupations, profile.Occupation) {
		score++
	}
	if len(rules.EmployeeTypes) == 0 || containsString(rules.EmployeeTypes, profile.EmployeeType) {
		score++
	}
	if len(rules.EducationLevels) == 0 || containsString(rules.EducationLevels, profile.EducationLevel) {
		score++
	}
	if !rules.DisabilityRequired || profile.IsDisabled {
		score++
	}
	return score
}

// calculateAge computes age in years from a YYYY-MM-DD date string
func calculateAge(dobStr string) int {
	if len(dobStr) < 10 {
		return 30
	}
	var year, month, day int
	fmt.Sscanf(dobStr[:10], "%d-%d-%d", &year, &month, &day)
	now := time.Now()
	age := now.Year() - year
	if now.Month() < time.Month(month) || (now.Month() == time.Month(month) && now.Day() < day) {
		age--
	}
	return age
}

// educationSufficient checks whether the user's education meets a job's minimum requirement
func educationSufficient(userEdu, jobReq string) bool {
	order := map[string]int{
		"None": 0, "Primary": 1, "10th Pass": 2, "12th Pass": 3,
		"Undergraduate": 4, "Graduate": 5, "Post Graduate": 6,
	}
	userLevel, userOk := order[userEdu]
	requiredLevel, jobOk := order[jobReq]
	if !jobOk {
		return true // Unknown requirement — include the job
	}
	if !userOk {
		return false
	}
	return userLevel >= requiredLevel
}

// ── Gemini API Types ──────────────────────────────────────────────────────────

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}
type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}
type geminiPart struct {
	Text string `json:"text"`
}
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// generateGeminiExplanation calls the Gemini API for a personalized explanation.
// Falls back to a rule-based explanation if the call fails or API key is missing.
func generateGeminiExplanation(profile models.UserProfile, title, description, benefits, itemType string) string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return generateFallbackExplanation(profile, title, itemType)
	}

	benefitsSummary := benefits
	if len(benefitsSummary) > 150 {
		benefitsSummary = benefitsSummary[:150]
	}

	prompt := fmt.Sprintf(
		"You are a government scheme advisor for Indian citizens. Given this citizen profile and a %s, write ONE concise sentence (max 20 words) explaining why this %s is suitable for them. Be specific and helpful.\n\nCitizen: %s, %s, from %s, %s caste, annual income ₹%.0f, %s, %s education.\n%s: %s\nBenefits: %s\n\nResponse (1 sentence only):",
		itemType, itemType,
		profile.FullName, profile.Gender, profile.State, profile.CasteCategory,
		profile.AnnualIncome, profile.Occupation, profile.EducationLevel,
		strings.Title(itemType), title, benefitsSummary,
	)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return generateFallbackExplanation(profile, title, itemType)
	}

	apiURL := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s",
		apiKey,
	)

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		log.Printf("Gemini API error for '%s': %v", title, err)
		return generateFallbackExplanation(profile, title, itemType)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return generateFallbackExplanation(profile, title, itemType)
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(body, &gemResp); err != nil || len(gemResp.Candidates) == 0 {
		return generateFallbackExplanation(profile, title, itemType)
	}
	if len(gemResp.Candidates[0].Content.Parts) == 0 {
		return generateFallbackExplanation(profile, title, itemType)
	}

	explanation := strings.TrimSpace(gemResp.Candidates[0].Content.Parts[0].Text)
	if explanation == "" {
		return generateFallbackExplanation(profile, title, itemType)
	}
	return explanation
}

// generateFallbackExplanation provides a rule-based explanation when Gemini is unavailable
func generateFallbackExplanation(profile models.UserProfile, title, itemType string) string {
	switch itemType {
	case "scheme":
		return fmt.Sprintf("Based on your profile as a %s from %s, you meet the eligibility criteria for this scheme.", profile.Occupation, profile.State)
	case "job":
		return fmt.Sprintf("Your %s qualification matches the requirements for this %s opening.", profile.EducationLevel, title)
	default:
		return fmt.Sprintf("This %s matches your profile and location in %s.", itemType, profile.State)
	}
}
