package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes and returns a connection to the SQLite database.
func InitDB() *sql.DB {
	// Set up logging
	logPath := filepath.Join(".", "app.log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)

	// Define database path and open connection
	dbPath := filepath.Join(".", "sqlite.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Check database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create tables if they do not exist
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(100) NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS routers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL,
			password TEXT
		);
	`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Ensure default admin user exists
	createDefaultUser(db)

	log.Println("Database initialized successfully")
	return db
}

// createDefaultUser ensures the default admin user exists
func createDefaultUser(db *sql.DB) {
	// Check if user already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check existing users: %v", err)
	}

	// If no admin user exists, create one
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash default password: %v", err)
		}

		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", "admin", string(hashedPassword))
		if err != nil {
			log.Fatalf("Failed to create default admin user: %v", err)
		}

		log.Println("Default admin user created")
	}
}
