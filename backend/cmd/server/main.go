package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/handlers"
	"yojana-portal/backend/internal/middleware"
	"yojana-portal/backend/internal/ws"
)

func main() {
	// Load environment variables from .env file (if present)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	// Initialize Database connection
	err = db.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.DB.Close()

	// Initialize WebSocket Hub for real-time notifications
	hub := ws.NewHub()
	go hub.Run()

	// Ensure uploads directory exists
	os.MkdirAll("uploads/resumes", 0755)

	// Setup Rate Limiters
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)
	translateLimiter := middleware.NewRateLimiter(30, time.Minute)
	eligibilityLimiter := middleware.NewRateLimiter(50, time.Minute)
	aiLimiter := middleware.NewRateLimiter(15, time.Minute)

	// Setup Routes with middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlers.HealthHandler)

	// Public routes (no auth required)
	mux.HandleFunc("/api/schemes", handlers.GetSchemesHandler)
	mux.HandleFunc("/api/schemes/", handlers.GetSchemeDetailsHandler)
	mux.HandleFunc("/api/jobs", handlers.GetJobsHandler)
	mux.HandleFunc("/api/jobs/", handlers.GetJobDetailsHandler)
	mux.HandleFunc("/api/companies", handlers.GetCompaniesHandler)
	mux.HandleFunc("/api/companies/", handlers.GetCompanyDetailsHandler)
	mux.HandleFunc("/api/private-jobs", handlers.GetPrivateJobsHandler)
	mux.HandleFunc("/api/private-jobs/", handlers.GetPrivateJobDetailsHandler)
	mux.HandleFunc("/api/newsletter/subscribe", handlers.NewsletterSubscribeHandler)
	mux.Handle("/api/eligibility-check", eligibilityLimiter.Limit(http.HandlerFunc(handlers.CheckEligibilityHandler)))
	mux.HandleFunc("/api/auth/register", handlers.RegisterHandler)
	mux.Handle("/api/auth/login", loginLimiter.Limit(http.HandlerFunc(handlers.LoginHandler)))
	mux.HandleFunc("/api/auth/logout", handlers.LogoutHandler)

	// Auth-protected user routes
	mux.Handle("/api/user/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.UserProfileHandler)))
	mux.Handle("/api/user/saved", middleware.AuthMiddleware(http.HandlerFunc(handlers.SavedSchemesHandler)))
	mux.Handle("/api/user/apply", middleware.AuthMiddleware(http.HandlerFunc(handlers.ApplySchemeHandler)))
	mux.Handle("/api/user/applications", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserApplicationsHandler)))
	mux.Handle("/api/user/apply-job", middleware.AuthMiddleware(http.HandlerFunc(handlers.ApplyJobHandler)))
	mux.Handle("/api/user/job-applications", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserJobApplicationsHandler)))
	mux.Handle("/api/user/recommendations", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetRecommendationsHandler)))
	mux.Handle("/api/user/notifications", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserNotificationsHandler)))
	mux.Handle("/api/user/notifications/read", middleware.AuthMiddleware(http.HandlerFunc(handlers.MarkNotificationsReadHandler)))
	mux.Handle("/api/translate", translateLimiter.Limit(middleware.AuthMiddleware(http.HandlerFunc(handlers.TranslateHandler))))
	mux.Handle("/api/user/apply-private-job", middleware.AuthMiddleware(http.HandlerFunc(handlers.ApplyPrivateJobHandler)))
	mux.Handle("/api/user/private-job-applications", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserPrivateJobApplicationsHandler)))

	// AI routes (auth protected, rate limited)
	mux.Handle("/api/ai/resume-analyze", aiLimiter.Limit(middleware.AuthMiddleware(http.HandlerFunc(handlers.ResumeAnalyzeHandler))))
	mux.Handle("/api/ai/career-advisor", aiLimiter.Limit(middleware.AuthMiddleware(http.HandlerFunc(handlers.CareerAdvisorHandler))))
	mux.Handle("/api/ai/skill-gap", aiLimiter.Limit(middleware.AuthMiddleware(http.HandlerFunc(handlers.SkillGapHandler))))
	mux.Handle("/api/ai/interview-prep", aiLimiter.Limit(middleware.AuthMiddleware(http.HandlerFunc(handlers.InterviewPrepHandler))))
	mux.Handle("/api/ai/cover-letter", aiLimiter.Limit(middleware.AuthMiddleware(http.HandlerFunc(handlers.CoverLetterHandler))))
	mux.Handle("/api/ai/smart-search", aiLimiter.Limit(http.HandlerFunc(handlers.SmartSearchHandler))) // Smart search itself is public

	// WebSocket endpoint for real-time notifications
	mux.Handle("/ws", middleware.AuthMiddleware(handlers.WebSocketHandler(hub)))

	// Admin-protected routes
	mux.Handle("/api/admin/analytics", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.GetAdminAnalyticsHandler))))
	mux.Handle("/api/admin/schemes", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminSchemesHandler))))
	mux.Handle("/api/admin/schemes/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminSchemeDetailsHandler))))
	mux.Handle("/api/admin/jobs", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminJobsHandler))))
	mux.Handle("/api/admin/jobs/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminJobDetailsHandler))))
	mux.Handle("/api/admin/categories", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminCategoriesHandler))))
	mux.Handle("/api/admin/categories/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminCategoryDetailsHandler))))
	mux.Handle("/api/admin/users", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminUsersHandler))))
	mux.Handle("/api/admin/users/toggle-active", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminUserToggleHandler))))
	mux.Handle("/api/admin/users/admin", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminCreateHandler))))
	mux.Handle("/api/admin/notifications", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminSendNotificationHandler))))
	mux.Handle("/api/admin/applications", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminApplicationsHandler))))
	mux.Handle("/api/admin/applications/status", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminApplicationStatusHandler))))
	mux.Handle("/api/admin/companies", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminCompaniesHandler))))
	mux.Handle("/api/admin/companies/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminCompanyDetailsHandler))))
	mux.Handle("/api/admin/private-jobs", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminPrivateJobsHandler))))
	mux.Handle("/api/admin/private-jobs/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminPrivateJobDetailsHandler))))
	mux.Handle("/api/admin/audit-logs", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.GetAdminAuditLogsHandler))))

	// Global middleware chain: Logging → CORS → Security Headers → Body Size Limit → Routes
	handler := middleware.LoggingMiddleware(
		middleware.CorsMiddleware(
			middleware.SecurityHeadersMiddleware(
				middleware.RequestSizeLimiter(2<<20)(mux), // 2MB max body
			),
		),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Yojana Portal Backend running on http://localhost:%s", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
