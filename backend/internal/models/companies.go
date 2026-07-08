package models

import "time"

// Company represents an employer organization
type Company struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	LogoURL     string    `json:"logo_url" db:"logo_url"`
	Website     string    `json:"website" db:"website"`
	Industry    string    `json:"industry" db:"industry"`
	CompanySize string    `json:"company_size" db:"company_size"`
	Location    string    `json:"location" db:"location"`
	FoundedYear int       `json:"founded_year" db:"founded_year"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	JobCount    int       `json:"job_count,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PrivateJob represents a private sector job listing
type PrivateJob struct {
	ID                     int       `json:"id" db:"id"`
	Title                  string    `json:"title" db:"title"`
	TitleHi                string    `json:"title_hi" db:"title_hi"`
	TitleMr                string    `json:"title_mr" db:"title_mr"`
	CompanyID              *int      `json:"company_id" db:"company_id"`
	CompanyName            string    `json:"company_name" db:"company_name"`
	CompanyLogoURL         string    `json:"company_logo_url" db:"company_logo_url"`
	Description            string    `json:"description" db:"description"`
	Requirements           string    `json:"requirements" db:"requirements"`
	Responsibilities       string    `json:"responsibilities" db:"responsibilities"`
	SalaryMin              float64   `json:"salary_min" db:"salary_min"`
	SalaryMax              float64   `json:"salary_max" db:"salary_max"`
	SalaryCurrency         string    `json:"salary_currency" db:"salary_currency"`
	JobType                string    `json:"job_type" db:"job_type"`
	WorkMode               string    `json:"work_mode" db:"work_mode"`
	Location               string    `json:"location" db:"location"`
	ExperienceMin          int       `json:"experience_min" db:"experience_min"`
	ExperienceMax          int       `json:"experience_max" db:"experience_max"`
	EducationQualification string    `json:"education_qualification" db:"education_qualification"`
	Skills                 []string  `json:"skills" db:"skills"`
	EmploymentType         string    `json:"employment_type" db:"employment_type"`
	ApplicationStartDate   string    `json:"application_start_date" db:"application_start_date"`
	ApplicationEndDate     string    `json:"application_end_date" db:"application_end_date"`
	OfficialWebsite        string    `json:"official_website" db:"official_website"`
	ApplyLink              string    `json:"apply_link" db:"apply_link"`
	Stipend                string    `json:"stipend" db:"stipend"`
	PrizePool              string    `json:"prize_pool" db:"prize_pool"`
	IsActive               bool      `json:"is_active" db:"is_active"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// NewsletterSubscriber represents an email newsletter subscriber
type NewsletterSubscriber struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	SubscribedAt time.Time `json:"subscribed_at" db:"subscribed_at"`
}

// AuditLog records admin actions
type AuditLog struct {
	ID           int       `json:"id" db:"id"`
	UserID       *int      `json:"user_id" db:"user_id"`
	Action       string    `json:"action" db:"action"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	ResourceID   *int      `json:"resource_id" db:"resource_id"`
	Details      string    `json:"details" db:"details"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
