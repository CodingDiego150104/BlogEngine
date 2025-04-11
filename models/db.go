package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("[error] failed to initialize database, got error %v", err)
	}

	// Auto migrate Post struct
	err = DB.AutoMigrate(&Post{})
	if err != nil {
		log.Fatalf("[error] failed to migrate Post model, got error %v", err)
	}
}
