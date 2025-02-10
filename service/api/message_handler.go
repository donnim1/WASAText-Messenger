package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// MessageRequest defines the request format for sending a message.
type MessageRequest struct {
	ReceiverID string `json:"receiverId"` // For private chats; leave empty for group messages.
	Content    string `json:"content"`    // The message text.
	IsGroup    bool   `json:"isGroup"`    // True if sending to a group.
	GroupID    string `json:"groupId"`    // Required if IsGroup is true.
}

// MessageResponse defines the response format.
type MessageResponse struct {
	MessageID string `json:"messageId"`
}

// sendMessage handles sending a message.
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Validate Authorization header.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Decode the request body.
	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	// Validate required fields.
	if req.Content == "" || (!req.IsGroup && req.ReceiverID == "") {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	// Call database function to send the message.
	messageID, err := rt.db.SendMessage(userID, req.ReceiverID, req.Content, req.IsGroup, req.GroupID)
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(MessageResponse{MessageID: messageID}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// forwardMessageRequest defines the payload for forwarding a message.
type forwardMessageRequest struct {
	TargetConversationID string `json:"targetConversationId"`
}

// forwardMessage handles forwarding a message.
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
	// Decode the request body.
	var req forwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	// Forward the message.
	newMessageID, err := rt.db.ForwardMessage(originalMessageID, req.TargetConversationID, userID)
	if err != nil {
		http.Error(w, "Failed to forward message: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the response.
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

// commentMessage handles adding a reaction to a message.
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
	// Remove the reaction.
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
	// Delete the message using authenticated user ID.
	err = rt.db.DeleteMessage(messageID, userID)
	if err != nil {
		http.Error(w, "Failed to delete message: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Message deleted successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
