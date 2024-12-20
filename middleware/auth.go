// middleware/auth.go
package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

// RateLimiter implements a simple in-memory rate limiting
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

// Global rate limiter instance
var limiter = NewRateLimiter()

// cleanOldRequests removes requests older than 1 minute
func (rl *RateLimiter) cleanOldRequests(key string) {
	now := time.Now()
	valid := []time.Time{}
	for _, t := range rl.requests[key] {
		if now.Sub(t) < time.Minute {
			valid = append(valid, t)
		}
	}
	rl.requests[key] = valid
}

// isAllowed checks if a request is allowed based on rate limiting
func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.cleanOldRequests(key)

	// Allow 60 requests per minute
	if len(rl.requests[key]) >= 60 {
		return false
	}

	rl.requests[key] = append(rl.requests[key], time.Now())
	return true
}

func RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		expectedKey := os.Getenv("BLOG_API_KEY")

		// Check if API key is present
		if apiKey == "" {
			respondError(w, "API key is required", http.StatusUnauthorized)
			return
		}

		// Rate limiting check
		if !limiter.isAllowed(apiKey) {
			respondError(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		// Verify API key
		if apiKey != expectedKey {
			respondError(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
