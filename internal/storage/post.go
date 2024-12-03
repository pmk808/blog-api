// internal/storage/post.go
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pmk808/blog-api/internal/model"
)

type PostStore struct {
	db *DB
}

func NewPostStore(db *DB) *PostStore {
	return &PostStore{db: db}
}

func (s *PostStore) CreatePost(ctx context.Context, post *model.PostCreate) (*model.Post, error) {
	query := `
        INSERT INTO posts (slug, title, content, description, tags, is_published, published_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, slug, title, content, description, tags, is_published, published_at, created_at, updated_at`

	var publishedAt *time.Time
	if post.IsPublished {
		now := time.Now()
		publishedAt = &now
	}

	// Create slug from title (you'll need to implement this)
	slug := createSlug(post.Title)

	var p model.Post
	err := s.db.db.QueryRowContext(
		ctx,
		query,
		slug,
		post.Title,
		post.Content,
		post.Description,
		post.Tags,
		post.IsPublished,
		publishedAt,
	).Scan(
		&p.ID,
		&p.Slug,
		&p.Title,
		&p.Content,
		&p.Description,
		&p.Tags,
		&p.IsPublished,
		&p.PublishedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PostStore) GetPost(ctx context.Context, slug string) (*model.Post, error) {
	query := `
        SELECT id, slug, title, content, description, tags, is_published, published_at, created_at, updated_at
        FROM posts
        WHERE slug = $1`

	var p model.Post
	err := s.db.db.QueryRowContext(ctx, query, slug).Scan(
		&p.ID,
		&p.Slug,
		&p.Title,
		&p.Content,
		&p.Description,
		&p.Tags,
		&p.IsPublished,
		&p.PublishedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PostStore) ListPosts(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	query := `
        SELECT id, slug, title, content, description, tags, is_published, published_at, created_at, updated_at
        FROM posts
        WHERE is_published = true
        ORDER BY published_at DESC
        LIMIT $1 OFFSET $2`

	rows, err := s.db.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p model.Post
		err := rows.Scan(
			&p.ID,
			&p.Slug,
			&p.Title,
			&p.Content,
			&p.Description,
			&p.Tags,
			&p.IsPublished,
			&p.PublishedAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

func (s *PostStore) UpdatePost(ctx context.Context, slug string, update *model.PostUpdate) (*model.Post, error) {
	// Start building the query dynamically based on which fields are being updated
	query := "UPDATE posts SET"
	args := []interface{}{slug}
	argCount := 1

	// Build the query based on which fields are being updated
	var updates []string
	if update.Title != nil {
		argCount++
		updates = append(updates, fmt.Sprintf(" title = $%d", argCount))
		args = append(args, *update.Title)
	}
	if update.Content != nil {
		argCount++
		updates = append(updates, fmt.Sprintf(" content = $%d", argCount))
		args = append(args, *update.Content)
	}
	if update.Description != nil {
		argCount++
		updates = append(updates, fmt.Sprintf(" description = $%d", argCount))
		args = append(args, *update.Description)
	}
	if update.Tags != nil {
		argCount++
		updates = append(updates, fmt.Sprintf(" tags = $%d", argCount))
		args = append(args, *update.Tags)
	}
	if update.IsPublished != nil {
		argCount++
		updates = append(updates, fmt.Sprintf(" is_published = $%d", argCount))
		args = append(args, *update.IsPublished)

		// If publishing for the first time, set published_at
		if *update.IsPublished {
			argCount++
			updates = append(updates, fmt.Sprintf(" published_at = $%d", argCount))
			now := time.Now()
			args = append(args, now)
		}
	}

	if len(updates) == 0 {
		return s.GetPost(ctx, slug)
	}

	query += strings.Join(updates, ",") + " WHERE slug = $1 RETURNING *"

	var p model.Post
	err := s.db.db.QueryRowContext(ctx, query, args...).Scan(
		&p.ID,
		&p.Slug,
		&p.Title,
		&p.Content,
		&p.Description,
		&p.Tags,
		&p.IsPublished,
		&p.PublishedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Helper function to create URL-friendly slugs
func createSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}
