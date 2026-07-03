package models

import "time"

// GovtJob represents a government job opening
type GovtJob struct {
	ID                     int            `json:"id" db:"id"`
	Title                  string         `json:"title" db:"title"`
	TitleHi                string         `json:"title_hi" db:"title_hi"`
	TitleMr                string         `json:"title_mr" db:"title_mr"`
	Organization           string         `json:"organization" db:"organization"`
	Department             string         `json:"department" db:"department"`
	Vacancies              int            `json:"vacancies" db:"vacancies"`
	CategoryVacancies      map[string]int `json:"category_vacancies" db:"category_vacancies"`
	EducationQualification string         `json:"education_qualification" db:"education_qualification"`
	ExperienceRequired     string         `json:"experience_required" db:"experience_required"`
	RequiredDocuments      []string       `json:"required_documents" db:"required_documents"`
	ApplicationStartDate   string         `json:"application_start_date" db:"application_start_date"`
	ApplicationEndDate     string         `json:"application_end_date" db:"application_end_date"`
	OfficialWebsite        string         `json:"official_website" db:"official_website"`
	ApplyLink              string         `json:"apply_link" db:"apply_link"`
	ApplicationFee         string         `json:"application_fee" db:"application_fee"`
	IsActive               bool           `json:"is_active" db:"is_active"`
	CreatedAt              time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at" db:"updated_at"`
}
