package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

// AppDatabase is the high-level interface for the DB.
type AppDatabase interface {
	CreateUser(username string) (string, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUserName(userID, newName string) error
	UpdateUserPhoto(userID, photoUrl string) error
	GetConversationsByUserID(userID string) ([]Conversation, error)
	Ping() error
}

// appdbimpl is the concrete implementation of AppDatabase.
type appdbimpl struct {
	db *sql.DB
}

// User represents a user record.
type User struct {
	ID       string
	Username string
	PhotoURL string // Optional profile photo URL.
}

// Conversation represents a conversation record.
type Conversation struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsGroup   bool   `json:"is_group"`
	CreatedAt string `json:"created_at"` // Added field to store the creation timestamp
}

// New creates a new database instance.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	// Create users table if not exists, including an optional photo_url column.
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		photo_url TEXT
	)`)
	if err != nil {
		return nil, fmt.Errorf("error creating users table: %w", err)
	}

	// Create conversations table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    name TEXT, -- Name of group (NULL for private chats)
    is_group BOOLEAN NOT NULL DEFAULT 0, -- 0 = Private Chat, 1 = Group Chat
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)`)
	if err != nil {
		return nil, fmt.Errorf("error creating conversations table: %w", err)
	}
	//create messages table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
	id TEXT PRIMARY KEY,
	conversation_id TEXT NOT NULL,
	sender_id TEXT NOT NULL,
	content TEXT NOT NULL, -- Message text or media URL
	reply_to TEXT NULL, -- If replying to another message
	sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (conversation_id) REFERENCES conversations(id),
	FOREIGN KEY (sender_id) REFERENCES users(id),
	FOREIGN KEY (reply_to) REFERENCES messages(id) ON DELETE CASCADE
)`)
	if err != nil {
		return nil, fmt.Errorf("error creating messages table: %w", err)
	}
	//create group_members table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS group_members (
    group_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	
)`)
	if err != nil {
		return nil, fmt.Errorf("error creating groups table: %w", err)
	}
	//create message reactions table if not exists

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS message_reactions (
    message_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    reaction TEXT NOT NULL, -- Example: "😂" or "🔥"
    reacted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	
)`)
	if err != nil {
		return nil, fmt.Errorf("error creating message_reactions table: %w", err)
	}
	return &appdbimpl{db: db}, nil
}

// CreateUser inserts a new user.
func (db *appdbimpl) CreateUser(username string) (string, error) {
	// Generate a new UUID for the user.
	userID, err := GenerateNewID()
	if err != nil {
		return "", fmt.Errorf("failed to create id: %w", err)
	}

	// Insert the user into the database.
	_, err = db.db.Exec("INSERT INTO users (id, username) VALUES (?, ?)", userID, username)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

// GetUserByUsername retrieves a user by username.
func (db *appdbimpl) GetUserByUsername(username string) (*User, error) {
	row := db.db.QueryRow("SELECT id, username, photo_url FROM users WHERE username = ?", username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PhotoURL)
	if err == sql.ErrNoRows {
		return nil, nil // User does not exist.
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &user, nil
}

// UpdateUserName changes the username for the specified user ID without changing the ID.
func (db *appdbimpl) UpdateUserName(userID, newName string) error {
	// Check if the new username is already taken.
	var existingID string
	err := db.db.QueryRow("SELECT id FROM users WHERE username = ?", newName).Scan(&existingID)
	if err != sql.ErrNoRows {
		if err == nil {
			// Username already exists.
			return errors.New("username already exists")
		}
		// Some other error occurred.
		return fmt.Errorf("error checking for existing username: %w", err)
	}

	// Update the username.
	_, err = db.db.Exec("UPDATE users SET username = ? WHERE id = ?", newName, userID)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

// UpdateUserPhoto updates the profile photo URL for the specified user ID.
func (db *appdbimpl) UpdateUserPhoto(userID, photoUrl string) error {
	_, err := db.db.Exec("UPDATE users SET photo_url = ? WHERE id = ?", photoUrl, userID)

	if err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}
	return nil
}

// Ping checks the database connection.
func (db *appdbimpl) Ping() error {
	return db.db.Ping()
}

// GenerateNewID generates a new UUID.
func GenerateNewID() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

// GetConversationsByUserID retrieves all conversations associated with a user.
func (db *appdbimpl) GetConversationsByUserID(userID string) ([]Conversation, error) {
	// Updated query to include the created_at column
	rows, err := db.db.Query(`
		SELECT c.id, c.name, c.is_group, c.created_at
		FROM conversations c
		JOIN group_members gm ON c.id = gm.group_id
		WHERE gm.user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch conversations: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conversation Conversation
		if err := rows.Scan(&conversation.ID, &conversation.Name, &conversation.IsGroup, &conversation.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conversation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return conversations, nil
}
