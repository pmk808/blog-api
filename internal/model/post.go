package model

import "time"

type Post struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug" gorm:"uniqueIndex"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
