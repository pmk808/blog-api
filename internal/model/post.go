// internal/model/post.go
package model

import (
	"time"
)

type Post struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Description string     `json:"description,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Tags        []string   `json:"tags"`
	IsPublished bool       `json:"is_published"`
}

type PostCreate struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"is_published"`
}

type PostUpdate struct {
	Title       *string   `json:"title"`
	Content     *string   `json:"content"`
	Description *string   `json:"description"`
	Tags        *[]string `json:"tags"`
	IsPublished *bool     `json:"is_published"`
}
