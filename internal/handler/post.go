package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmk808/blog-api/internal/model"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

// NewPostHandler creates a new PostHandler with the provided database connection
func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{
		db: db,
	}
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	var posts []model.Post
	if result := h.db.Find(&posts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var post model.Post

	if result := h.db.Where("slug = ?", slug).First(&post); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// CreatePost handles new post creation
func (h *PostHandler) CreatePost(c *gin.Context) {
	var newPost model.Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if newPost.Title == "" || newPost.Slug == "" || newPost.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Check for existing slug
	var existingPost model.Post
	if result := h.db.Where("slug = ?", newPost.Slug).First(&existingPost); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Slug already exists"})
		return
	}

	if result := h.db.Create(&newPost); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPost)
}

// UpdatePost updates an existing post
func (h *PostHandler) UpdatePost(c *gin.Context) {
	slug := c.Param("slug")
	var existingPost model.Post

	// Find existing post
	if result := h.db.Where("slug = ?", slug).First(&existingPost); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Bind update data
	var updateData model.Post
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent slug changes
	if updateData.Slug != "" && updateData.Slug != existingPost.Slug {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot change slug"})
		return
	}

	// Update fields
	existingPost.Title = updateData.Title
	existingPost.Content = updateData.Content

	if result := h.db.Save(&existingPost); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, existingPost)
}

// DeletePost deletes a post by slug
func (h *PostHandler) DeletePost(c *gin.Context) {
	slug := c.Param("slug")
	var post model.Post

	// Find post first to check existence
	if result := h.db.Where("slug = ?", slug).First(&post); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result := h.db.Delete(&post); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
