# Blog API Project Overview

A lightweight, efficient blog API backend with built-in content management capabilities.

## Project Structure
```
blog-api/
├── cmd/
│   └── api/
│       └── main.go           // Application entry point
├── internal/
│   ├── handler/             // HTTP handlers
│   │   ├── post.go
│   │   └── create.go        // Post creation handler
│   ├── storage/             // Database operations
│   │   ├── post.go
│   │   └── db.go           // Database connection and queries
│   └── model/              // Data models
│       └── post.go         // Post model definition
├── migrations/             // SQL migrations
└── ui/                     // Content management interface
    └── index.html         // Markdown editor
```

## Features

1. Blog Post Management
   - Create, read, update, delete operations
   - Draft/published status support
   - Tag support
   - Full-text search capability
   - Markdown content support
   - Rich text preview

2. API Security
   - Secure endpoints for content management
   - Protected routes for modifications
   - Environment-based configuration

3. Database
   - PostgreSQL for reliable data storage
   - SQL migrations support
   - Efficient querying and indexing

4. API Endpoints
```
GET    /api/posts       // List posts
GET    /api/posts/:slug // Get single post
POST   /create         // Create post
PUT    /posts/:slug    // Update post
DELETE /posts/:slug    // Delete post
GET    /api/tags       // List all tags
```

## Tech Stack
- Go 1.21+
- PostgreSQL
- Chi router
- Native `database/sql`
- Server-side Markdown parsing

## Content Management
- Built-in Markdown editor
- Real-time preview
- Tag management
- Draft/publish functionality
- Direct content posting

## Deployment
- Fly.io for hosting
- PostgreSQL on Fly.io
- GitHub Actions for CI/CD

## Getting Started

1. Prerequisites
   - Go 1.21 or higher
   - PostgreSQL 14+
   - Make (optional, for using Makefile commands)

2. Configuration
   ```bash
   # Set required environment variables
   export DATABASE_URL=postgresql://user:pass@localhost:5432/blogdb
   export BLOG_API_KEY=your-secret-key
   ```

3. Running Locally
   ```bash
   # Start the server
   go run cmd/api/main.go
   ```

4. Development
   ```bash
   # Run migrations
   make migrate-up

   # Run tests
   make test
   ```