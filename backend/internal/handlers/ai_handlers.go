package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/middleware"
	"yojana-portal/backend/internal/models"
)

// Gemini API structures for AI Handers
type aiGeminiRequest struct {
	Contents []aiGeminiContent `json:"contents"`
}

type aiGeminiContent struct {
	Parts []aiGeminiPart `json:"parts"`
}

type aiGeminiPart struct {
	Text       string        `json:"text,omitempty"`
	InlineData *aiInlineData `json:"inlineData,omitempty"`
}

type aiInlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type aiGeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Helper to call Gemini API and get response text
func callGeminiAPI(parts []aiGeminiPart) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	reqBody := aiGeminiRequest{
		Contents: []aiGeminiContent{
			{Parts: parts},
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	apiURL := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s",
		apiKey,
	)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(body))
	}

	var gemResp aiGeminiResponse
	if err := json.Unmarshal(body, &gemResp); err != nil {
		return "", err
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini API")
	}

	return gemResp.Candidates[0].Content.Parts[0].Text, nil
}

// Clean markdown block wrappers if Gemini responds with them
func cleanJSONResponse(input string) string {
	cleaned := strings.TrimSpace(input)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimSuffix(cleaned, "```")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
		cleaned = strings.TrimSuffix(cleaned, "```")
	}
	return strings.TrimSpace(cleaned)
}

