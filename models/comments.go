package models

import (
	"time"
)

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	PostID    uint   `gorm:"index;not null"`
	Author    string `gorm:"size:100;not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
}
