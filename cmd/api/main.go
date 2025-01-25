package main

import (
	"os"

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
