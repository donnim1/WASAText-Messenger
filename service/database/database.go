package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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
	GetConversationBetween(userID1, userID2 string) (*Conversation, error)
	GetConversationsByUserID(userID string) ([]Conversation, error)
	GetConversation(conversationID string) (*Conversation, []Message, error)

	SendMessage(senderID, receiverID, content string, isGroup bool, groupID, conversationID string, replyTo string) (string, string, error)
	ForwardMessage(originalMessageID, targetConversationID, senderID string) (string, error)
	CommentMessage(messageID, userID, reaction string) error
	UncommentMessage(messageID, userID string) error
	DeleteMessage(messageID, senderID string) error

	UpdateMessageStatus(messageID, status, userID string) error

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
	PhotoUrl sql.NullString // Now handles NULL values; optional profile photo URL.
}

// Conversation represents a conversation record.
type Conversation struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	IsGroup            bool           `json:"is_group"`
	CreatedAt          string         `json:"created_at"`
	Messages           []Message      `json:"messages"`
	PhotoUrl           string         `json:"photoUrl"` // Timestamp of when the conversation was created
	Members            []User         `json:"members"`
	LastMessageContent sql.NullString `json:"last_message_content"` // New field for the last message content
	LastMessageSentAt  sql.NullString `json:"last_message_sent_at"` // New field for the last message sent time
}

type Message struct {
	ID             string       `json:"ID"`
	ConversationID string       `json:"ConversationID"`
	SenderID       string       `json:"SenderID"`
	Content        string       `json:"Content"`
	ReplyTo        string       `json:"ReplyTo,omitempty"` // Changed to string
	SentAt         string       `json:"SentAt"`
	Reactions      []Reaction   `json:"reactions"` // New field for reactions
	Status         string       `json:"status"`    // "pending", "sent", "delivered", "read"
	DeliveredAt    sql.NullTime `json:"deliveredAt,omitempty"`
	ReadAt         sql.NullTime `json:"readAt,omitempty"`
}

type Reaction struct {
	Reaction string `json:"reaction"`
	UserName string `json:"userName"`
	UserID   string `json:"userID"`
}

// GetConversationBetween looks for an existing conversation that includes both userID1 and userID2.
func (db *appdbimpl) GetConversationBetween(userID1, userID2 string) (*Conversation, error) {
	// Example SQL: Adjust based on your conversation/members table structure.
	query := `
        SELECT c.id
        FROM conversations c
        INNER JOIN group_members cm1 ON c.id = cm1.group_id
        INNER JOIN group_members cm2 ON c.id = cm2.group_id
        WHERE cm1.user_id = ? AND cm2.user_id = ?
		AND c.is_group = 0
        LIMIT 1
    `
	var convID string
	err := db.db.QueryRow(query, userID1, userID2).Scan(&convID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("conversation not found")
		}
		return nil, err
	}

	// Retrieve messages for the conversation.
	msgQuery := `
        SELECT id, sender_id, content, sent_at 
        FROM messages 
        WHERE conversation_id = ?
        ORDER BY sent_at ASC
    `
	rows, err := db.db.Query(msgQuery, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.Content, &msg.SentAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	// Check for any errors encountered during iteration.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &Conversation{
		ID:       convID,
		Messages: messages,
	}, nil
}

