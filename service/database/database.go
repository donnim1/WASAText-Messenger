package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(username string) (string, error)
	Ping() error
}

type appdbimpl struct {
	db *sql.DB  // Changed from 'c' to 'db' for clarity
}

// New creates a new database instance
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	// Create users table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL
	)`)
	if err != nil {
		return nil, fmt.Errorf("error creating users table: %w", err)
	}

	return &appdbimpl{db: db}, nil
}

func (db *appdbimpl) CreateUser(username string) (string, error) {
	// Implementation from previous steps
	// This is just a placeholder - use your actual UUID generation
	// and database operations here
	return "generated-uuid", nil
}

func (db *appdbimpl) Ping() error {
	return db.db.Ping()
}