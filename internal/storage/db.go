// internal/storage/db.go
package storage

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/lib/pq"
)

type DB struct {
    db *sql.DB
}

func NewDB() (*DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
    }

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

    return &DB{db: db}, nil
}

func (db *DB) Close() error {
    return db.db.Close()
}

// GetDB returns the underlying sql.DB instance
func (db *DB) GetDB() *sql.DB {
    return db.db
}