package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

// Context key type to avoid collisions
type contextKey string

const userIDKey contextKey = "user_id"

// jwtSecret holds the loaded secret; initialized once
var (
	jwtSecret string
	secretOnce sync.Once
)

// loadJWTSecret reads JWT_SECRET from environment once. Panics if missing.
func loadJWTSecret() {
	secretOnce.Do(func() {
		jwtSecret = os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("FATAL: JWT_SECRET environment variable is not set. Server cannot start without it.")
		}
	})
}

// GetJWTSecret returns the loaded JWT secret for signing tokens
func GetJWTSecret() string {
	loadJWTSecret()
	return jwtSecret
}

// AuthMiddleware validates JWT Bearer tokens and injects user_id into context
func AuthMiddleware(next http.Handler) http.Handler {
	loadJWTSecret()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Missing or invalid authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid or expired session token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid token format claims",
			})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Missing user session identification in token",
			})
			return
		}

		// Inject user_id into request context
		ctx := context.WithValue(r.Context(), userIDKey, int(userIDFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extracts user_id from the request context
func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in request context")
	}
	return userID, nil
}
