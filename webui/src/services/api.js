import axios from './axios.js';

// Update the current user's username
export function updateUsername(newUsername) {
  return axios.put('/user/username', { newName: newUsername });
}

// Update the current user's profile photo
export function updatePhoto(newPhotoUrl) {
  return axios.put('/user/photo', { photoUrl: newPhotoUrl });
}

// List all users (for example, for a user directory)
export function listUsers() {
    return axios.get('/users');
  }

// Send a message (private or group)
export function sendMessage({ receiverId, content, isGroup, groupId }) {
  return axios.post('/message', { receiverId, content, isGroup, groupId });
}

// Forward a message
export function forwardMessage(messageId, targetConversationId) {
  return axios.post(`/message/${messageId}/forward`, { targetConversationId });
}

// Comment on a message
export function commentMessage(messageId, reaction) {
  return axios.post(`/message/${messageId}/comment`, { reaction });
}

// Remove (uncomment) a message comment
export function uncommentMessage(messageId) {
  return axios.delete(`/message/${messageId}/comment`);
}

// Delete a message
export function deleteMessage(messageId) {
  // Note: For DELETE requests with a body, we use "data" in the config.
  return axios.delete(`/message/${messageId}/delete`, { data: {} });
}

// ----- Group Endpoints ----- //

// List all groups the authenticated user is a member of.
export function listUserGroups() {
  return axios.get('/groups');
}

// For example, your other functions might be:
export function createGroup(data) {
  return axios.post('/groups/create', data);
}

export function addUserToGroup(groupId, targetUserId) {
  return axios.post('/groups/add', { groupId, userId: targetUserId });
}

export function leaveGroup(groupId) {
  return axios.delete('/groups/leave', { data: { groupId } });
}

// Update group name
export function setGroupName(groupId, newName) {
  return axios.put('/groups/name', { groupId, newName });
}

// Update group photo
export function setGroupPhoto(groupId, photoUrl) {
  return axios.put('/groups/photo', { groupId, photoUrl });
}
