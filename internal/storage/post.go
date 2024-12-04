// internal/storage/post.go
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/pmk808/blog-api/internal/model"
)

func createSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove all characters except letters, numbers, and hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

type PostStore struct {
	db *DB
}

func NewPostStore(db *DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

// CreatePost creates a new blog post with all its sections
func (s *PostStore) CreatePost(ctx context.Context, post *model.PostCreate) (*model.Post, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert main post
	var createdPost model.Post
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO posts (
            slug, title, intro_question, intro_hook, 
            tldr_points, impact_points, insight_points,
            tags, is_published
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, slug, title, created_at, updated_at`,
		createSlug(post.Title),
		post.Title,
		post.Intro.Question,
		post.Intro.Hook,
		post.Summary.Points,
		post.Impact.Points,
		post.Insights.Points,
		post.Tags,
		post.IsPublished,
	).Scan(
		&createdPost.ID,
		&createdPost.Slug,
		&createdPost.Title,
		&createdPost.CreatedAt,
		&createdPost.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error inserting post: %w", err)
	}

	// Insert deep dive sections
	for i, section := range post.Content.Sections {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO content_sections (
                post_id, title, content, points, examples, display_order
            ) VALUES ($1, $2, $3, $4, $5, $6)`,
			createdPost.ID,
			section.Title,
			section.Content,
			section.Points,
			section.Examples,
			i,
		)
		if err != nil {
			return nil, fmt.Errorf("error inserting content section: %w", err)
		}
	}

	// Insert resources
	for _, resource := range post.Resources {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO resources (post_id, title, url, type)
            VALUES ($1, $2, $3, $4)`,
			createdPost.ID,
			resource.Title,
			resource.URL,
			resource.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("error inserting resource: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return s.GetPost(ctx, createdPost.Slug)
}

// GetPost retrieves a post with all its sections by slug
func (s *PostStore) GetPost(ctx context.Context, slug string) (*model.Post, error) {
	// Get main post data
	var post model.Post
	err := s.db.db.QueryRowContext(
		ctx,
		`SELECT 
            id, slug, title, intro_question, intro_hook,
            tldr_points, impact_points, insight_points,
            tags, is_published, created_at, updated_at,
            published_at
        FROM posts WHERE slug = $1`,
		slug,
	).Scan(
		&post.ID,
		&post.Slug,
		&post.Title,
		&post.Intro.Question,
		&post.Intro.Hook,
		&post.Summary.Points,
		&post.Impact.Points,
		&post.Insights.Points,
		&post.Tags,
		&post.IsPublished,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.PublishedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	}

	// Get content sections
	rows, err := s.db.db.QueryContext(
		ctx,
		`SELECT title, content, points, examples
        FROM content_sections
        WHERE post_id = $1
        ORDER BY display_order`,
		post.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting content sections: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var section model.ContentSection
		if err := rows.Scan(
			&section.Title,
			&section.Content,
			&section.Points,
			&section.Examples,
		); err != nil {
			return nil, fmt.Errorf("error scanning content section: %w", err)
		}
		post.Content.Sections = append(post.Content.Sections, section)
	}

	// Get resources
	rows, err = s.db.db.QueryContext(
		ctx,
		`SELECT title, url, type
        FROM resources
        WHERE post_id = $1`,
		post.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting resources: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var resource model.Resource
		if err := rows.Scan(
			&resource.Title,
			&resource.URL,
			&resource.Type,
		); err != nil {
			return nil, fmt.Errorf("error scanning resource: %w", err)
		}
		post.Resources = append(post.Resources, resource)
	}

	return &post, nil
}

func (s *PostStore) ListPosts(ctx context.Context, page, pageSize int) ([]*model.Post, error) {
	offset := (page - 1) * pageSize

	// First get the basic post data
	query := `
        SELECT 
            id, slug, title, intro_question, intro_hook,
            tldr_points, impact_points, insight_points,
            tags, is_published, created_at, updated_at,
            published_at
        FROM posts
        WHERE is_published = true
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2`

	rows, err := s.db.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing posts: %w", err)
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.Slug,
			&post.Title,
			&post.Intro.Question,
			&post.Intro.Hook,
			&post.Summary.Points,
			&post.Impact.Points,
			&post.Insights.Points,
			&post.Tags,
			&post.IsPublished,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.PublishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}

		// Get content sections for each post
		sections, err := s.getContentSections(ctx, post.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting content sections: %w", err)
		}
		post.Content.Sections = sections

		// Get resources for each post
		resources, err := s.getResources(ctx, post.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting resources: %w", err)
		}
		post.Resources = resources

		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *PostStore) UpdatePost(ctx context.Context, slug string, update *model.PostUpdate) (*model.Post, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// First check if post exists
	var postID string
	err = tx.QueryRowContext(ctx, "SELECT id FROM posts WHERE slug = $1", slug).Scan(&postID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error checking post existence: %w", err)
	}

	// Update main post fields
	query := `UPDATE posts SET
        title = COALESCE($1, title),
        intro_question = COALESCE($2, intro_question),
        intro_hook = COALESCE($3, intro_hook),
        tldr_points = COALESCE($4, tldr_points),
        impact_points = COALESCE($5, impact_points),
        insight_points = COALESCE($6, insight_points),
        tags = COALESCE($7, tags),
        is_published = COALESCE($8, is_published),
        updated_at = CURRENT_TIMESTAMP
        WHERE slug = $9`

	_, err = tx.ExecContext(ctx, query,
		getNullableValue(update.Title),
		getNullableValue(update.Intro.Question),
		getNullableValue(update.Intro.Hook),
		getNullableValue(update.Summary.Points),
		getNullableValue(update.Impact.Points),
		getNullableValue(update.Insights.Points),
		getNullableValue(update.Tags),
		getNullableValue(update.IsPublished),
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating post: %w", err)
	}

	// Update content sections if provided
	if update.Content != nil {
		// Delete existing sections
		_, err = tx.ExecContext(ctx, "DELETE FROM content_sections WHERE post_id = $1", postID)
		if err != nil {
			return nil, fmt.Errorf("error deleting existing sections: %w", err)
		}

		// Insert new sections
		for i, section := range update.Content.Sections {
			_, err = tx.ExecContext(ctx,
				`INSERT INTO content_sections (
                    post_id, title, content, points, examples, display_order
                ) VALUES ($1, $2, $3, $4, $5, $6)`,
				postID,
				section.Title,
				section.Content,
				section.Points,
				section.Examples,
				i,
			)
			if err != nil {
				return nil, fmt.Errorf("error inserting content section: %w", err)
			}
		}
	}

	// Update resources if provided
	if update.Resources != nil {
		// Delete existing resources
		_, err = tx.ExecContext(ctx, "DELETE FROM resources WHERE post_id = $1", postID)
		if err != nil {
			return nil, fmt.Errorf("error deleting existing resources: %w", err)
		}

		// Insert new resources
		for _, resource := range *update.Resources {
			_, err = tx.ExecContext(ctx,
				`INSERT INTO resources (post_id, title, url, type)
                VALUES ($1, $2, $3, $4)`,
				postID,
				resource.Title,
				resource.URL,
				resource.Type,
			)
			if err != nil {
				return nil, fmt.Errorf("error inserting resource: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return s.GetPost(ctx, slug)
}

// Helper functions
func getNullableValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	return v
}

// Helper method to get content sections
func (s *PostStore) getContentSections(ctx context.Context, postID string) ([]model.ContentSection, error) {
	rows, err := s.db.db.QueryContext(ctx,
		`SELECT title, content, points, examples
        FROM content_sections
        WHERE post_id = $1
        ORDER BY display_order`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []model.ContentSection
	for rows.Next() {
		var section model.ContentSection
		if err := rows.Scan(
			&section.Title,
			&section.Content,
			&section.Points,
			&section.Examples,
		); err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}
	return sections, nil
}

// Helper method to get resources
func (s *PostStore) getResources(ctx context.Context, postID string) ([]model.Resource, error) {
	rows, err := s.db.db.QueryContext(ctx,
		`SELECT title, url, type
        FROM resources
        WHERE post_id = $1`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []model.Resource
	for rows.Next() {
		var resource model.Resource
		if err := rows.Scan(
			&resource.Title,
			&resource.URL,
			&resource.Type,
		); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}
