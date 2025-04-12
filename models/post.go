package models

import (
	"time"
)

type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:255;not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	Comments  []Comment `gorm:"foreignKey:PostID"`
}
