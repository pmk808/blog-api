// internal/handler/post.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pmk808/blog-api/internal/model"
	"github.com/pmk808/blog-api/internal/storage"
)

type PostHandler struct {
	store *storage.PostStore
}

func NewPostHandler(store *storage.PostStore) *PostHandler {
	return &PostHandler{store: store}
}

// ListPosts handles GET /api/posts
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	// Simple pagination, defaults to first 10 posts
	limit := 10
	offset := 0

	posts, err := h.store.ListPosts(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	respondJSON(w, posts)
}

// GetPost handles GET /api/posts/{slug}
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Invalid slug", http.StatusBadRequest)
		return
	}

	post, err := h.store.GetPost(r.Context(), slug)
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	respondJSON(w, post)
}

// CreatePost handles POST /create
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var input model.PostCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if input.Title == "" || input.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	post, err := h.store.CreatePost(r.Context(), &input)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	respondJSON(w, post)
}

// UpdatePost handles PUT /posts/{slug}
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Invalid slug", http.StatusBadRequest)
		return
	}

	var input model.PostUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.store.UpdatePost(r.Context(), slug, &input)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	respondJSON(w, post)
}

// Helper function
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
