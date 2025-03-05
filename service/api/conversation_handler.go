package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/donnim1/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// ConversationResponse is the payload returned by GetConversationByReceiver.
type ConversationResponse struct {
	ConversationID string             `json:"conversationId"`
	Messages       []database.Message `json:"messages"`
}

// GetConversationByReceiver handles requests to fetch a conversation by receiverID.
// It expects the current user's ID to be provided via a header (e.g. "X-User-Id").
// Route: GET /conversationsfor/:receiverId
func (rt *_router) GetConversationByReceiver(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	receiverId := ps.ByName("receiverId")
	// Extract current user ID from the Authorization header.
	currentUserId, err := rt.getAuthenticatedUserID(r)
	if err != nil || currentUserId == "" {
		http.Error(w, "Missing current user ID", http.StatusBadRequest)
		return
	}

	// Use a database helper to fetch a conversation between the two users.
	conv, err := rt.db.GetConversationBetween(currentUserId, receiverId)
	if err != nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	res := ConversationResponse{
		ConversationID: conv.ID,
		Messages:       conv.Messages,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("Error encoding conversation response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// getMyConversationsResponse defines the structure of the response after fetching conversations.
type getMyConversationsResponse struct {
	Conversations []Conversation `json:"conversations"`
}

// Conversation defines the API's conversation structure.
type Conversation struct {
	ID                 string `json:"id"`       // Conversation ID
	Name               string `json:"name"`     // Name of the conversation (group name or null for private chats)
	IsGroup            bool   `json:"is_group"` // True if it's a group chat, otherwise false
	CreatedAt          string `json:"created_at"`
	PhotoUrl           string `json:"group_photo"`
	LastMessageContent string `json:"last_message_content"` // new field
	LastMessageSentAt  string `json:"last_message_sent_at"`
	// new field
	// Timestamp of when the conversation was created
}

// getMyConversations retrieves all conversations for the authenticated user.
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Get the authenticated user ID from the Authorization header.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Fetch conversations from the database using the authenticated user ID.
	convs, err := rt.db.GetConversationsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve conversations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Convert each database conversation into the API's Conversation type.
	var apiConversations []Conversation
	for _, conv := range convs {
		// For private chats (is_group false), if the conversation name is empty,
		// try to get the chat partner's username and photo.
		if conv.Name == "" && !conv.IsGroup {
			partner, err := rt.db.GetChatPartner(conv.ID, userID)
			if err == nil && partner != nil {
				conv.Name = partner.Username
				if partner.PhotoUrl.Valid {
					conv.PhotoUrl = partner.PhotoUrl.String
				} else {
					conv.PhotoUrl = ""
				}
			}
		}
		// Parse and format the created_at timestamp.
		var formattedCreatedAt string
		t, err := time.Parse(time.RFC3339, conv.CreatedAt)
		if err != nil {
			formattedCreatedAt = conv.CreatedAt
		} else {
			formattedCreatedAt = t.Format("2006-01-02 15:04:05")
		}
		apiConversations = append(apiConversations, Conversation{
			ID:                 conv.ID,
			Name:               conv.Name,
			IsGroup:            conv.IsGroup,
			CreatedAt:          formattedCreatedAt,
			PhotoUrl:           conv.PhotoUrl,
			LastMessageContent: conv.LastMessageContent.String, // New field.
			LastMessageSentAt:  conv.LastMessageSentAt.String,  // New field. // This field will contain either group_photo or private chat photo.
		})
	}

	// 4. Return the list of conversations.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(getMyConversationsResponse{
		Conversations: apiConversations,
	}); err != nil {
		log.Printf("Error encoding getMyConversations response: %v", err)
	}
}

// getConversation handles GET requests to /conversations/:conversationId.
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 1. Validate the Authorization header.
	currentUserId, err := rt.getAuthenticatedUserID(r)
	if err != nil || currentUserId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Extract the conversation ID.
	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	// 3. Retrieve conversation details and messages from the database.
	conv, messages, err := rt.db.GetConversation(conversationID)
	if err != nil {
		http.Error(w, "Failed to retrieve conversation: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if conv == nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// 4. For private chats with empty name, attempt to fetch the chat partner's details.
	if !conv.IsGroup && conv.Name == "" {
		partner, err := rt.db.GetChatPartner(conv.ID, currentUserId)
		if err == nil && partner != nil {
			conv.Name = partner.Username
			if partner.PhotoUrl.Valid {
				conv.PhotoUrl = partner.PhotoUrl.String
			}
		}
	}

	// 5. Parse and format the created_at timestamp.
	var formattedCreatedAt string
	t, err := time.Parse(time.RFC3339, conv.CreatedAt)
	if err != nil {
		formattedCreatedAt = conv.CreatedAt
	} else {
		formattedCreatedAt = t.Format("2006-01-02 15:04:05")
	}

	// 6. Build the API conversation structure.
	apiConv := Conversation{
		ID:        conv.ID,
		Name:      conv.Name,
		IsGroup:   conv.IsGroup,
		CreatedAt: formattedCreatedAt,
		PhotoUrl:  conv.PhotoUrl,
	}

	// 7. Build and return the response.
	response := struct {
		Conversation Conversation       `json:"conversation"`
		Messages     []database.Message `json:"messages"`
	}{
		Conversation: apiConv,
		Messages:     messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding getConversation response: %v", err)
	}
}
