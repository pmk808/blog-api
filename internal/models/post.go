package models

import (
	"context"
	"database/sql"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) UpsertPost(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (title, slug, content, tags, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (slug) 
		DO UPDATE SET 
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			tags = EXCLUDED.tags,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		post.Title,
		post.Slug,
		post.Content,
		post.Tags,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *PostRepository) GetBySlug(ctx context.Context, slug string) (*Post, error) {
	post := &Post{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, title, slug, content, tags, created_at, updated_at FROM posts WHERE slug = $1",
		slug,
	).Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return post, err
}

func (r *PostRepository) List(ctx context.Context) ([]*Post, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, title, slug, content, tags, created_at, updated_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}
