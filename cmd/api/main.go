package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pmk808/blog-api/internal/handler"
	"github.com/pmk808/blog-api/internal/middleware"
	"github.com/pmk808/blog-api/internal/storage"
)

func main() {
	// Get DB connection string from env
	dbConn := os.Getenv("DB_CONN")
	db, err := storage.NewDB(dbConn)
	if err != nil {
		panic("failed to connect to database")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://your-vercel-app.vercel.app", // Your production domain
			"http://localhost:3000",              // Local development
		},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	postHandler := handler.NewPostHandler(db)

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY environment variable not set")
	}

	// Public routes
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:slug", postHandler.GetPostBySlug)

	// Protected routes group
	authorized := r.Group("/")
	authorized.Use(middleware.APIKeyAuth(apiKey))
	{
		authorized.POST("/posts", postHandler.CreatePost)
		authorized.PUT("/posts/:slug", postHandler.UpdatePost)
		authorized.DELETE("/posts/:slug", postHandler.DeletePost)
	}

	r.Run() // Listen on 0.0.0.0:8080
}
