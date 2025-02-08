<template>
	<div class="profile-container">
	  <!-- Sidebar: Displays current user's photo and username -->
	  <div class="profile-sidebar">
		<img :src="currentPhotoUrl || defaultPhoto" alt="Profile Photo" class="profile-image" />
		<h2>{{ currentUsername }}</h2>
	  </div>
	  <!-- Main Section: Forms to update username and photo (no preview here) -->
	  <div class="profile-main">
		<h1>Update Profile</h1>
		<form @submit.prevent="updateUsernameHandler" class="username-form">
		  <label for="username">New Username:</label>
		  <input id="username" v-model="newUsername" placeholder="Enter new username" required minlength="3" maxlength="16" />
		  <button type="submit">Update Username</button>
		</form>
		<form @submit.prevent="updatePhotoHandler" class="photo-form">
		  <label for="photo">New Profile Photo URL:</label>
		  <input id="photo" v-model="newPhotoUrl" placeholder="Enter photo URL" required />
		  <button type="submit">Update Photo</button>
		</form>
		<div v-if="message" class="message">{{ message }}</div>
		<div v-if="error" class="error">{{ error }}</div>
	  </div>
	</div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { updateUsername, updatePhoto } from '@/services/api.js';
  
  export default {
	name: 'MyProfileView',
	setup() {
	  // Load current profile info from localStorage
	  const currentUsername = ref(localStorage.getItem('username') || 'Anonymous');
	  const currentPhotoUrl = ref(localStorage.getItem('photoUrl') || '');
	  const defaultPhoto = ref("https://via.placeholder.com/150?text=No+Photo");
  
	  // Variables for the update forms
	  const newUsername = ref('');
	  const newPhotoUrl = ref('');
	  const message = ref('');
	  const error = ref('');
  
	  async function updateUsernameHandler() {
		message.value = "";
		error.value = "";
		try {
		  const response = await updateUsername(newUsername.value);
		  message.value = response.data.message;
		  // Update sidebar and localStorage
		  currentUsername.value = newUsername.value;
		  localStorage.setItem('username', newUsername.value);
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
		  // Update sidebar and localStorage
		  currentPhotoUrl.value = newPhotoUrl.value;
		  localStorage.setItem('photoUrl', newPhotoUrl.value);
		} catch (err) {
		  error.value = "Failed to update photo";
		  console.error(err);
		}
	  }
  
	  onMounted(() => {
		// Reload profile data when component mounts
		currentUsername.value = localStorage.getItem('username') || 'Anonymous';
		currentPhotoUrl.value = localStorage.getItem('photoUrl') || '';
	  });
  
	  return {
		currentUsername,
		currentPhotoUrl,
		defaultPhoto,
		newUsername,
		newPhotoUrl,
		message,
		error,
		updateUsernameHandler,
		updatePhotoHandler,
	  };
	},
  };
  </script>
  
  <style scoped>
  .profile-container {
	display: flex;
	gap: 20px;
	padding: 20px;
	background: #e8f5e9; /* Light green background */
	border-radius: 8px;
	max-width: 900px;
	margin: 20px auto;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .profile-sidebar {
	flex: 0 0 250px;
	text-align: center;
	padding: 20px;
	background: #66bb6a; /* Medium green */
	border-radius: 8px;
	color: white;
  }
  
  .profile-sidebar .profile-image {
	width: 150px;
	height: 150px;
	object-fit: cover;
	border-radius: 50%;
	border: 2px solid white;
	margin-bottom: 10px;
  }
  
  .profile-main {
	flex: 1;
	padding: 20px;
	background: white;
	border-radius: 8px;
  }
  
  .username-form,
  .photo-form {
	margin-bottom: 20px;
  }
  
  label {
	display: block;
	font-weight: bold;
	margin-bottom: 5px;
  }
  
  input {
	width: 100%;
	padding: 10px;
	margin-bottom: 10px;
	border: 1px solid #ccc;
	border-radius: 4px;
	font-size: 1rem;
  }
  
  button {
	background-color: #43a047; /* Green */
	color: white;
	padding: 10px 15px;
	border: none;
	border-radius: 4px;
	cursor: pointer;
	font-size: 1rem;
  }
  
  button:hover {
	background-color: #388e3c;
  }
  
  .message {
	color: green;
	margin-top: 10px;
	font-weight: bold;
  }
  
  .error {
	color: red;
	margin-top: 10px;
	font-weight: bold;
  }
  </style>
  