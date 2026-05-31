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

	// Setup Rate Limiters
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)
	translateLimiter := middleware.NewRateLimiter(30, time.Minute)
	eligibilityLimiter := middleware.NewRateLimiter(50, time.Minute)

	// Setup Routes with middleware
	mux := http.NewServeMux()

	// Public routes (no auth required)
	mux.HandleFunc("/api/schemes", handlers.GetSchemesHandler)
	mux.HandleFunc("/api/schemes/", handlers.GetSchemeDetailsHandler) // Handles /api/schemes/:id
	mux.HandleFunc("/api/jobs", handlers.GetJobsHandler)
	mux.HandleFunc("/api/jobs/", handlers.GetJobDetailsHandler) // Handles /api/jobs/:id
	mux.Handle("/api/eligibility-check", eligibilityLimiter.Limit(http.HandlerFunc(handlers.CheckEligibilityHandler)))
	mux.HandleFunc("/api/auth/register", handlers.RegisterHandler)
	mux.Handle("/api/auth/login", loginLimiter.Limit(http.HandlerFunc(handlers.LoginHandler)))

	// Auth-protected routes (JWT middleware validates token and injects user_id)
	mux.Handle("/api/user/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.UserProfileHandler)))
	mux.Handle("/api/user/saved", middleware.AuthMiddleware(http.HandlerFunc(handlers.SavedSchemesHandler)))
	mux.Handle("/api/user/apply", middleware.AuthMiddleware(http.HandlerFunc(handlers.ApplySchemeHandler)))
	mux.Handle("/api/user/applications", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserApplicationsHandler)))
	mux.Handle("/api/translate", translateLimiter.Limit(middleware.AuthMiddleware(http.HandlerFunc(handlers.TranslateHandler))))

	// Admin-protected routes (JWT validates, then AdminMiddleware verifies is_admin claim)
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
	mux.Handle("/api/admin/notifications", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminNotificationsHandler))))
	mux.Handle("/api/admin/applications", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminApplicationsHandler))))
	mux.Handle("/api/admin/applications/status", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.AdminApplicationStatusHandler))))


	// Global middleware chain: Logging → CORS → Routes
	handler := middleware.LoggingMiddleware(middleware.CorsMiddleware(mux))

	// Fetch PORT from environment
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
