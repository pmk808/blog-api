// internal/model/post.go
package model

import (
    "time"
)

// Post represents the main blog post structure
type Post struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Slug        string     `json:"slug"`
    Intro       IntroSection    `json:"intro"`      // WTF section
    Summary     TLDRSection    `json:"summary"`     // TLDR section
    Content     DeepDiveSection `json:"content"`    // Deep Dive section
    Impact      ImpactSection   `json:"impact"`     // Why Should I Care section
    Insights    InsightsSection `json:"insights"`   // Mind = Blown section
    Resources   []Resource      `json:"resources"`  // Learn More section
    IsPublished bool           `json:"is_published"`
    PublishedAt *time.Time     `json:"published_at,omitempty"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    Tags        []string       `json:"tags"`
}

// IntroSection represents the "WTF?" section
type IntroSection struct {
    Question string `json:"question"`
    Hook     string `json:"hook"`
}

// TLDRSection represents the "TLDR" section
type TLDRSection struct {
    Points []string `json:"points"`
}

// DeepDiveSection represents the "Deep Dive" section
type DeepDiveSection struct {
    Sections []ContentSection `json:"sections"`
}

// ContentSection represents a subsection in the Deep Dive
type ContentSection struct {
    Title    string   `json:"title"`
    Content  string   `json:"content"`
    Points   []string `json:"points"`
    Examples []string `json:"examples,omitempty"`
}

// ImpactSection represents the "Why Should I Care?" section
type ImpactSection struct {
    Points []string `json:"points"`
}

// InsightsSection represents the "Mind = Blown" section
type InsightsSection struct {
    Points []string `json:"points"`
}

// Resource represents an item in the "Learn More" section
type Resource struct {
    Title string `json:"title"`
    URL   string `json:"url"`
    Type  string `json:"type"` // "documentation", "research", "tutorial", etc.
}

// PostCreate represents the structure for creating a new post
type PostCreate struct {
    Title       string          `json:"title"`
    Intro       IntroSection    `json:"intro"`
    Summary     TLDRSection    `json:"summary"`
    Content     DeepDiveSection `json:"content"`
    Impact      ImpactSection   `json:"impact"`
    Insights    InsightsSection `json:"insights"`
    Resources   []Resource      `json:"resources"`
    Tags        []string        `json:"tags"`
    IsPublished bool           `json:"is_published"`
}

// PostUpdate represents the structure for updating an existing post
type PostUpdate struct {
    Title       *string          `json:"title,omitempty"`
    Intro       *IntroSection    `json:"intro,omitempty"`
    Summary     *TLDRSection    `json:"summary,omitempty"`
    Content     *DeepDiveSection `json:"content,omitempty"`
    Impact      *ImpactSection   `json:"impact,omitempty"`
    Insights    *InsightsSection `json:"insights,omitempty"`
    Resources   *[]Resource      `json:"resources,omitempty"`
    Tags        *[]string        `json:"tags,omitempty"`
    IsPublished *bool           `json:"is_published,omitempty"`
}