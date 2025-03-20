package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// MessageRequest defines the request format for sending a message.
type MessageRequest struct {
	ConversationID string `json:"conversationId"` // Conversation ID
	ReceiverID     string `json:"receiverId"`     // Receiver ID
	Content        string `json:"content"`
	IsGroup        bool   `json:"isGroup"`
	GroupID        string `json:"groupId"`           // Group ID
	ReplyTo        string `json:"replyTo,omitempty"` // Optional reply-to field
}

// MessageResponse defines the response format.
type MessageResponse struct {
	MessageID      string `json:"messageId"`
	ConversationID string `json:"conversationId"`
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Validate Authorization header.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode request body.
	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields.
	if req.Content == "" || (!req.IsGroup && req.ConversationID == "" && req.ReceiverID == "") {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Call the updated SendMessage function.
	messageID, conversationID, err := rt.db.SendMessage(userID, req.ReceiverID, req.Content, req.IsGroup, req.GroupID, req.ConversationID, req.ReplyTo)
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response with both messageId and conversationId.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		MessageID      string `json:"messageId"`
		ConversationID string `json:"conversationId"`
	}{
		MessageID:      messageID,
		ConversationID: conversationID,
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// forwardMessageRequest defines the payload for forwarding a message.
type forwardMessageRequest struct {
	TargetConversationID string `json:"targetConversationId"`
	// SenderID is removed because we use the token's userID.
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Validate Authorization header.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	originalMessageID := ps.ByName("messageId")
	if originalMessageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	var req forwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	newMessageID, err := rt.db.ForwardMessage(originalMessageID, req.TargetConversationID, userID)
	if err != nil {
		http.Error(w, "Failed to forward message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"messageId": newMessageID}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// commentMessageRequest defines the payload for commenting (reacting) on a message.
type commentMessageRequest struct {
	Reaction string `json:"reaction"`
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get authenticated user ID.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	// Decode the request body.
	var req commentMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if req.Reaction == "" {
		http.Error(w, "Reaction cannot be empty", http.StatusBadRequest)
		return
	}

	// Insert the reaction.
	err = rt.db.CommentMessage(messageID, userID, req.Reaction)
	if err != nil {
		http.Error(w, "Failed to add comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get authenticated user ID.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	err = rt.db.UncommentMessage(messageID, userID)
	if err != nil {
		http.Error(w, "Failed to remove comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Comment removed successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get authenticated user ID.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	// Use the authenticated userID directly for deletion.
	err = rt.db.DeleteMessage(messageID, userID)
	if err != nil {
		http.Error(w, "Failed to delete message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Message deleted successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (rt *_router) updateMessageStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get message ID and new status from URL parameters.
	messageID := ps.ByName("messageId")
	status := ps.ByName("status") // Expected values: "delivered" or "read"
	if messageID == "" || status == "" {
		http.Error(w, "Message ID and status are required", http.StatusBadRequest)
		return
	}

	// Retrieve the authenticated user's ID.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Update the message status in the database.
	if err := rt.db.UpdateMessageStatus(messageID, status, userID); err != nil {
		http.Error(w, "Failed to update message status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
