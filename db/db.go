package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes the SQLite database connection and runs migrations.
// Database is stored at ~/.work_manager/work_manager.db — persists across rebuilds.
func InitDB() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	dbDir := filepath.Join(homeDir, ".work_manager")
	dbPath := filepath.Join(dbDir, "work_manager.db")

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	var openErr error
	DB, openErr = sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if openErr != nil {
		return fmt.Errorf("failed to open database: %w", openErr)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(1) // SQLite works best with a single writer
	DB.SetMaxIdleConns(1)

	// Verify connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(DB); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Printf("Database initialized at: %s\n", dbPath)
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
