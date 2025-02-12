import axios from './axios.js';

// Login Function
export async function login(username) {
  try {
    const response = await axios.post('/session', { name: username });
    return response.data; // { identifier, username, photoUrl }
  } catch (error) {
    throw error.response?.data || 'Login failed. Please try again.';
  }
}

// Update Username
export function updateUsername(newUsername) {
  return axios.put('/user/username', { newName: newUsername });
}
 
// Update Profile Photo (using FormData)
export function updatePhoto(formData) {
  return axios.put('/user/photo', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });
}

// List all Users
export function listUsers() {
  return axios.get('/users');
}

// Conversations
export function getMyConversations() {
  return axios.get('/conversation/myconversations');
}

export function getConversation(conversationId) {
  return axios.get(`/conversations/${conversationId}`);
}

// Messaging Endpoints
export function sendMessage({ conversationId, receiverId, content, isGroup, groupId }) {
  // Now include conversationId (which may be empty if not yet created)
  return axios.post('/messages', { conversationId, receiverId, content, isGroup, groupId });
}

export function forwardMessage(messageId, targetConversationId) {
  return axios.post(`/messages/${messageId}/forward`, { targetConversationId });
}

export function commentMessage(messageId, reaction) {
  return axios.post(`/messages/${messageId}/comments`, { reaction });
}

export function uncommentMessage(messageId) {
  return axios.delete(`/messages/${messageId}/uncomment`);
}

export function deleteMessage(messageId) {
  return axios.delete(`/messages/${messageId}/delete`);
}

// Group Management Endpoints
export function listUserGroups() {
  return axios.get('/groups');
}

export function createGroup(data) {
  return axios.post('/groups', data);
}

export function addUserToGroup(groupId, targetUserId) {
  return axios.post(`/groups/${groupId}/members`, { userId: targetUserId });
}

export function leaveGroup(groupId) {
  // Updated endpoint to match backend change.
  return axios.delete(`/groups/${groupId}/leave`);
}

export function setGroupName(groupId, payload) {
  // Expects payload { newName: '...' }
  return axios.put(`/groups/${groupId}/name`, payload);
}

export function setGroupPhoto(groupId, payload) {
  // Expects payload { photoUrl: '...' }
  return axios.put(`/groups/${groupId}/photo`, payload);
}

export async function addUserToGroupByUsername(groupId, username) {
  // This example assumes your backend exposes an endpoint to add by username.
  return axios.post(`/groups/members/${groupId}/by-username`, { username });
}


