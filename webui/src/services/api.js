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

export async function getConversation(conversationId) {
  try {
    const response = await axios.get(`/conversations/${conversationId}`);
    return response;
  } catch (error) {
    throw error.response?.data || 'Failed to load conversation.';
  }
}

export function getConversationByReceiver(receiverId) {
  // Assumes the backend endpoint GET /conversations/for/:receiverId is implemented.
  return axios.get(`/conversationsfor/${receiverId}`);
}

export async function getConversations() {
  try {
    const response = await axios.get("/conversations"); // adjust endpoint if needed
    return response;
  } catch (error) {
    throw error.response?.data || new Error("Failed to load conversations.");
  }
}

// Messaging Endpoints
export function sendMessage({ conversationId, receiverId, content, isGroup, groupId, replyTo }) {
  return axios.post('/messages', { conversationId, receiverId, content, isGroup, groupId, replyTo });
}

export function forwardMessageApi(messageId, targetConversationId) {
  console.log(`API calling /messages/${messageId}/forward with:`, { targetConversationId });
  return axios.post(`/messages/${messageId}/forward`, {
    targetConversationId: targetConversationId  // Make sure the property name matches exactly
  });
}

export async function commentMessage(messageId, reaction) {
  return axios.post(`/messages/${messageId}/comments`, { reaction });
}

/**
 * Remove a reaction (uncomment) from a message.
 * Updated to call the correct endpoint.
 * @param {string} messageId - The ID of the message.
 * @returns {Promise} Axios promise.
 */
export async function uncommentMessage(messageId) {
  return axios.delete(`/messages/${messageId}/uncomment`);
}

export function deleteMessage(messageId) {
  return axios.delete(`/messages/${messageId}`);
}

// Group Management Endpoints
export function listUserGroups() {
  return axios.get('/groups');
}

export function createGroup(data) {
  return axios.post('/group', data);
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
  return axios.post(`/groups/${groupId}/members`, { username });
}

export async function uploadImage(formData) {
  return axios.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });
}

/**
 * Upload a new group photo.
 * @param {string} groupId - The group's ID.
 * @param {FormData} formData - FormData containing the photo file.
 * @returns {Promise} - Axios response with the updated photo URL.
 */
export async function uploadGroupImage(groupId, formData) {
  return axios.put(`/groups/${groupId}/photo`, formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
}

export function updateMessageStatus(messageId, status) {
  return axios.post(`/messages/${messageId}/status/${status}`);
}