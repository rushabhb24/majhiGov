package models

import (
	"time"

	"github.com/lib/pq"
)

// User represents credentials and auth state
type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserProfile represents personal details used in eligibility checking
type UserProfile struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	FullName       string    `json:"full_name" db:"full_name"`
	DateOfBirth    string    `json:"date_of_birth" db:"date_of_birth"` // YYYY-MM-DD
	Gender         string    `json:"gender" db:"gender"`
	State          string    `json:"state" db:"state"`
	District       string    `json:"district" db:"district"`
	CasteCategory  string    `json:"caste_category" db:"caste_category"`
	AnnualIncome   float64   `json:"annual_income" db:"annual_income"`
	Occupation     string    `json:"occupation" db:"occupation"`
	EmployeeType   string    `json:"employee_type" db:"employee_type"`
	EducationLevel string    `json:"education_level" db:"education_level"`
	IsDisabled     bool      `json:"is_disabled" db:"is_disabled"`
	DisabilityType string    `json:"disability_type" db:"disability_type"`
	AvatarURL      string    `json:"avatar_url" db:"avatar_url"`
	Email          string    `json:"email" db:"email"`
	Phone          string    `json:"phone" db:"phone"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// SchemeCategory represents a sector (e.g. Farmers, Students)
type SchemeCategory struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	NameHi      string    `json:"name_hi" db:"name_hi"`
	NameMr      string    `json:"name_mr" db:"name_mr"`
	Icon        string    `json:"icon" db:"icon"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Scheme represents a government scheme with multilingual fields
type Scheme struct {
	ID                   int                 `json:"id" db:"id"`
	Title                string              `json:"title" db:"title"`
	TitleHi              string              `json:"title_hi" db:"title_hi"`
	TitleMr              string              `json:"title_mr" db:"title_mr"`
	Description          string              `json:"description" db:"description"`
	DescriptionHi        string              `json:"description_hi" db:"description_hi"`
	DescriptionMr        string              `json:"description_mr" db:"description_mr"`
	CategoryID           int                 `json:"category_id" db:"category_id"`
	CategoryName         string              `json:"category_name" db:"category_name"`     // From JOIN
	CategoryNameHi       string              `json:"category_name_hi" db:"category_name_hi"` // From JOIN
	CategoryNameMr       string              `json:"category_name_mr" db:"category_name_mr"` // From JOIN
	GovernmentLevel      string              `json:"government_level" db:"government_level"` // central, state
	State                *string             `json:"state" db:"state"`                       // null if central
	Benefits             string              `json:"benefits" db:"benefits"`
	ApplicationStartDate string              `json:"application_start_date" db:"application_start_date"`
	ApplicationEndDate   string              `json:"application_end_date" db:"application_end_date"`
	OfficialWebsite      string              `json:"official_website" db:"official_website"`
	ApplyLink            string              `json:"apply_link" db:"apply_link"`
	IsActive             bool                `json:"is_active" db:"is_active"`
	CreatedAt            time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at" db:"updated_at"`
	Documents            []SchemeDocument    `json:"documents,omitempty"`
	FAQs                 []SchemeFAQ         `json:"faqs,omitempty"`
	Eligibility          *EligibilityCriteria `json:"eligibility,omitempty"`
}