// GetGroupsByUserID retrieves all group conversations associated with a user.
func (db *appdbimpl) GetGroupsByUserID(userID string) ([]Conversation, error) {
	// Query to get groups that the user is a member of (only groups).
	query := `
     SELECT c.id, c.name, c.is_group, c.created_at, c.group_photo, 
            NULL as last_message_content, NULL as last_message_sent_at
     FROM conversations c
     INNER JOIN group_members gm ON c.id = gm.group_id
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
		var name, groupPhoto sql.NullString
		var lastMsgContent, lastMsgSentAt sql.NullString
		if err := rows.Scan(&conv.ID, &name, &conv.IsGroup, &conv.CreatedAt, &groupPhoto, &lastMsgContent, &lastMsgSentAt); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}

		if name.Valid {
			conv.Name = name.String
		}
		if groupPhoto.Valid {
			conv.PhotoUrl = groupPhoto.String
		}
		if lastMsgContent.Valid {
			conv.LastMessageContent = lastMsgContent
		}
		if lastMsgSentAt.Valid {
			conv.LastMessageSentAt = lastMsgSentAt
		}

		// Retrieve members for this group.
		memberQuery := `
            SELECT u.id, u.username, u.photo_url
            FROM users u
            INNER JOIN group_members gm ON u.id = gm.user_id
            WHERE gm.group_id = ?
        `
		memberRows, err := db.db.Query(memberQuery, conv.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch group members: %w", err)
		}

		var members []User
		for memberRows.Next() {
			var member User
			if err := memberRows.Scan(&member.ID, &member.Username, &member.PhotoUrl); err != nil {
				memberRows.Close()
				return nil, fmt.Errorf("failed to scan group member: %w", err)
			}
			members = append(members, member)
		}
		// Explicitly close after processing.
		memberRows.Close()
		conv.Members = members

		groups = append(groups, conv)
	}
	if err = rows.Err(); err != nil {
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
		if err := rows.Scan(&user.ID, &user.Username, &user.PhotoUrl); err != nil {
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
	// Using sql.NullString for optional fields.
	var name, groupPhoto sql.NullString
	err := db.db.QueryRow(`
        SELECT id, name, is_group, created_at, group_photo 
        FROM conversations 
        WHERE id = ?`, conversationID).
		Scan(&conv.ID, &name, &conv.IsGroup, &conv.CreatedAt, &groupPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil // Conversation not found.
	} else if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve conversation: %w", err)
	}

	if name.Valid {
		conv.Name = name.String
	} else {
		conv.Name = ""
	}

	if groupPhoto.Valid {
		conv.PhotoUrl = groupPhoto.String
	} else {
		conv.PhotoUrl = ""
	}

	// Retrieve messages for this conversation.
	var messages []Message
	rows, err := db.db.Query(`
	SELECT id, conversation_id, sender_id, content, reply_to, sent_at, status, deliveredAt, readAt 
	FROM messages 
	WHERE conversation_id = ? 
	ORDER BY sent_at ASC`, conversationID)
	if err != nil {
		return &conv, messages, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		var replyTo sql.NullString
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &replyTo, &msg.SentAt, &msg.Status, &msg.DeliveredAt, &msg.ReadAt); err != nil {
			return &conv, messages, fmt.Errorf("failed to scan message: %w", err)
		}
		if replyTo.Valid {
			msg.ReplyTo = replyTo.String
		} else {
			msg.ReplyTo = ""
		}
		messages = append(messages, msg)
	}
	// Check for iteration errors.
	if err := rows.Err(); err != nil {
		return &conv, messages, fmt.Errorf("rows iteration error: %w", err)
	}

	// Retrieve reactions for all messages in one query if there are any messages.
	if len(messages) > 0 {
		// Build placeholder string for SQL IN clause.
		placeholders := "?"
		args := []interface{}{messages[0].ID}
		for i := 1; i < len(messages); i++ {
			placeholders += ",?"
			args = append(args, messages[i].ID)
		}

		// Query reactions for all message IDs and join to retrieve the username.
		reactionRows, err := db.db.Query(
			"SELECT mr.message_id, mr.reaction, u.id, u.username FROM message_reactions mr JOIN users u ON mr.user_id = u.id WHERE mr.message_id IN ("+placeholders+")",
			args...,
		)
		if err != nil {
			return &conv, messages, fmt.Errorf("failed to query reactions: %w", err)
		}
		defer reactionRows.Close()

		// Build a mapping from message ID to a slice of reaction objects.
		reactionsMap := make(map[string][]Reaction)
		for reactionRows.Next() {
			var messageID, reaction, userID, userName string
			if err := reactionRows.Scan(&messageID, &reaction, &userID, &userName); err != nil {
				return &conv, messages, fmt.Errorf("failed to scan reaction: %w", err)
			}
			reactionsMap[messageID] = append(reactionsMap[messageID], Reaction{
				Reaction: reaction,
				UserID:   userID,
				UserName: userName,
			})
		}

		// Attach reactions to each message.
		for i, msg := range messages {
			if r, ok := reactionsMap[msg.ID]; ok {
				// Make sure Message.Reactions is defined to accept these reaction objects.
				messages[i].Reactions = r
			} else {
				messages[i].Reactions = []Reaction{}
			}
		}
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
		status TEXT NOT NULL DEFAULT 'sent',
        deliveredAt DATETIME,
        readAt DATETIME,
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

	// In your New() function, after creating other tables, add:
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS message_read_receipts (
		message_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		read_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (message_id, user_id),
		FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return nil, fmt.Errorf("error creating message_read_receipts table: %w", err)
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
	err := row.Scan(&user.ID, &user.Username, &user.PhotoUrl)
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
	query := `
    SELECT 
      c.id, 
      c.name, 
      c.is_group, 
      c.created_at, 
      c.group_photo,
      COALESCE(
        (SELECT content FROM messages WHERE conversation_id = c.id ORDER BY sent_at DESC LIMIT 1), 
        ''
      ) AS last_message_content,
      COALESCE(
        (SELECT sent_at FROM messages WHERE conversation_id = c.id ORDER BY sent_at DESC LIMIT 1), 
        ''
      ) AS last_message_sent_at
    FROM conversations c
    JOIN group_members gm ON c.id = gm.group_id
    WHERE gm.user_id = ?
  `
	rows, err := db.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch conversations: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		var name sql.NullString
		var groupPhoto sql.NullString
		if err := rows.Scan(
			&conv.ID,
			&name,
			&conv.IsGroup,
			&conv.CreatedAt,
			&groupPhoto,
			&conv.LastMessageContent, // Now scanned as string (with COALESCE, never NULL)
			&conv.LastMessageSentAt,  // Now scanned as string too.
		); err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		if name.Valid {
			conv.Name = name.String
		} else {
			conv.Name = ""
		}
		if groupPhoto.Valid {
			conv.PhotoUrl = groupPhoto.String
		} else {
			conv.PhotoUrl = ""
		}
		conversations = append(conversations, conv)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return conversations, nil
}

// SendMessage inserts a new message and returns the generated messageID and conversationID.
// If conversationID is empty, creates a new conversation for the users.
func (db *appdbimpl) SendMessage(userID, receiverID, content string, isGroup bool, groupID, conversationID, replyTo string) (string, string, error) {
	// For private messages, check if a conversation already exists.
	if !isGroup {
		if conversationID == "" {
			existingConv, err := db.GetPrivateConversation(userID, receiverID)
			if err != nil {
				return "", "", fmt.Errorf("error checking for existing conversation: %w", err)
			}
			if existingConv != nil {
				conversationID = existingConv.ID
			} else {
				conversationID, err = db.createConversation(userID, receiverID)
				if err != nil {
					return "", "", fmt.Errorf("failed to create conversation: %w", err)
				}
			}
		}
	}

	newMessageID, err := GenerateNewID()
	if err != nil {
		return "", "", fmt.Errorf("GenerateNewID error: %w", err)
	}

	currentTime := time.Now().UTC().Format(time.RFC3339)

	// Updated query to include reply_to column.
	query := `INSERT INTO messages (id, conversation_id, sender_id, content, reply_to, sent_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.db.Exec(query, newMessageID, conversationID, userID, content, replyTo, currentTime)
	if err != nil {
		return "", "", fmt.Errorf("failed to insert message, query error: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", "", fmt.Errorf("failed to check affected rows: %w", err)
	}
	if affected == 0 {
		return "", "", fmt.Errorf("no rows affected")
	}

	return newMessageID, conversationID, nil
}

// createConversation creates a new conversation between two users and returns the new conversation ID.
func (db *appdbimpl) createConversation(userID, receiverID string) (string, error) {
	// Generate a unique conversation ID.
	newConversationID, err := GenerateNewID()
	if err != nil {
		return "", fmt.Errorf("failed to generate conversation ID: %w", err)
	}

	// Insert a new conversation record into the conversations table.
	// For a private conversation, you might leave the name blank.
	query := `INSERT INTO conversations (id, name, is_group, created_at) VALUES (?, ?, 0, ?)`
	currentTime := time.Now().UTC().Format(time.RFC3339)
	_, err = db.db.Exec(query, newConversationID, "", currentTime)
	if err != nil {
		return "", fmt.Errorf("failed to insert conversation: %w", err)
	}

	// Add both users as members of this private conversation.
	_, err = db.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", newConversationID, userID)
	if err != nil {
		return "", fmt.Errorf("failed to add creator to conversation: %w", err)
	}
	_, err = db.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", newConversationID, receiverID)
	if err != nil {
		return "", fmt.Errorf("failed to add receiver to conversation: %w", err)
	}

	return newConversationID, nil
}

// (Duplicate GenerateNewID function removed)

// ForwardMessage forwards a message to another conversation.
func (db *appdbimpl) ForwardMessage(originalMessageID, targetConversationID, senderID string) (string, error) {
	// Retrieve the original message content.
	var originalContent string
	err := db.db.QueryRow("SELECT content FROM messages WHERE id = ?", originalMessageID).
		Scan(&originalContent)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve original message: %w", err)
	}

	newMessageID, err := GenerateNewID()
	if err != nil {
		return "", fmt.Errorf("failed to generate new message ID: %w", err)
	}

	trimmed := strings.TrimSpace(originalContent)
	var newContent string
	if strings.HasPrefix(trimmed, "data:image") {
		// Wrap the image with a caption as HTML.
		newContent = fmt.Sprintf(`<div class="forward-caption">Forwarded from you:</div><img src="%s" alt="Image message" class="sent-image" />`, trimmed)
	} else {
		newContent = fmt.Sprintf("Forwarded from you: %s", originalContent)
	}

	currentTime := time.Now().UTC().Format(time.RFC3339)
	_, err = db.db.Exec(
		"INSERT INTO messages (id, conversation_id, sender_id, content, reply_to, sent_at) VALUES (?, ?, ?, ?, ?, ?)",
		newMessageID, targetConversationID, senderID, newContent, nil, currentTime)
	if err != nil {
		return "", fmt.Errorf("failed to insert forwarded message: %w", err)
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
	log.Printf("Deleting reaction for message %s, user %s", messageID, userID)

	query := "DELETE FROM message_reactions WHERE message_id = ? AND user_id = ?"
	result, err := db.db.Exec(query, messageID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deletion: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("reaction not found")
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
	// Verify the conversation is a group.
	var isGroup bool
	err := db.db.QueryRow("SELECT is_group FROM conversations WHERE id = ?", groupID).Scan(&isGroup)
	if err != nil {
		return fmt.Errorf("group verification failed: %w", err)
	}
	if !isGroup {
		return errors.New("cannot add members to private chats")
	}

	// Check if the user is already a member.
	var count int
	err = db.db.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("membership check failed: %w", err)
	}
	if count > 0 {
		// Already a member; consider this a success.
		return nil
	}

	// Insert the new group membership.
	_, err = db.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
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
	err := row.Scan(&user.ID, &user.Username, &user.PhotoUrl)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetPrivateConversation looks for an existing private conversation between two users.
func (db *appdbimpl) GetPrivateConversation(userID, receiverID string) (*Conversation, error) {
	// Adjust this SQL according to your schema and how you store private conversations.
	query := `
    SELECT c.id, c.name, c.is_group, c.created_at 
    FROM conversations c
    JOIN group_members gm ON c.id = gm.group_id
    WHERE c.is_group = 0 
      AND gm.user_id = ?
      AND c.id IN (
          SELECT group_id FROM group_members WHERE user_id = ?
      )
    LIMIT 1
    `
	row := db.db.QueryRow(query, userID, receiverID)
	var conv Conversation
	if err := row.Scan(&conv.ID, &conv.Name, &conv.IsGroup, &conv.CreatedAt); err != nil {
		// If the error is due to no rows, return nil without an error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &conv, nil
}
func (db *appdbimpl) UpdateMessageStatus(messageID, status, userID string) error {
	currentTime := time.Now().UTC().Format(time.RFC3339)

	if status == "delivered" {
		// For delivered status, update as usual.
		query := "UPDATE messages SET status = ?, deliveredAt = ? WHERE id = ? AND status != 'read'"
		_, err := db.db.Exec(query, status, currentTime, messageID)
		return err
	} else if status == "read" {
		// Insert (or replace) a read receipt for the current user.
		_, err := db.db.Exec(
			"INSERT OR REPLACE INTO message_read_receipts (message_id, user_id, read_at) VALUES (?, ?, ?)",
			messageID, userID, currentTime,
		)
		if err != nil {
			return fmt.Errorf("failed to insert read receipt: %w", err)
		}

		// Retrieve the conversation ID for this message.
		var conversationID string
		err = db.db.QueryRow("SELECT conversation_id FROM messages WHERE id = ?", messageID).Scan(&conversationID)
		if err != nil {
			return fmt.Errorf("failed to retrieve conversation id: %w", err)
		}

		// Check if this conversation is a group chat.
		var isGroup bool
		err = db.db.QueryRow("SELECT is_group FROM conversations WHERE id = ?", conversationID).Scan(&isGroup)
		if err != nil {
			return fmt.Errorf("failed to check conversation type: %w", err)
		}

		if !isGroup {
			// For private chats, update immediately.
			query := "UPDATE messages SET status = ?, readAt = ? WHERE id = ?"
			_, err := db.db.Exec(query, status, currentTime, messageID)
			return err
		} else {
			// For group chats, count total group members.
			var totalMembers int
			err = db.db.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = ?", conversationID).Scan(&totalMembers)
			if err != nil {
				return fmt.Errorf("failed to count group members: %w", err)
			}

			// Count the number of read receipts for this message.
			var totalRead int
			err = db.db.QueryRow("SELECT COUNT(*) FROM message_read_receipts WHERE message_id = ?", messageID).Scan(&totalRead)
			if err != nil {
				return fmt.Errorf("failed to count read receipts: %w", err)
			}

			// Logging the counts for debugging purposes.
			log.Printf("UpdateMessageStatus: messageID=%s, totalMembers=%d, totalRead=%d", messageID, totalMembers, totalRead)

			// Update the message's status to "read" when totalRead is greater than or equal to totalMembers - 1 (since the sender does not send a receipt).
			if totalRead >= (totalMembers - 1) {
				query := "UPDATE messages SET status = ?, readAt = ? WHERE id = ?"
				_, err := db.db.Exec(query, status, currentTime, messageID)
				return err
			}

			// Not all members (except sender) have read the message yet.
			return nil
		}
	}
	return fmt.Errorf("unsupported status update: %s", status)
}
