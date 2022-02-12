package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDb() {
	log.Println("Connecting DB...")

	db, err := gorm.Open(sqlite.Open("grpcdemo.db"), &gorm.Config{})

	if err != nil {
		panic("Connot connect to db")
	}

	DB = db
}
