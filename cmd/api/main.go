package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pmk808/blog-api/internal/handler"
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

	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:slug", postHandler.GetPostBySlug)
	r.POST("/posts", postHandler.CreatePost)
	r.PUT("/posts/:slug", postHandler.UpdatePost)
	r.DELETE("/posts/:slug", postHandler.DeletePost)

	r.Run() // Listen on 0.0.0.0:8080
}
