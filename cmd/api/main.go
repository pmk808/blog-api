// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/pmk808/blog-api/internal/handler"
	"github.com/pmk808/blog-api/internal/storage"
	custommiddleware "github.com/pmk808/blog-api/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           Blog API
// @version         1.0
// @description     API for managing technical blog content
// @termsOfService  http://swagger.io/terms/

// @contact.name   Your Name
// @contact.url    http://your-url.com
// @contact.email  your@email.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}
	// Initialize database
	db, err := storage.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize storage and handlers
	postStore := storage.NewPostStore(db)
	postHandler := handler.NewPostHandler(postStore)

	// Create Chi router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	r.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-API-Key"))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Public routes
	r.Route("/api", func(r chi.Router) {
		// Add versioning prefix
		r.Route("/v1", func(r chi.Router) {
			// Posts routes
			r.Route("/posts", func(r chi.Router) {
				r.Get("/", postHandler.ListPosts)     // GET /api/v1/posts
				r.Get("/{slug}", postHandler.GetPost) // GET /api/v1/posts/{slug}
			})

			// Protected routes under /api/v1/admin
			r.Route("/admin", func(r chi.Router) {
				r.Use(custommiddleware.RequireAPIKey)
				r.Post("/posts", postHandler.CreatePost)       // POST /api/v1/admin/posts
				r.Put("/posts/{slug}", postHandler.UpdatePost) // PUT /api/v1/admin/posts/{slug}
			})
		})
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Configure server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