// ResumeAnalyzeHandler analyzes a resume (supports text paste or PDF file upload)
func ResumeAnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, _ := middleware.GetUserIDFromContext(r.Context())

	var resumeText string
	var pdfBase64 string
	var mimeType string
	var filename string
	var fileSize int64

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form for file upload
		err := r.ParseMultipartForm(10 << 20) // 10MB limit
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "Failed to parse multipart form: "+err.Error())
			return
		}

		file, handler, err := r.FormFile("resume")
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "Missing 'resume' file parameter")
			return
		}
		defer file.Close()

		filename = handler.Filename
		fileSize = handler.Size
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".pdf" && ext != ".txt" {
			writeJSONError(w, http.StatusBadRequest, "Only .pdf and .txt files are supported")
			return
		}

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to read uploaded file")
			return
		}

		if ext == ".pdf" {
			mimeType = "application/pdf"
			pdfBase64 = base64.StdEncoding.EncodeToString(fileBytes)
		} else {
			resumeText = string(fileBytes)
		}

		// Optionally save to disk if user is logged in
		if userID > 0 {
			os.MkdirAll("uploads/resumes", 0755)
			diskPath := filepath.Join("uploads/resumes", fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), ext))
			err = os.WriteFile(diskPath, fileBytes, 0644)
			if err == nil {
				_, dbErr := db.DB.Exec(`
					INSERT INTO resume_uploads (user_id, original_filename, file_path, text_content, file_size_bytes, uploaded_at)
					VALUES ($1, $2, $3, $4, $5, NOW())`,
					userID, filename, diskPath, resumeText, fileSize,
				)
				if dbErr != nil {
					log.Printf("[ResumeAnalyze] DB error saving upload metadata: %v", dbErr)
				}
			} else {
				log.Printf("[ResumeAnalyze] Failed to write file to disk: %v", err)
			}
		}
	} else {
		// Parse JSON paste request
		var req struct {
			ResumeText string `json:"resume_text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ResumeText == "" {
			writeJSONError(w, http.StatusBadRequest, "Valid resume_text field is required")
			return
		}
		resumeText = req.ResumeText
	}

	// Prepare Gemini API Prompt
	var parts []aiGeminiPart
	systemPrompt := "You are a professional ATS (Applicant Tracking System) parser and career coach. " +
		"Analyze the following resume and return a JSON object with: " +
		"\"ats_score\" (integer between 0 and 100), " +
		"\"strength\" (string: \"weak\", \"moderate\", or \"strong\"), " +
		"\"missing_keywords\" (array of strings), " +
		"\"suggestions\" (array of strings), " +
		"\"summary\" (string summary of candidate's profile), " +
		"\"skills_found\" (array of strings)." +
		"Ensure the response is ONLY a valid JSON block, no explanations outside of the JSON."

	parts = append(parts, aiGeminiPart{Text: systemPrompt})

	if pdfBase64 != "" {
		parts = append(parts, aiGeminiPart{
			InlineData: &aiInlineData{
				MimeType: mimeType,
				Data:     pdfBase64,
			},
		})
	} else {
		parts = append(parts, aiGeminiPart{Text: "Resume text:\n" + resumeText})
	}

	rawResp, err := callGeminiAPI(parts)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gemini API error: "+err.Error())
		return
	}

	cleanedJSON := cleanJSONResponse(rawResp)

	// Verify it's valid JSON
	var parsedResp map[string]interface{}
	if err := json.Unmarshal([]byte(cleanedJSON), &parsedResp); err != nil {
		// Return raw response as fallback text wrapper
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{
			"ats_score":        70,
			"strength":         "moderate",
			"missing_keywords": []string{"Optimization", "Scalability"},
			"suggestions":      []string{"Add quantitative metrics.", "Include professional certifications."},
			"summary":          "Detailed resume provided. Unable to parse structured JSON natively.",
			"skills_found":     []string{"General Professional Skills"},
			"raw_ai_text":      rawResp,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cleanedJSON))
}

// CareerAdvisorHandler generates career roadmap, roles, learning paths
func CareerAdvisorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
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
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to load user profile: "+err.Error())
		return
	}

	// Build prompt
	prompt := fmt.Sprintf(
		"You are a professional career advisor in India. Provide a detailed career guide for this citizen:\n"+
			"Name: %s\nAge: %d\nGender: %s\nState: %s\nEducation: %s\nOccupation: %s\nEmployee Type: %s\nAnnual Income: ₹%.2f\nDisabled: %t\n\n"+
			"Respond with a JSON object containing:\n"+
			"1. \"suitable_roles\": array of objects with \"title\" (string), \"match_score\" (int 0-100), and \"description\" (string)\n"+
			"2. \"roadmap\": array of objects with \"phase\" (string, e.g. '0-6 months'), \"duration\" (string), and \"milestones\" (array of strings)\n"+
			"3. \"required_skills\": array of strings\n"+
			"4. \"certifications\": array of objects with \"name\" (string) and \"provider\" (string)\n"+
			"5. \"expected_salary_range\": object with \"min\" (numeric), \"max\" (numeric), and \"currency\" (string, e.g., 'INR')\n"+
			"6. \"learning_resources\": array of objects with \"name\" (string), \"url\" (string), and \"type\" (string: course/book/tutorial)\n\n"+
			"Do not include markdown tags outside of the JSON block.",
		profile.FullName, calculateAge(profile.DateOfBirth), profile.Gender, profile.State,
		profile.EducationLevel, profile.Occupation, profile.EmployeeType, profile.AnnualIncome, profile.IsDisabled,
	)

	parts := []aiGeminiPart{{Text: prompt}}
	rawResp, err := callGeminiAPI(parts)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gemini API error: "+err.Error())
		return
	}

	cleanedJSON := cleanJSONResponse(rawResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cleanedJSON))
}

// SkillGapHandler compares user profile vs job details
func SkillGapHandler(w http.ResponseWriter, r *http.Request) {
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
		JobTitle       string   `json:"job_title"`
		JobDescription string   `json:"job_description"`
		RequiredSkills []string `json:"required_skills"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var profile models.UserProfile
	profileQuery := `
		SELECT p.education_level, p.occupation, p.employee_type
		FROM user_profiles p WHERE p.user_id = $1`
	_ = db.DB.QueryRow(profileQuery, userID).Scan(
		&profile.EducationLevel, &profile.Occupation, &profile.EmployeeType,
	)

	prompt := fmt.Sprintf(
		"Compare this citizen profile:\n"+
			"Education: %s\nCurrent Occupation: %s\n\n"+
			"Against this Job post:\n"+
			"Title: %s\nDescription: %s\nRequired Skills: %s\n\n"+
			"Respond with a JSON object containing:\n"+
			"1. \"match_score\": integer (0-100)\n"+
			"2. \"matching_skills\": array of strings\n"+
			"3. \"missing_skills\": array of objects with \"skill\" (string), \"priority\" (string: 'high'/'medium'/'low'), and \"estimated_learning_time\" (string)\n"+
			"4. \"courses\": array of objects with \"name\" (string), \"provider\" (string), and \"url\" (string)\n"+
			"5. \"projects\": array of objects with \"title\" (string), \"description\" (string), and \"skills_covered\" (array of strings)\n\n"+
			"Ensure response is strictly valid JSON.",
		profile.EducationLevel, profile.Occupation, req.JobTitle, req.JobDescription, strings.Join(req.RequiredSkills, ", "),
	)

	parts := []aiGeminiPart{{Text: prompt}}
	rawResp, err := callGeminiAPI(parts)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gemini API error: "+err.Error())
		return
	}

	cleanedJSON := cleanJSONResponse(rawResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cleanedJSON))
}

