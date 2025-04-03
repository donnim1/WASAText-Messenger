package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"github.com/donnim1/WASAText/service/api/reqcontext"

	"github.com/donnim1/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

// groupListResponse defines the JSON response for listing groups.
type groupListResponse struct {
	Groups []Conversation `json:"groups"`
}

// listGroups handles GET /groups and returns all groups the authenticated user is a member of.
func (rt *_router) listGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	// Validate the Authorization header.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve groups from the database.
	groups, err := rt.db.GetGroupsByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve groups: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert []database.Conversation to []Conversation.
	var convs []Conversation
	for _, group := range groups {
		convs = append(convs, convertDBConversationToConversation(group))
		// If using sql.Rows elsewhere, remember to call rows.Err() after the loop.
	}

	// Process conversations and format the response.
	var apiConversations []Conversation
	for _, conv := range convs {
		// For private chats (is_group false), if the conversation name is empty,
		// try to get the chat partner's username.
		if conv.Name == "" && !conv.IsGroup {
			partner, err := rt.db.GetChatPartner(conv.ID, userID)
			if err == nil && partner != nil {
				conv.Name = partner.Username
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
			LastMessageContent: conv.LastMessageContent, // adjust if needed
			LastMessageSentAt:  conv.LastMessageSentAt,  // adjust if needed
			Members:            conv.Members,            // include the members from the DB
		})
	}

	// Return the groups as JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groupListResponse{Groups: apiConversations}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// createGroupRequest defines the expected JSON payload for creating a group.
type createGroupRequest struct {
	GroupName  string `json:"groupName"`  // Group name is required.
	GroupPhoto string `json:"groupPhoto"` // Optional group photo URL.
}

// createGroupResponse defines the response for group creation.
type createGroupResponse struct {
	GroupID string `json:"groupId"`
}

// createGroup handles POST /groups/create to create a new group.
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	creatorID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body.
	var req createGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if req.GroupName == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	groupID, err := rt.db.CreateGroup(creatorID, req.GroupName, req.GroupPhoto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create group: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createGroupResponse{GroupID: groupID}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// addToGroupRequest defines the required fields for adding a user to a group.
type addToGroupRequest struct {
	GroupID  string `json:"groupId"`
	Username string `json:"username"` // Required: target user's username.
}

// addToGroup handles POST /groups/add to add a user to a group.
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	_, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	var req addToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.GroupID == "" {
		req.GroupID = groupID
	}
	if req.GroupID == "" || req.Username == "" {
		http.Error(w, "Group ID and Username are required", http.StatusBadRequest)
		return
	}

	user, err := rt.db.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	if err := rt.db.AddToGroup(req.GroupID, user.ID); err != nil {
		http.Error(w, "Failed to add user to group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User added to group successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// leaveGroupRequest defines the payload for leaving a group.
type leaveGroupRequest struct {
	GroupID string `json:"groupId"`
}

// leaveGroup handles DELETE /groups/leave to remove a user from a group.
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := ps.ByName("groupId")
	var req leaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && groupID == "" {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if req.GroupID == "" {
		req.GroupID = groupID
	}
	if req.GroupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}
	if err := rt.db.LeaveGroup(req.GroupID, userID); err != nil {
		http.Error(w, "Failed to leave group: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// setGroupName handles PUT /groups/name to update a group's name.
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	_, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	var payload struct {
		NewName string `json:"newName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if payload.NewName == "" {
		http.Error(w, "New name is required", http.StatusBadRequest)
		return
	}

	if err := rt.db.SetGroupName(groupID, payload.NewName); err != nil {
		http.Error(w, "Failed to update group name: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Group name updated successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// convertDBConversationToConversation converts a database.Conversation to an API Conversation.
func convertDBConversationToConversation(dbConv database.Conversation) Conversation {
	return Conversation{
		ID:       dbConv.ID,
		Name:     dbConv.Name,
		PhotoUrl: dbConv.PhotoUrl,
		Members:  dbConv.Members,
		// Add other fields as needed.
	}
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	_, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := ps.ByName("groupId")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Attempt file upload: check if a file with key "photo" is provided.
	file, header, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()

		uploadDir := "uploads/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Printf("❌ Failed to create upload directory: %v", err)
			http.Error(w, "Server error: cannot create upload directory", http.StatusInternalServerError)
			return
		}

		filePath := filepath.Join(uploadDir, header.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			log.Printf("❌ Failed to create file: %v", err)
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			log.Printf("❌ Failed to write file: %v", err)
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		photoUrl := fmt.Sprintf("/%s", filepath.ToSlash(filePath))
		if err := rt.db.SetGroupPhoto(groupID, photoUrl); err != nil {
			http.Error(w, "Failed to update group photo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"photoUrl": photoUrl}); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
		}
		log.Println("✅ Group photo successfully updated (via file upload):", photoUrl)
		return
	}

	// Fallback: JSON-based photo URL update.
	var payload struct {
		PhotoUrl string `json:"photoUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if payload.PhotoUrl == "" {
		http.Error(w, "Photo URL is required", http.StatusBadRequest)
		return
	}

	if err := rt.db.SetGroupPhoto(groupID, payload.PhotoUrl); err != nil {
		http.Error(w, "Failed to update group photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"photoUrl": payload.PhotoUrl}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Println("✅ Group photo successfully updated (via URL):", payload.PhotoUrl)
}
