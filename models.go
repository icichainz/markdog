package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // SQLite driver
)

// DB is the global database connection
var DB *gorm.DB

// Document represents the document model
type Document struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `gorm:"null"`
	Content string `gorm:"type:text"`
	UserID  uint   `gorm:"not null"` 
	HashedContent string `gorm:"not null"`
}

// User model
type User struct {
	ID        uint       `gorm:"primaryKey"`
	Username  string     `gorm:"not null;unique"`           // Ensure username is unique
	LastLogin time.Time  `gorm:"default:CURRENT_TIMESTAMP"` // Last login time
	IsActive  bool       `gorm:"default:true"`              // Is active flag
	Documents []Document `gorm:"foreignKey:UserID"`         // One-to-many relationship
}


// Initialize initializes the SQLite database and sets up the required tables
func DBInitialize() error {
	var err error
	DB, err = gorm.Open("sqlite3", "app.db")
	if err != nil {
		return err // Return the error instead of logging and terminating
	}

	// Migrate the schema, this will create the documents table if it doesn't exist
	if err := DB.AutoMigrate(&Document{}).Error; err != nil {
		return err
	}

	if err := DB.AutoMigrate(&User{}).Error; err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// Close closes the database connection
func DBClose() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}
}
