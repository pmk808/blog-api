package util

import (
	"errors"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

const (
	MaxPostSize = 1024 * 1024 // 1MB
)

// SanitizeMarkdown sanitizes markdown content to prevent XSS attacks
func SanitizeMarkdown(content string) string {
	p := bluemonday.UGCPolicy()
	// Allow certain HTML elements commonly used in markdown
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6", "p", "br", "strong", "em", "code", "pre")
	// Allow certain attributes
	p.AllowAttrs("class").OnElements("code", "pre")

	return p.Sanitize(content)
}

// ValidatePost validates post content
func ValidatePost(title, content, slug string) error {
	if len(content) > MaxPostSize {
		return errors.New("post content exceeds maximum size limit")
	}

	if strings.TrimSpace(title) == "" {
		return errors.New("title cannot be empty")
	}

	if strings.TrimSpace(slug) == "" {
		return errors.New("slug cannot be empty")
	}

	// Validate slug format (alphanumeric and hyphens only)
	for _, char := range slug {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return errors.New("slug can only contain lowercase letters, numbers, and hyphens")
		}
	}

	return nil
}
