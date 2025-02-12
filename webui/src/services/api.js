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
// Replace the old sendMessage function with:
export function sendMessage(payload) {
  return axios.post('/messages', payload);
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
  return axios.post('/groups/create', data);
}

export function addUserToGroup(groupId, targetUserId) {
  return axios.post(`/groups/${groupId}/members`, { userId: targetUserId });
}

export function leaveGroup(groupId) {
  return axios.delete(`/groups/${groupId}/leave`);
}

export function setGroupName(groupId, newName) {
  return axios.put(`/groups/${groupId}/name`, { newName });
}

export function setGroupPhoto(groupId, photoUrl) {
  return axios.put(`/groups/${groupId}/photo`, { photoUrl });
}


