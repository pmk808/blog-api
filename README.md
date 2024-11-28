Blog API Project Overview:

Structure:
```
blog-api/
├── cmd/
│   └── server/
│       └── main.go           // Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/         // HTTP request handlers
│   │   └── middleware/       // Auth, logging middleware
│   ├── models/              // Database models
│   ├── repository/          // Database operations
│   └── service/             // Business logic
├── pkg/
│   ├── database/            // DB connection, migrations 
│   └── utils/               // Helper functions
├── config/                  // Configuration files
└── migrations/              // SQL migrations
```

Features:
1. Blog Post Management
   - CRUD operations
   - Post status (draft/published)
   - Tags/categories
   - Search functionality

2. Authentication
   - JWT-based auth
   - Admin/author roles
   - Protected routes

3. Database
   - PostgreSQL
   - GORM for ORM
   - Migrations support

4. API Endpoints:
```
GET    /api/posts       // List posts
GET    /api/posts/:id   // Single post
POST   /api/posts       // Create post
PUT    /api/posts/:id   // Update post
DELETE /api/posts/:id   // Delete post
GET    /api/tags        // List tags
POST   /api/auth/login  // Login
```

Tech Stack:
- Go 1.21+
- Fiber framework
- GORM
- PostgreSQL
- JWT

Deployment:
- Fly.io for hosting
- PostgreSQL on Fly.io
- CI/CD with GitHub Actions