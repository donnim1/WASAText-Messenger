<!-- src/views/MyProfileView.vue -->
<template>
	<div class="profile-view">
	  <h1>My Profile</h1>
	  <form @submit.prevent="updateUsername">
		<h2>Update Username</h2>
		<input v-model="newUsername" placeholder="Enter new username" required />
		<button type="submit">Update Username</button>
	  </form>
	  <form @submit.prevent="updatePhoto">
		<h2>Update Profile Photo</h2>
		<input v-model="newPhotoUrl" placeholder="Enter photo URL" required />
		<button type="submit">Update Photo</button>
	  </form>
	  <div v-if="message" class="message">{{ message }}</div>
	  <div v-if="error" class="error">{{ error }}</div>
	</div>
  </template>
  
  <script>
  import { ref } from 'vue';
  import { updateUsername, updatePhoto } from '@/services/api.js';
  
  export default {
	name: 'MyProfileView',
	setup() {
	  const newUsername = ref("");
	  const newPhotoUrl = ref("");
	  const message = ref("");
	  const error = ref("");
  
	  async function updateUsernameHandler() {
		message.value = "";
		error.value = "";
		try {
		  const response = await updateUsername(newUsername.value);
		  message.value = response.data.message;
		} catch (err) {
		  error.value = "Failed to update username";
		  console.error(err);
		}
	  }
  
	  async function updatePhotoHandler() {
		message.value = "";
		error.value = "";
		try {
		  const response = await updatePhoto(newPhotoUrl.value);
		  message.value = response.data.message;
		} catch (err) {
		  error.value = "Failed to update photo";
		  console.error(err);
		}
	  }
  
	  return {
		newUsername,
		newPhotoUrl,
		message,
		error,
		updateUsername: updateUsernameHandler,
		updatePhoto: updatePhotoHandler,
	  };
	},
  };
  </script>
  
  <style scoped>
  .message {
	color: green;
  }
  .error {
	color: red;
  }
  </style>
  