// internal/handler/post.go
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pmk808/blog-api/internal/model"
	"github.com/pmk808/blog-api/internal/storage"
)

type PostHandler struct {
	store *storage.PostStore
}

func NewPostHandler(store *storage.PostStore) *PostHandler {
	return &PostHandler{
		store: store,
	}
}

func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	posts, err := h.store.ListPosts(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	respondJSON(w, posts, http.StatusOK)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		respondError(w, "Invalid slug", http.StatusBadRequest)
		return
	}

	var input model.PostUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.store.UpdatePost(r.Context(), slug, &input)
	if err != nil {
		respondError(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	if post == nil {
		respondError(w, "Post not found", http.StatusNotFound)
		return
	}

	respondJSON(w, post, http.StatusOK)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var input model.PostCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if err := validatePost(&input); err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.store.CreatePost(r.Context(), &input)
	if err != nil {
		respondError(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	respondJSON(w, post, http.StatusCreated)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		respondError(w, "Invalid slug", http.StatusBadRequest)
		return
	}

	post, err := h.store.GetPost(r.Context(), slug)
	if err != nil {
		respondError(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	if post == nil {
		respondError(w, "Post not found", http.StatusNotFound)
		return
	}

	respondJSON(w, post, http.StatusOK)
}

// validatePost performs basic validation on post input
func validatePost(post *model.PostCreate) error {
	if post.Title == "" {
		return fmt.Errorf("title is required")
	}

	if post.Intro.Question == "" || post.Intro.Hook == "" {
		return fmt.Errorf("intro question and hook are required")
	}

	if len(post.Summary.Points) == 0 {
		return fmt.Errorf("at least one TLDR point is required")
	}

	if len(post.Content.Sections) == 0 {
		return fmt.Errorf("at least one content section is required")
	}

	if len(post.Impact.Points) == 0 {
		return fmt.Errorf("at least one impact point is required")
	}

	if len(post.Insights.Points) == 0 {
		return fmt.Errorf("at least one insight point is required")
	}

	return nil
}

// Helper functions for responses
func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Error: message,
	})
}
