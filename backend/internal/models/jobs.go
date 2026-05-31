package models

import "time"

type GovtJob struct {
	ID                     int      `json:"id"`
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
	CreatedAt              time.Time`json:"created_at"`
	UpdatedAt              time.Time`json:"updated_at"`
}
