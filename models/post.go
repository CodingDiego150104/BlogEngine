package main

import (
	"time"
)

type Post struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `validate:"required"`
	Content string `validate:"required"`
	Date    time.Time
}
