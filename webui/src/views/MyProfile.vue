<template>
	<div class="profile-card">
	  <!-- Sidebar: Displays current user's photo and username -->
	  <div class="profile-sidebar">
		<img :src="currentPhotoUrl || defaultPhoto" alt="Profile Photo" class="profile-image" />
		<h2 class="profile-username">{{ currentUsername }}</h2>
	  </div>
	  <!-- Main Section: Forms to update username and photo -->
	  <div class="profile-main">
		<h1>Update Profile</h1>
		<form @submit.prevent="updateUsernameHandler" class="update-form">
		  <label for="username">Update Username</label>
		  <input
			id="username"
			v-model="newUsername"
			type="text"
			placeholder="Enter new username"
			required
			minlength="3"
			maxlength="16"
		  />
		  <button type="submit">Update Username</button>
		</form>
		<form @submit.prevent="updatePhotoHandler" class="update-form">
		  <label for="photo">Update Profile Photo URL</label>
		  <input
			id="photo"
			v-model="newPhotoUrl"
			type="url"
			placeholder="Enter new photo URL"
			required
		  />
		  <button type="submit">Update Photo</button>
		</form>
		<div v-if="message" class="success-message">{{ message }}</div>
		<div v-if="error" class="error-message">{{ error }}</div>
	  </div>
	</div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { updateUsername, updatePhoto } from '@/services/api.js';
  
  export default {
	name: 'MyProfileView',
	setup() {
	  // Read current profile info from localStorage; if not present, use defaults.
	  const currentUsername = ref(localStorage.getItem('username') || '');
	  const currentPhotoUrl = ref(localStorage.getItem('photoUrl') || '');
	  const defaultPhoto = ref("https://static.vecteezy.com/system/resources/thumbnails/009/734/564/small/default-avatar-profile-icon-of-social-media-user-vector.jpg");
  
	  // Fields for updating profile data.
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
		  // Update the current username and persist it.
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
		  // Update the current photo and persist it.
		  currentPhotoUrl.value = newPhotoUrl.value;
		  localStorage.setItem('photoUrl', newPhotoUrl.value);
		} catch (err) {
		  error.value = "Failed to update photo";
		  console.error(err);
		}
	  }
  
	  onMounted(() => {
		// Reload profile details from localStorage on component mount.
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
  .profile-card {
	max-width: 500px;
	margin: 40px auto;
	padding: 20px;
	border: 1px solid #ddd;
	border-radius: 8px;
	background: #fff;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	display: flex;
	gap: 20px;
  }
  
  .profile-sidebar {
	flex: 0 0 200px;
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
  
  .profile-username {
	font-size: 1.8rem;
	margin-top: 10px;
  }
  
  .profile-main {
	flex: 1;
	padding: 20px;
	background: #f9fbe7; /* Light green background */
	border-radius: 8px;
  }
  
  .update-form {
	margin-bottom: 20px;
	display: flex;
	flex-direction: column;
  }
  
  .update-form label {
	font-weight: bold;
	margin-bottom: 5px;
  }
  
  .update-form input {
	padding: 8px;
	border: 1px solid #ccc;
	border-radius: 4px;
	font-size: 1rem;
	margin-bottom: 10px;
  }
  
  .update-form button {
	padding: 10px;
	background-color: #43a047;
	color: white;
	border: none;
	border-radius: 4px;
	cursor: pointer;
	font-size: 1rem;
  }
  
  .update-form button:hover {
	background-color: #388e3c;
  }
  
  .success-message {
	color: green;
	font-weight: bold;
	text-align: center;
	margin-top: 10px;
  }
  
  .error-message {
	color: red;
	font-weight: bold;
	text-align: center;
	margin-top: 10px;
  }
  </style>
  