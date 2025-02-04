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
	GetConversation(conversationID string) (*Conversation, []Message, error)

	SendMessage(senderID string, receiverID string, content string, isGroup bool, groupID string) (string, error)
	ForwardMessage(originalMessageID, targetConversationID, senderID string) (string, error)
	CommentMessage(messageID, userID, reaction string) error
	UncommentMessage(messageID, userID string) error
	DeleteMessage(messageID, senderID string) error

	AddToGroup(groupID, userID string) error
	LeaveGroup(groupID, userID string) error
	SetGroupName(groupID, newName string) error
	SetGroupPhoto(groupID, photoUrl string) error

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

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	Content        string
	ReplyTo        sql.NullString // use NullString if replies are optional
	SentAt         string         // using string for simplicity; you may use time.Time
}

// GetConversation retrieves a conversation and all its messages.
func (db *appdbimpl) GetConversation(conversationID string) (*Conversation, []Message, error) {
	// Retrieve conversation details.
	var conv Conversation
	err := db.db.QueryRow("SELECT id, name, is_group, created_at FROM conversations WHERE id = ?", conversationID).
		Scan(&conv.ID, &conv.Name, &conv.IsGroup, &conv.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil, nil // Conversation not found.
	} else if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve conversation: %w", err)
	}

	// Retrieve messages for this conversation.
	rows, err := db.db.Query("SELECT id, conversation_id, sender_id, content, reply_to, sent_at FROM messages WHERE conversation_id = ? ORDER BY sent_at ASC", conversationID)
	if err != nil {
		return &conv, nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.ReplyTo, &msg.SentAt); err != nil {
			return &conv, messages, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return &conv, messages, fmt.Errorf("error iterating messages: %w", err)
	}

	return &conv, messages, nil
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
		group_photo TEXT, -- New column for the group photo URL
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

/// SendMessage automatically creates a conversation if it doesn't exist and inserts a message.
func (db *appdbimpl) SendMessage(senderID, receiverID, content string, isGroup bool, groupID string) (string, error) {
	var conversationID string

	// **1. Check if conversation exists**
	if isGroup {
		// If it's a group, use the group ID as the conversation ID
		conversationID = groupID
	} else {
		// For private messages, check if a conversation already exists
		query := `
			SELECT id FROM conversations 
			WHERE is_group = 0 
			AND id IN (
				SELECT group_id FROM group_members 
				WHERE user_id IN (?, ?)
				GROUP BY group_id HAVING COUNT(user_id) = 2
			)
			LIMIT 1;
		`
		err := db.db.QueryRow(query, senderID, receiverID).Scan(&conversationID)
		if err == sql.ErrNoRows {
			// **2. If conversation does not exist, create a new one**
			conversationID, _ = GenerateNewID() // Generate a new UUID
			_, err = db.db.Exec("INSERT INTO conversations (id, name, is_group) VALUES (?, NULL, 0)", conversationID)
			if err != nil {
				return "", fmt.Errorf("failed to create conversation: %w", err)
			}

			// Add both users to the `group_members` table
			_, err = db.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?), (?, ?)", conversationID, senderID, conversationID, receiverID)
			if err != nil {
				return "", fmt.Errorf("failed to add users to group_members: %w", err)
			}
		} else if err != nil {
			return "", fmt.Errorf("error checking conversation: %w", err)
		}
	}

	// **3. Insert message into messages table**
	messageID, _ := GenerateNewID()
	_, err := db.db.Exec("INSERT INTO messages (id, conversation_id, sender_id, content) VALUES (?, ?, ?, ?)",
		messageID, conversationID, senderID, content)
	if err != nil {
		return "", fmt.Errorf("failed to insert message: %w", err)
	}

	return messageID, nil
}


// ForwardMessage forwards a message to another conversation.
func (db *appdbimpl) ForwardMessage(originalMessageID, targetConversationID, senderID string) (string, error) {
	// Retrieve the original message.
	var originalMessage Message
	err := db.db.QueryRow("SELECT content, reply_to FROM messages WHERE id = ?", originalMessageID).
		Scan(&originalMessage.Content, &originalMessage.ReplyTo)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("original message not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve original message: %w", err)
	}

	// Generate a new ID for the forwarded message.
	newMessageID, err := GenerateNewID()
	if err != nil {
		return "", fmt.Errorf("failed to generate new message ID: %w", err)
	}

	// Insert the forwarded message into the messages table.
	_, err = db.db.Exec("INSERT INTO messages (id, conversation_id, sender_id, content, reply_to) VALUES (?, ?, ?, ?, ?)",
		newMessageID, targetConversationID, senderID, originalMessage.Content, originalMessage.ReplyTo)
	if err != nil {
		return "", fmt.Errorf("failed to forward message: %w", err)
	}

	return newMessageID, nil
}

// CommentMessage inserts a reaction (comment) for a message into the message_reactions table.
func (db *appdbimpl) CommentMessage(messageID, userID, reaction string) error {
	// Insert the reaction. Using INSERT OR REPLACE to allow updating an existing reaction.
	_, err := db.db.Exec(
		"INSERT OR REPLACE INTO message_reactions (message_id, user_id, reaction) VALUES (?, ?, ?)",
		messageID, userID, reaction,
	)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}
	return nil
}

// UncommentMessage removes a reaction (comment) for a message from the message_reactions table.
func (db *appdbimpl) UncommentMessage(messageID, userID string) error {
	_, err := db.db.Exec(
		"DELETE FROM message_reactions WHERE message_id = ? AND user_id = ?",
		messageID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove comment: %w", err)
	}
	return nil
}

// DeleteMessage removes a message from the messages table. (In a real system, you might also mark it as deleted.)
func (db *appdbimpl) DeleteMessage(messageID, senderID string) error {
	// Optionally verify that the sender is the one deleting the message.
	res, err := db.db.Exec("DELETE FROM messages WHERE id = ? AND sender_id = ?", messageID, senderID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("no message deleted (perhaps invalid message ID or sender mismatch)")
	}
	return nil
}

// AddToGroup adds a user to a group by inserting a record into the group_members table.
func (db *appdbimpl) AddToGroup(groupID, userID string) error {
	_, err := db.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}
	return nil
}

// LeaveGroup removes a user from a group.
func (db *appdbimpl) LeaveGroup(groupID, userID string) error {
	res, err := db.db.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("user not a member of the group")
	}
	return nil
}

// SetGroupName updates the name of a group (in the conversations table).
func (db *appdbimpl) SetGroupName(groupID, newName string) error {
	// Only update if the conversation is a group (is_group = 1)
	_, err := db.db.Exec("UPDATE conversations SET name = ? WHERE id = ? AND is_group = 1", newName, groupID)
	if err != nil {
		return fmt.Errorf("failed to update group name: %w", err)
	}
	return nil
}

// SetGroupPhoto updates the group photo for a group conversation.
// Ensure that your conversations table has a column called group_photo.
func (db *appdbimpl) SetGroupPhoto(groupID, photoUrl string) error {
	_, err := db.db.Exec("UPDATE conversations SET group_photo = ? WHERE id = ? AND is_group = 1", photoUrl, groupID)
	if err != nil {
		return fmt.Errorf("failed to update group photo: %w", err)
	}
	return nil
}
