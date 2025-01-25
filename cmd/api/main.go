package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pmk808/blog-api/internal/handler"
	"github.com/pmk808/blog-api/internal/middleware"
	"github.com/pmk808/blog-api/internal/storage"
)

func main() {
	dbConn := os.Getenv("DB_CONN")
	db, err := storage.NewDB(dbConn)
	if err != nil {
		panic("failed to connect to database")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5174",
			"https://portfolio-mc-dev.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	postHandler := handler.NewPostHandler(db)

	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "X-API-Key,Content-Type,Authorization")
		c.Status(http.StatusNoContent)
	})

	// Routes
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:slug", postHandler.GetPostBySlug)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.APIKeyAuth(os.Getenv("API_KEY")))
	{
		authorized.POST("/posts", postHandler.CreatePost)
		authorized.PUT("/posts/:slug", postHandler.UpdatePost)
		authorized.DELETE("/posts/:slug", postHandler.DeletePost)
	}

	r.Run()
}
