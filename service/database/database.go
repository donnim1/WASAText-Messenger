package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

// AppDatabase is the high-level interface for the DB.
type AppDatabase interface {
	CreateUser(username string) (string, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUserName(userID, newName string) error
	UpdateUserPhoto(userID, photoUrl string) error

	// ListUsers retrieves all users.
	ListUsers() ([]User, error)

	// In the AppDatabase interface:
	GetChatPartner(conversationID, currentUserID string) (*User, error)

	GetConversationsByUserID(userID string) ([]Conversation, error)
	GetConversation(conversationID string) (*Conversation, []Message, error)

	SendMessage(senderID, receiverID, content string, isGroup bool, groupID, conversationID string) (string, string, error)
	ForwardMessage(originalMessageID, targetConversationID, senderID string) (string, error)
	CommentMessage(messageID, userID, reaction string) error
	UncommentMessage(messageID, userID string) error
	DeleteMessage(messageID, senderID string) error

	// CreateGroup creates a new group conversation and adds the creator as a member.
	CreateGroup(creatorID, groupName, groupPhoto string) (string, error)
	// Register the GET /groups endpoint.
	GetGroupsByUserID(userID string) ([]Conversation, error)

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
	PhotoURL sql.NullString // Now handles NULL values; optional profile photo URL.
}

// Conversation represents a conversation record.
type Conversation struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsGroup   bool   `json:"is_group"`
	CreatedAt string `json:"created_at"` // Timestamp of when the conversation was created
}

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	Content        string
	ReplyTo        sql.NullString // If replies are optional
	SentAt         string         // using string for simplicity; you may use time.Time
}

// GetGroupsByUserID retrieves all group conversations associated with a user.
func (db *appdbimpl) GetGroupsByUserID(userID string) ([]Conversation, error) {
	query := `
      SELECT c.id, c.name, c.is_group, c.created_at, c.group_photo
      FROM conversations c
      JOIN group_members gm ON c.id = gm.group_id
      WHERE gm.user_id = ? AND c.is_group = 1
    `
	rows, err := db.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups: %w", err)
	}
	defer rows.Close()

	var groups []Conversation
	for rows.Next() {
		var conv Conversation
		// If you want to include group_photo, consider adding it to the Conversation struct.
		if err := rows.Scan(&conv.ID, &conv.Name, &conv.IsGroup, &conv.CreatedAt, new(interface{})); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, conv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return groups, nil
}

func (db *appdbimpl) CreateGroup(creatorID, groupName, groupPhoto string) (string, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return "", fmt.Errorf("transaction start failed: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			log.Printf("tx.Rollback() error: %v", rbErr)
		}
	}()

	// Generate group ID
	groupID, err := GenerateNewID()
	if err != nil {
		return "", fmt.Errorf("ID generation failed: %w", err)
	}

	// Create conversation
	_, err = tx.Exec(`INSERT INTO conversations 
        (id, name, is_group, group_photo) 
        VALUES (?, ?, 1, ?)`,
		groupID, groupName, groupPhoto)
	if err != nil {
		return "", fmt.Errorf("conversation creation failed: %w", err)
	}

	// Add creator as member
	_, err = tx.Exec(`INSERT INTO group_members 
        (group_id, user_id) 
        VALUES (?, ?)`,
		groupID, creatorID)
	if err != nil {
		return "", fmt.Errorf("member addition failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("transaction commit failed: %w", err)
	}

	return groupID, nil
}

