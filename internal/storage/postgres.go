package storage

import (
	"github.com/pmk808/blog-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Auto-migrate the Post model
	db.AutoMigrate(&model.Post{})
	return db, nil
}