// EligibilityCriteria represents eligibility constraints per scheme
type EligibilityCriteria struct {
	ID                 int            `json:"id" db:"id"`
	SchemeID           int            `json:"scheme_id" db:"scheme_id"`
	MinAge             int            `json:"min_age" db:"min_age"`
	MaxAge             int            `json:"max_age" db:"max_age"`
	Gender             string         `json:"gender" db:"gender"` // all, male, female, other
	CasteCategories    pq.StringArray `json:"caste_categories" db:"caste_categories"`
	MinIncome          float64        `json:"min_income" db:"min_income"`
	MaxIncome          float64        `json:"max_income" db:"max_income"`
	States             pq.StringArray `json:"states" db:"states"` // null/empty = all India
	Occupations        pq.StringArray `json:"occupations" db:"occupations"`
	EmployeeTypes      pq.StringArray `json:"employee_types" db:"employee_types"`
	EducationLevels    pq.StringArray `json:"education_levels" db:"education_levels"`
	DisabilityRequired bool           `json:"disability_required" db:"disability_required"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
}

// SchemeDocument represents what papers are required
type SchemeDocument struct {
	ID             int       `json:"id" db:"id"`
	SchemeID       int       `json:"scheme_id" db:"scheme_id"`
	DocumentName   string    `json:"document_name" db:"document_name"`
	DocumentNameHi string    `json:"document_name_hi" db:"document_name_hi"`
	DocumentNameMr string    `json:"document_name_mr" db:"document_name_mr"`
	IsMandatory    bool      `json:"is_mandatory" db:"is_mandatory"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// SchemeFAQ represents frequently asked questions
type SchemeFAQ struct {
	ID         int       `json:"id" db:"id"`
	SchemeID   int       `json:"scheme_id" db:"scheme_id"`
	Question   string    `json:"question" db:"question"`
	Answer     string    `json:"answer" db:"answer"`
	QuestionHi string    `json:"question_hi" db:"question_hi"`
	AnswerHi   string    `json:"answer_hi" db:"answer_hi"`
	QuestionMr string    `json:"question_mr" db:"question_mr"`
	AnswerMr   string    `json:"answer_mr" db:"answer_mr"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// UserSavedScheme represents bookmarks
type UserSavedScheme struct {
	ID       int       `json:"id" db:"id"`
	UserID   int       `json:"user_id" db:"user_id"`
	SchemeID int       `json:"scheme_id" db:"scheme_id"`
	SavedAt  time.Time `json:"saved_at" db:"saved_at"`
}

// UserAppliedScheme represents applications
type UserAppliedScheme struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	SchemeID  int       `json:"scheme_id" db:"scheme_id"`
	Status    string    `json:"status" db:"status"` // pending, approved, rejected
	AppliedAt time.Time `json:"applied_at" db:"applied_at"`
	Notes     string    `json:"notes" db:"notes"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// OTP represents verification tokens
type OTP struct {
	ID        int       `json:"id" db:"id"`
	Contact   string    `json:"contact" db:"contact"`
	OTPCode   string    `json:"otp_code" db:"otp_code"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	IsUsed    bool      `json:"is_used" db:"is_used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Notification represents alert messages
type Notification struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Message   string    `json:"message" db:"message"`
	Type      string    `json:"type" db:"type"`
	IsRead    bool      `json:"is_read" db:"is_read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserProfileRequest represents dynamic anonymous checking structure
type UserProfileRequest struct {
	Age            int     `json:"age"`
	Gender         string  `json:"gender"`
	State          string  `json:"state"`
	Caste          string  `json:"caste"`
	AnnualIncome   float64 `json:"annual_income"`
	Occupation     string  `json:"occupation"`
	EmployeeType   string  `json:"employee_type"`
	EducationLevel string  `json:"education_level"`
	IsDisabled     bool    `json:"is_disabled"`
}

// EligibilityStatus holds matching details for a specific scheme
type EligibilityStatus struct {
	Scheme     Scheme   `json:"scheme"`
	IsEligible bool     `json:"is_eligible"`
	Reasons    []string `json:"reasons"`
}

// EligibilityResponse represents lists of eligible and ineligible schemes
type EligibilityResponse struct {
	Eligible    []EligibilityStatus `json:"eligible"`
	NotEligible []EligibilityStatus `json:"not_eligible"`
}

// RegisterRequest carries credentials and profile fields for registering
type RegisterRequest struct {
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Password       string  `json:"password"`
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
}

// LoginRequest carries user login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse carries session JWT and logged-in profile data
type AuthResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Token   string       `json:"token,omitempty"`
	Profile *UserProfile `json:"profile,omitempty"`
}