func (db *appdbimpl) ListUsers() ([]User, error) {
	rows, err := db.db.Query("SELECT id, username, photo_url FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.PhotoURL); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return users, nil
}

// GetConversation retrieves a conversation and all its messages.
func (db *appdbimpl) GetConversation(conversationID string) (*Conversation, []Message, error) {
	// Retrieve conversation details.
	var conv Conversation
	// Use a sql.NullString for the name column.
	var name sql.NullString
	err := db.db.QueryRow("SELECT id, name, is_group, created_at FROM conversations WHERE id = ?", conversationID).
		Scan(&conv.ID, &name, &conv.IsGroup, &conv.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil // Conversation not found.
	} else if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve conversation: %w", err)
	}
	// Convert sql.NullString to string.
	if name.Valid {
		conv.Name = name.String
	} else {
		conv.Name = ""
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

	// Create conversations table if not exists.
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

	// Create messages table if not exists.
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

	// Create group_members table if not exists.
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

	// Create message reactions table if not exists.
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

	// In database.go - New() function after creating tables:
	_, err = db.Exec(`CREATE TRIGGER IF NOT EXISTS check_private_members
	BEFORE INSERT ON group_members
	FOR EACH ROW
	WHEN (SELECT is_group FROM conversations WHERE id = NEW.group_id) = 0
	BEGIN
		SELECT RAISE(ABORT, 'private chat cannot have more than two members')
		WHERE (SELECT COUNT(*) FROM group_members WHERE group_id = NEW.group_id) >= 2;
	END;`)
	if err != nil {
		return nil, fmt.Errorf("error creating check_private_members trigger: %w", err)
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
	if errors.Is(err, sql.ErrNoRows) {
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
	if err == nil {
		return errors.New("username already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
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
	// Execute update query
	res, err := db.db.Exec("UPDATE users SET photo_url = ? WHERE id = ?", photoUrl, userID)
	if err != nil {
		log.Printf("❌ Database update error: %v", err)
		return fmt.Errorf("failed to update photo: %w", err)
	}

	// Check if any row was updated
	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("❌ Error checking rows affected: %v", err)
		return fmt.Errorf("failed to update photo: %w", err)
	}
	if affected == 0 {
		log.Println("⚠️ No user found to update photo")
		return fmt.Errorf("no user found with given ID")
	}

	log.Println("✅ Photo successfully updated in database")
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
		var conv Conversation
		var name sql.NullString
		if err := rows.Scan(&conv.ID, &name, &conv.IsGroup, &conv.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		if name.Valid {
			conv.Name = name.String
		} else {
			conv.Name = ""
		}
		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return conversations, nil
}

// SendMessage inserts a new message and returns the generated messageID and conversationID.
// If conversationID is empty, creates a new conversation for the users.
func (db *appdbimpl) SendMessage(userID, receiverID, content string, isGroup bool, groupID, conversationID string) (string, string, error) {
	// If this is a private message and conversationID is empty, create a new conversation.
	if !isGroup && conversationID == "" {
		newConvID, err := db.createConversation(userID, receiverID)
		if err != nil {
			return "", "", err
		}
		conversationID = newConvID
	}

	// For group messages, you might already have groupID. (Handle accordingly)

	// Simulate creating a new message record.
	newMessageID := generateNewMessageID() // Assume this function generates a unique message ID.
	currentTime := time.Now().UTC().Format(time.RFC3339)

	// Insert the new message into your messages table.
	// (Replace the following with your actual DB insert query)
	fmt.Printf("Inserting message: ID=%s, ConversationID=%s, UserID=%s, Content=%s, SentAt=%s\n",
		newMessageID, conversationID, userID, content, currentTime)

	// If the insert fails, return an error.
	// e.g., err := db.insertMessage(query, params...)
	// if err != nil {
	//     return "", "", err
	// }

	// Return new messageID and the conversationID.
	return newMessageID, conversationID, nil
}

// createConversation creates a new conversation between two users and returns the new conversation ID.
func (db *appdbimpl) createConversation(userID, receiverID string) (string, error) {
	newConversationID := generateNewConversationID() // Assume this function generates a unique conversation ID.
	// Insert a new conversation record into your conversations table.
	// (Replace the following with your actual DB insert query)
	fmt.Printf("Creating new conversation: ID=%s between user %s and receiver %s\n",
		newConversationID, userID, receiverID)
	// If the insert fails, return an error.
	// e.g., err := db.insertConversation(query, params...)
	// if err != nil {
	//     return "", err
	// }
	return newConversationID, nil
}

// Dummy functions to simulate ID generation.
func generateNewMessageID() string {
	// TODO: Implement a proper unique ID generator.
	return "msg123456789"
}

func generateNewConversationID() string {
	// TODO: Implement a proper unique ID generator.
	return "conv123456789"
}

// ForwardMessage forwards a message to another conversation.
func (db *appdbimpl) ForwardMessage(originalMessageID, targetConversationID, senderID string) (string, error) {
	var originalMessage Message
	err := db.db.QueryRow("SELECT content, reply_to FROM messages WHERE id = ?", originalMessageID).
		Scan(&originalMessage.Content, &originalMessage.ReplyTo)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("original message not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve original message: %w", err)
	}
	newMessageID, err := GenerateNewID()
	if err != nil {
		return "", fmt.Errorf("failed to generate new message ID: %w", err)
	}
	_, err = db.db.Exec("INSERT INTO messages (id, conversation_id, sender_id, content, reply_to) VALUES (?, ?, ?, ?, ?)",
		newMessageID, targetConversationID, senderID, originalMessage.Content, originalMessage.ReplyTo)
	if err != nil {
		return "", fmt.Errorf("failed to forward message: %w", err)
	}
	return newMessageID, nil
}

// CommentMessage inserts a reaction (comment) for a message into the message_reactions table.
func (db *appdbimpl) CommentMessage(messageID, userID, reaction string) error {
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

// DeleteMessage removes a message from the messages table.
func (db *appdbimpl) DeleteMessage(messageID, senderID string) error {
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

// In database.go - AddToGroup function:
func (db *appdbimpl) AddToGroup(groupID, userID string) error {
	// Verify target is a group
	var isGroup bool
	err := db.db.QueryRow("SELECT is_group FROM conversations WHERE id = ?", groupID).Scan(&isGroup)
	if err != nil {
		return fmt.Errorf("group verification failed: %w", err)
	}
	if !isGroup {
		return errors.New("cannot add members to private chats")
	}

	// Check existing membership
	var exists bool
	err = db.db.QueryRow("SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?",
		groupID, userID).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("membership check failed: %w", err)
	}
	if exists {
		return nil // Already member, return success
	}

	// Proceed with insertion
	_, err = db.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)",
		groupID, userID)
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
	_, err := db.db.Exec("UPDATE conversations SET name = ? WHERE id = ? AND is_group = 1", newName, groupID)
	if err != nil {
		return fmt.Errorf("failed to update group name: %w", err)
	}
	return nil
}

// SetGroupPhoto updates the group photo for a group conversation.
func (db *appdbimpl) SetGroupPhoto(groupID, photoUrl string) error {
	_, err := db.db.Exec("UPDATE conversations SET group_photo = ? WHERE id = ? AND is_group = 1", photoUrl, groupID)
	if err != nil {
		return fmt.Errorf("failed to update group photo: %w", err)
	}
	return nil
}

// GetChatPartner returns the user (other than currentUserID) in the private conversation.
func (db *appdbimpl) GetChatPartner(conversationID, currentUserID string) (*User, error) {
	// This query joins the group_members and users tables to find the other user.
	row := db.db.QueryRow(
		`SELECT u.id, u.username, u.photo_url
		 FROM group_members gm
		 JOIN users u ON gm.user_id = u.id
		 WHERE gm.group_id = ? AND u.id != ?
		 LIMIT 1`,
		conversationID, currentUserID,
	)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.PhotoURL)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
