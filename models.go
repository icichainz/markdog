package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// DB is the global database connection
var DB *sql.DB

// Initialize initializes the SQLite database and sets up the required tables
func Initialize() {
	var err error
	DB, err = sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create table if it doesn't exist
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS documents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

// Close closes the database connection
func Close() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}
}
