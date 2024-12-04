-- migrations/002_structured_content.sql

-- Modify existing posts table
ALTER TABLE posts
ADD COLUMN intro_question TEXT,
ADD COLUMN intro_hook TEXT,
ADD COLUMN tldr_points TEXT[],
ADD COLUMN impact_points TEXT[],
ADD COLUMN insight_points TEXT[];

-- Create table for deep dive sections
CREATE TABLE content_sections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    points TEXT[],
    examples TEXT[],
    display_order INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create table for resources
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_content_sections_post_id ON content_sections(post_id);
CREATE INDEX idx_content_sections_order ON content_sections(post_id, display_order);
CREATE INDEX idx_resources_post_id ON resources(post_id);

-- Update updated_at trigger for content_sections
CREATE TRIGGER update_content_sections_updated_at
    BEFORE UPDATE ON content_sections
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();