/*
Package database is the middleware between the app database and the code.
All data (de)serialization (save/load) from a persistent database are handled here.
Database-specific logic should never escape this package.

This version supports the WASAText project user management. It provides functions to:
  - Create a new user if one does not exist (or return an existing user’s identifier)
  - Update a user’s username
  - Retrieve a user’s username

This implementation assumes an SQLite database. If you are using MySQL, adjust the table
existence query and SQL dialect accordingly.

Note: The Go version for this project is 1.17.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// AppDatabase is the high-level interface for the WASAText database.
// It currently supports basic user operations.
type AppDatabase interface {
	// CreateOrGetUser searches for a user by name. If the user exists, it returns the user’s identifier.
	// Otherwise, it creates a new user and returns the new identifier.
	CreateOrGetUser(name string) (string, error)
	// UpdateUserName changes the username for a given user identifier.
	UpdateUserName(userID string, newName string) error
	// GetUserName returns the username for a given user identifier.
	GetUserName(userID string) (string, error)

	// Ping checks if the database connection is alive.
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the provided database connection.
// It also applies a migration: if the "users" table does not exist, it creates it.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	// Check if the "users" table exists.
	// For SQLite; for MySQL, you might run "SHOW TABLES LIKE 'users'" instead.
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// The "users" table does not exist; create it.
		sqlStmt := `
		CREATE TABLE users (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating users table: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("error checking for users table: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

// Ping checks if the database connection is still alive.
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// CreateOrGetUser returns the existing identifier for a user with the given name.
// If no such user exists, a new user is created and its identifier is returned.
func (db *appdbimpl) CreateOrGetUser(name string) (string, error) {
	var id string
	err := db.c.QueryRow("SELECT id FROM users WHERE name = ?", name).Scan(&id)
	if err == nil {
		// User already exists.
		return id, nil
	} else if err != sql.ErrNoRows {
		return "", fmt.Errorf("error querying user: %w", err)
	}

	// User does not exist; create a new one.
	id = uuid.New().String()
	_, err = db.c.Exec("INSERT INTO users (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		return "", fmt.Errorf("error inserting new user: %w", err)
	}
	return id, nil
}

// UpdateUserName updates the username for the user with the given identifier.
func (db *appdbimpl) UpdateUserName(userID string, newName string) error {
	res, err := db.c.Exec("UPDATE users SET name = ? WHERE id = ?", newName, userID)
	if err != nil {
		return fmt.Errorf("error updating username: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("no user found with id %s", userID)
	}
	return nil
}

// GetUserName retrieves the username for the user with the specified identifier.
func (db *appdbimpl) GetUserName(userID string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("error retrieving username: %w", err)
	}
	return name, nil
}
