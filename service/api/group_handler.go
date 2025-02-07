package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// addToGroupRequest defines the expected JSON payload for adding a user to a group.
// We no longer need to pass UserID in the body because we extract it from the Authorization header.
type addToGroupRequest struct {
	GroupID string `json:"groupId"`
}

// addToGroup handles POST /groups/add to add a user to a group.
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Validate the Authorization header and get the authenticated user ID.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req addToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if req.GroupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}
	// Use the authenticated user ID instead of a request body field.
	if err := rt.db.AddToGroup(req.GroupID, userID); err != nil {
		http.Error(w, "Failed to add user to group: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User added to group successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// leaveGroupRequest defines the expected JSON payload for leaving a group.
// We no longer need to pass UserID in the body because we extract it from the Authorization header.
type leaveGroupRequest struct {
	GroupID string `json:"groupId"`
}

// leaveGroup handles DELETE /groups/leave to remove a user from a group.
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Validate the Authorization header.
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req leaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
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

// setGroupNameRequest defines the expected JSON payload for updating a group's name.
type setGroupNameRequest struct {
	GroupID string `json:"groupId"`
	NewName string `json:"newName"`
}

// setGroupName handles PUT /groups/name to update a group's name.
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Validate the Authorization header.
	_, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req setGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if req.GroupID == "" || req.NewName == "" {
		http.Error(w, "Group ID and new name are required", http.StatusBadRequest)
		return
	}
	if err := rt.db.SetGroupName(req.GroupID, req.NewName); err != nil {
		http.Error(w, "Failed to update group name: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Group name updated successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// setGroupPhotoRequest defines the expected JSON payload for updating a group's photo.
type setGroupPhotoRequest struct {
	GroupID  string `json:"groupId"`
	PhotoURL string `json:"photoUrl"`
}

// setGroupPhoto handles PUT /groups/photo to update a group's photo.
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Validate the Authorization header.
	_, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req setGroupPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if req.GroupID == "" || req.PhotoURL == "" {
		http.Error(w, "Group ID and photo URL are required", http.StatusBadRequest)
		return
	}
	if err := rt.db.SetGroupPhoto(req.GroupID, req.PhotoURL); err != nil {
		http.Error(w, "Failed to update group photo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Group photo updated successfully"}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