// InterviewPrepHandler generates interview preparation guides
func InterviewPrepHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		JobTitle        string `json:"job_title"`
		JobDescription  string `json:"job_description"`
		ExperienceYears int    `json:"experience_years"`
		InterviewType   string `json:"interview_type"` // hr, technical, behavioral, mixed
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	prompt := fmt.Sprintf(
		"Generate interview prep questions for:\n"+
			"Job Title: %s\nDescription: %s\nExperience level: %d years\nInterview Type: %s\n\n"+
			"Respond with a JSON object containing:\n"+
			"1. \"hr_questions\": array of objects with \"question\" (string) and \"tips\" (string)\n"+
			"2. \"technical_questions\": array of objects with \"question\" (string), \"expected_answer\" (string), and \"difficulty\" (string)\n"+
			"3. \"behavioral_questions\": array of objects with \"question\" (string) and \"star_method_hint\" (string)\n"+
			"4. \"coding_questions\": array of objects with \"question\" (string) and \"hints\" (array of strings)\n\n"+
			"Only output valid JSON.",
		req.JobTitle, req.JobDescription, req.ExperienceYears, req.InterviewType,
	)

	parts := []aiGeminiPart{{Text: prompt}}
	rawResp, err := callGeminiAPI(parts)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gemini API error: "+err.Error())
		return
	}

	cleanedJSON := cleanJSONResponse(rawResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cleanedJSON))
}

// CoverLetterHandler generates personalized cover letter
func CoverLetterHandler(w http.ResponseWriter, r *http.Request) {
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
		CompanyName    string `json:"company_name"`
		JobTitle       string `json:"job_title"`
		JobDescription string `json:"job_description"`
		Tone           string `json:"tone"` // professional, enthusiastic, creative
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var profile models.UserProfile
	profileQuery := `
		SELECT p.full_name, p.education_level, p.occupation, p.state
		FROM user_profiles p WHERE p.user_id = $1`
	_ = db.DB.QueryRow(profileQuery, userID).Scan(
		&profile.FullName, &profile.EducationLevel, &profile.Occupation, &profile.State,
	)

	prompt := fmt.Sprintf(
		"Generate a tailored cover letter from %s (Education: %s, Current Occupation: %s, State: %s)\n"+
			"To Company: %s\nFor Role: %s\nJob Description: %s\nTone: %s\n\n"+
			"Respond with a JSON object containing:\n"+
			"1. \"subject_line\": string\n"+
			"2. \"cover_letter_text\": string (use \\n for newlines)\n\n"+
			"Strictly return valid JSON.",
		profile.FullName, profile.EducationLevel, profile.Occupation, profile.State,
		req.CompanyName, req.JobTitle, req.JobDescription, req.Tone,
	)

	parts := []aiGeminiPart{{Text: prompt}}
	rawResp, err := callGeminiAPI(parts)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gemini API error: "+err.Error())
		return
	}

	cleanedJSON := cleanJSONResponse(rawResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cleanedJSON))
}

