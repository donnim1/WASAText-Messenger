package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/donnim1/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

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

// getMyConversations retrieves all the conversations for the authenticated user.
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Get the authenticated user ID from the Authorization header.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Fetch conversations from the database using the authenticated user ID.
	conversations, err := rt.db.GetConversationsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve conversations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Convert database.Conversation to the API's Conversation type.
	var apiConversations []Conversation
	for _, conv := range conversations {
		apiConversations = append(apiConversations, Conversation{
			ID:        conv.ID,
			Name:      conv.Name,
			IsGroup:   conv.IsGroup,
			CreatedAt: conv.CreatedAt, // Make sure this field exists in your db.Conversation.
		})
	}

	// 4. Return the list of conversations.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(getMyConversationsResponse{
		Conversations: apiConversations,
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// getConversation handles GET requests to /conversations/:conversationId.
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 1. Validate the Authorization header.
	_, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Extract the conversation ID from the URL parameters.
	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Conversation ID is required", http.StatusBadRequest)
		return
	}

	// 3. Call the database function to get conversation details and messages.
	conv, messages, err := rt.db.GetConversation(conversationID)
	if err != nil {
		http.Error(w, "Failed to retrieve conversation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Check if conversation was found.
	if conv == nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// 5. Build the response.
	response := struct {
		Conversation *database.Conversation `json:"conversation"`
		Messages     []database.Message     `json:"messages"`
	}{
		Conversation: conv,
		Messages:     messages,
	}

	// 6. Return the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
