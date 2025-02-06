package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/donnim1/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// GetMyConversationsRequest defines the expected request for fetching a user's conversations.
type getMyConversationsRequest struct {
	UserID string `json:"id"` // User identifier (in real scenario, fetched from JWT or session)
}

// GetMyConversationsResponse defines the structure of the response after fetching the conversations.
type getMyConversationsResponse struct {
	Conversations []Conversation `json:"conversations"`
}

// Conversation defines a structure for each conversation in the response.
type Conversation struct {
	ID        string `json:"id"`         // Conversation ID
	Name      string `json:"name"`       // Name of the conversation (group name or null for private chats)
	IsGroup   bool   `json:"is_group"`   // True if it's a group chat, otherwise false
	CreatedAt string `json:"created_at"` // Timestamp of when the conversation was created
}

// getMyConversations retrieves all the conversations for a given user.
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Extract userID from request (it would typically be from JWT or session)
	var req getMyConversationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 2. Fetch conversations from the database
	conversations, err := rt.db.GetConversationsByUserID(req.UserID)
	if err != nil {
		http.Error(w, "Failed to retrieve conversations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Convert database.Conversation to api.Conversation
	var apiConversations []Conversation
	for _, conv := range conversations {
		apiConversations = append(apiConversations, Conversation{
			ID:        conv.ID,
			Name:      conv.Name,
			IsGroup:   conv.IsGroup,
			CreatedAt: conv.CreatedAt, // Ensure this is part of the db.Conversation
		})
	}

	// 4. Return the list of conversations
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(getMyConversationsResponse{
		Conversations: apiConversations,
	}); err != nil {
		// Log the error. The response is already sent, so this is only for debugging.
		log.Printf("Error encoding response: %v", err)
	}

}

// getConversation handles GET requests to /conversations/:conversationId.
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	// Call the database function to get conversation details and messages.
	conv, messages, err := rt.db.GetConversation(conversationID)
	if err != nil {
		http.Error(w, "Failed to retrieve conversation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if conversation was found.
	if conv == nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// Build the response.
	response := struct {
		Conversation *database.Conversation `json:"conversation"`
		Messages     []database.Message     `json:"messages"`
	}{
		Conversation: conv,
		Messages:     messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error if encoding fails (the response is already written).
		// In a production system, you might want to handle this differently.
		log.Printf("Error encoding response: %v", err)
	}
}
