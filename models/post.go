package main

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	Author    string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
