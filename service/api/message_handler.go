package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// MessageRequest defines the request format for sending a message.
type MessageRequest struct {
	SenderID   string `json:"senderId"`   // Who is sending the message
	ReceiverID string `json:"receiverId"` // Who is receiving (for private chats)
	Content    string `json:"content"`    // Message content
	IsGroup    bool   `json:"isGroup"`    // True for group messages
	GroupID    string `json:"groupId"`    // Group ID if sending in a group
}

// MessageResponse defines the response format.
type MessageResponse struct {
	MessageID string `json:"messageId"`
}

// sendMessage handles sending a message and auto-creates a conversation if needed.
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Parse request body
	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 2. Validate request
	if req.SenderID == "" || req.Content == "" || (!req.IsGroup && req.ReceiverID == "") {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 3. Call database function to send message
	messageID, err := rt.db.SendMessage(req.SenderID, req.ReceiverID, req.Content, req.IsGroup, req.GroupID)
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(MessageResponse{
		MessageID: messageID,
	})
}

// forwardMessageRequest defines the payload for forwarding a message.
type forwardMessageRequest struct {
	TargetConversationID string `json:"targetConversationId"`
	SenderID             string `json:"senderId"`
}

// forwardMessage handles POST /message/:messageId/forward to forward an existing message.
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// Forward the message.
	newMessageID, err := rt.db.ForwardMessage(originalMessageID, req.TargetConversationID, req.SenderID)
	if err != nil {
		http.Error(w, "Failed to forward message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"messageId": newMessageID})
}

// commentMessageRequest defines the payload for commenting (reacting) on a message.
type commentMessageRequest struct {
	UserID   string `json:"userId"`
	Reaction string `json:"reaction"`
}

// commentMessage handles POST /message/:messageId/comment to add a reaction.
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

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
	err := rt.db.CommentMessage(messageID, req.UserID, req.Reaction)
	if err != nil {
		http.Error(w, "Failed to add comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
}

// uncommentMessage handles DELETE /message/:messageId/comment to remove a reaction.
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	// Assume the user ID is provided as a query parameter for simplicity.
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	err := rt.db.UncommentMessage(messageID, userID)
	if err != nil {
		http.Error(w, "Failed to remove comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment removed successfully"})
}

// deleteMessage handles DELETE /message/:messageId to delete a message.
type deleteMessageRequest struct {
	SenderID string `json:"senderId"`
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	var req deleteMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	err := rt.db.DeleteMessage(messageID, req.SenderID)
	if err != nil {
		http.Error(w, "Failed to delete message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message deleted successfully"})
}
