import axios from './axios.js';

// Update the current user's username
export function updateUsername(newUsername) {
  return axios.put('/user/username', { newName: newUsername });
}

// Update the current user's profile photo
export function updatePhoto(newPhotoUrl) {
  return axios.put('/user/photo', { photoUrl: newPhotoUrl });
}

// (Optional) Add additional API functions here as needed