// SmartSearchHandler parses natural language query and fetches jobs and schemes
func SmartSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Query == "" {
		writeJSONError(w, http.StatusBadRequest, "Valid query is required")
		return
	}

	// 1. Ask Gemini to extract structured filters from NLP query
	prompt := fmt.Sprintf(
		"Analyze this job/scheme search query: \"%s\"\n\n"+
			"Translate it into structured search filters as a JSON object containing:\n"+
			"1. \"keywords\": string (comma separated search terms, e.g. 'React, Developer' or 'Farmers')\n"+
			"2. \"location\": string (city or state name, e.g. 'Mumbai' or 'Maharashtra')\n"+
			"3. \"job_type\": string ('govt', 'private', 'internship', 'walkin', 'hackathon', or 'all')\n"+
			"4. \"work_mode\": string ('remote', 'hybrid', 'onsite', or 'all')\n"+
			"5. \"salary_max\": numeric (maximum salary limit, 0 if not specified)\n"+
			"6. \"experience_max\": integer (maximum experience limit, -1 if not specified)\n"+
			"7. \"education\": string ('10th Pass', '12th Pass', 'Graduate', 'Post Graduate', or 'all')\n\n"+
			"Ensure output is ONLY valid JSON.",
		req.Query,
	)

	parts := []aiGeminiPart{{Text: prompt}}
	rawResp, err := callGeminiAPI(parts)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Gemini parsing error: "+err.Error())
		return
	}

	cleanedJSON := cleanJSONResponse(rawResp)

	type SearchFilters struct {
		Keywords      string  `json:"keywords"`
		Location      string  `json:"location"`
		JobType       string  `json:"job_type"`
		WorkMode      string  `json:"work_mode"`
		SalaryMax     float64 `json:"salary_max"`
		ExperienceMax int     `json:"experience_max"`
		Education     string  `json:"education"`
	}

	var filters SearchFilters
	if err := json.Unmarshal([]byte(cleanedJSON), &filters); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to parse search filters: "+err.Error())
		return
	}

	// 2. Query Government Schemes using keywords/location
	schemes := []models.Scheme{}
	schemeQuery := `
		SELECT s.id, s.title, s.title_hi, s.title_mr, s.description, s.description_hi, s.description_mr,
		       s.category_id, c.name, c.name_hi, c.name_mr, s.government_level, s.state, s.benefits,
		       s.application_start_date, s.application_end_date, s.official_website, s.apply_link, s.is_active
		FROM schemes s
		JOIN scheme_categories c ON s.category_id = c.id
		WHERE s.is_active = true`

	var schemeArgs []interface{}
	argIdx := 1

	if filters.Keywords != "" {
		terms := strings.Split(filters.Keywords, ",")
		if len(terms) > 0 {
			term := strings.TrimSpace(terms[0])
			schemeQuery += fmt.Sprintf(" AND (s.title ILIKE $%d OR s.description ILIKE $%d OR c.name ILIKE $%d)", argIdx, argIdx, argIdx)
			schemeArgs = append(schemeArgs, "%"+term+"%")
			argIdx++
		}
	}
	if filters.Location != "" {
		schemeQuery += fmt.Sprintf(" AND (s.state ILIKE $%d OR s.government_level = 'central')", argIdx)
		schemeArgs = append(schemeArgs, "%"+filters.Location+"%")
		argIdx++
	}

	schemeQuery += " LIMIT 5"
	sRows, err := db.DB.Query(schemeQuery, schemeArgs...)
	if err == nil {
		defer sRows.Close()
		for sRows.Next() {
			var s models.Scheme
			err := sRows.Scan(
				&s.ID, &s.Title, &s.TitleHi, &s.TitleMr, &s.Description, &s.DescriptionHi, &s.DescriptionMr,
				&s.CategoryID, &s.CategoryName, &s.CategoryNameHi, &s.CategoryNameMr, &s.GovernmentLevel, &s.State, &s.Benefits,
				&s.ApplicationStartDate, &s.ApplicationEndDate, &s.OfficialWebsite, &s.ApplyLink, &s.IsActive,
			)
			if err == nil {
				schemes = append(schemes, s)
			}
		}
	}

	// 3. Query Government Jobs
	govtJobs := []models.GovtJob{}
	if filters.JobType == "all" || filters.JobType == "govt" {
		jobQuery := `
			SELECT id, title, organization, department, vacancies, education_qualification,
			       experience_required, application_start_date, application_end_date, official_website, apply_link
			FROM govt_jobs WHERE is_active = true`

		var jobArgs []interface{}
		jobArgIdx := 1

		if filters.Keywords != "" {
			terms := strings.Split(filters.Keywords, ",")
			if len(terms) > 0 {
				term := strings.TrimSpace(terms[0])
				jobQuery += fmt.Sprintf(" AND (title ILIKE $%d OR organization ILIKE $%d)", jobArgIdx, jobArgIdx)
				jobArgs = append(jobArgs, "%"+term+"%")
				jobArgIdx++
			}
		}
		if filters.Education != "" && filters.Education != "all" {
			jobQuery += fmt.Sprintf(" AND education_qualification = $%d", jobArgIdx)
			jobArgs = append(jobArgs, filters.Education)
			jobArgIdx++
		}

		jobQuery += " LIMIT 5"
		jRows, err := db.DB.Query(jobQuery, jobArgs...)
		if err == nil {
			defer jRows.Close()
			for jRows.Next() {
				var j models.GovtJob
				var startDate, endDate time.Time
				err := jRows.Scan(
					&j.ID, &j.Title, &j.Organization, &j.Department, &j.Vacancies, &j.EducationQualification,
					&j.ExperienceRequired, &startDate, &endDate, &j.OfficialWebsite, &j.ApplyLink,
				)
				if err == nil {
					j.ApplicationStartDate = startDate.Format("2006-01-02")
					j.ApplicationEndDate = endDate.Format("2006-01-02")
					govtJobs = append(govtJobs, j)
				}
			}
		}
	}

	// 4. Query Private Sector Jobs
	privateJobs := []models.PrivateJob{}
	if filters.JobType == "all" || filters.JobType != "govt" {
		pQuery := `
			SELECT pj.id, pj.title, c.name, c.logo_url, pj.job_type, pj.work_mode, pj.location,
			       pj.salary_min, pj.salary_max, pj.experience_min, pj.experience_max,
			       pj.education_qualification, pj.apply_link
			FROM private_jobs pj
			LEFT JOIN companies c ON pj.company_id = c.id
			WHERE pj.is_active = true`

		var pArgs []interface{}
		pArgIdx := 1

		if filters.Keywords != "" {
			terms := strings.Split(filters.Keywords, ",")
			if len(terms) > 0 {
				term := strings.TrimSpace(terms[0])
				pQuery += fmt.Sprintf(" AND (pj.title ILIKE $%d OR c.name ILIKE $%d)", pArgIdx, pArgIdx)
				pArgs = append(pArgs, "%"+term+"%")
				pArgIdx++
			}
		}
		if filters.JobType != "all" && filters.JobType != "govt" {
			pQuery += fmt.Sprintf(" AND pj.job_type = $%d", pArgIdx)
			pArgs = append(pArgs, filters.JobType)
			pArgIdx++
		}
		if filters.Location != "" {
			pQuery += fmt.Sprintf(" AND pj.location ILIKE $%d", pArgIdx)
			pArgs = append(pArgs, "%"+filters.Location+"%")
			pArgIdx++
		}
		if filters.WorkMode != "" && filters.WorkMode != "all" {
			pQuery += fmt.Sprintf(" AND pj.work_mode = $%d", pArgIdx)
			pArgs = append(pArgs, filters.WorkMode)
			pArgIdx++
		}

		pQuery += " LIMIT 5"
		pRows, err := db.DB.Query(pQuery, pArgs...)
		if err == nil {
			defer pRows.Close()
			for pRows.Next() {
				var j models.PrivateJob
				err := pRows.Scan(
					&j.ID, &j.Title, &j.CompanyName, &j.CompanyLogoURL, &j.JobType, &j.WorkMode, &j.Location,
					&j.SalaryMin, &j.SalaryMax, &j.ExperienceMin, &j.ExperienceMax,
					&j.EducationQualification, &j.ApplyLink,
				)
				if err == nil {
					privateJobs = append(privateJobs, j)
				}
			}
		}
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success":      true,
		"filters":      filters,
		"schemes":      schemes,
		"govt_jobs":    govtJobs,
		"private_jobs": privateJobs,
	})
}
