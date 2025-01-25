# Blog API

A lightweight blog API backend built with Go and PostgreSQL.

## Project Structure
```
blog-api/
├── internal/
│   ├── handler/            
│   │   └── post.go         # Post handlers (CRUD operations)
│   ├── middleware/         
│   │   └── auth.go         # API key authentication
│   ├── model/              
│   │   └── post.go         # Post data model
│   └── storage/            
│       └── postgres.go     # Database connection and setup
└── main.go                 # Application entry point
```

## Features
- CRUD operations for blog posts
- API key authentication for protected routes
- PostgreSQL database integration
- CORS support

## API Endpoints
```
GET    /posts          # List all posts
GET    /posts/:slug    # Get post by slug
POST   /posts          # Create post (protected)
PUT    /posts/:slug    # Update post (protected)
DELETE /posts/:slug    # Delete post (protected)
```

## Tech Stack
- Go with Gin framework
- PostgreSQL with GORM
- UUID for post IDs
- API key authentication

## Getting Started

1. Prerequisites
   - Go 1.21+
   - PostgreSQL

2. Environment Variables
   ```
   DB_CONN=postgresql://[user]:[password]@[host]:[port]/[dbname]
   API_KEY=[your-api-key]
   ```

3. Running Locally
   ```bash
   go run main.go
   ```

## Deployment
Currently deployed on Render.com with PostgreSQL database support.

## CORS Configuration
Configured to allow requests from:
- `http://localhost:5174`
- `https://portfolio-mc-dev.vercel.app`