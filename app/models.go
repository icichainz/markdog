package app

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // SQLite driver
	"github.com/google/uuid"
)

// DB is the global database connection
var DB *gorm.DB

// Document represents the document model
type Document struct {
	ID            uint      `gorm:"primaryKey"`
	UId           string    `gorm:"type:varchar(36);uniqueIndex;not null"`
	Title         string    `gorm:"not null;default:'Untitled'"`  // Default title
	Content       string    `gorm:"type:text"`
	HashedContent string    `gorm:"not null"`
	UserID        uint      `gorm:"not null;index"`              // Added index for better query performance
	User          User      `gorm:"foreignkey:UserID"`           // Add reference to User
	IsPublic      bool      `gorm:"default:false"`              // Privacy setting
	Version       int       `gorm:"default:1"`                   // Document version
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `gorm:"index"`                     // Soft delete support
}

// User model
type User struct {
	ID            uint       `gorm:"primaryKey"`
	UId           string     `gorm:"type:varchar(36);uniqueIndex;not null"`
	Username      string     `gorm:"not null;unique;index"`
	CreatedAt      time.Time 
	UpdatedAt      time.Time        
}

// BeforeCreate is a GORM hook that sets the UId before creating a new document
func (d *Document) BeforeCreate(scope *gorm.Scope) error {	
	if d.UId == "" {
		d.UId = uuid.New().String() // Generate a new UUID if UId is not set
	}
	return nil
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
