// middleware/auth.go
package middleware

import (
	"log"
	"net/http"
	"os"
)

func RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		expectedKey := os.Getenv("BLOG_API_KEY")

		log.Printf("Received API Key: %s", apiKey)
		log.Printf("Expected API Key: %s", expectedKey)

		if apiKey == "" || apiKey != expectedKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
