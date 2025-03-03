package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB inicializa e retorna a conexão com o banco de dados SQLite
func InitDB() *sql.DB {
	// Define o caminho absoluto para o arquivo de log
	logPath := filepath.Join(".", "app.log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close() // Fecha o arquivo no final da execução

	log.SetOutput(logFile)

	// Define o caminho absoluto para o banco de dados SQLite
	dbPath := filepath.Join(".", "sqlite.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Verifica se a conexão foi estabelecida corretamente
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

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

	log.Println("Database initialized successfully")
	return db
}
