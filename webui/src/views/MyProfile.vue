<template>
  <div class="profile-view">
    <div class="profile-container">
      <!-- Left Panel: Profile Image and Info -->
      <div class="profile-panel">
        <div class="panel-header">
          <h1>My Profile</h1>
        </div>
        <div class="profile-content">
          <div class="profile-image-container">
            <!-- Add unique key to force image refresh + timestamp to prevent caching -->
            <img :src="currentPhotoUrl || defaultPhoto" :key="imageKey" alt="Profile Photo" class="profile-image" />
          </div>
          <h2 class="profile-username">{{ currentUsername }}</h2>
        </div>
      </div>

      <!-- Right Panel: Update Forms -->
      <div class="update-panel">
        <div class="panel-header">
          <h2>Settings</h2>
        </div>
        <div class="forms-container">
          <form @submit.prevent="updateUsernameHandler" class="update-form">
            <div class="form-group">
              <label for="username-update">Update Username</label>
              <input
                id="username-update"
                v-model="newUsername"
                type="text"
                placeholder="Enter new username"
                required
                minlength="3"
                maxlength="16"
              />
            </div>
            <button type="submit">Update Username</button>
          </form>

          <form @submit.prevent="uploadPhotoHandler" class="update-form">
            <div class="form-group">
              <label for="photo-upload">Update Profile Photo</label>
              <input 
                id="photo-upload" 
                name="photo" 
                type="file" 
                accept="image/*" 
                @change="handleFileChange" 
                required 
                class="file-input"
              />
            </div>
            <button type="submit">Upload Photo</button>
          </form>

          <div v-if="message" class="message success">{{ message }}</div>
          <div v-if="error" class="message error">{{ error }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from "vue";
import { updateUsername, updatePhoto } from "@/services/api.js";

export default {
  name: "MyProfile",
  setup() {
    const defaultPhoto =
      "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";
    // Read current user info from localStorage
    const currentUsername = ref(localStorage.getItem("username") || "Username");
    const currentPhotoUrl = ref(localStorage.getItem("photoUrl") || "");
    
    // Add imageKey to force re-render when photo changes
    const imageKey = ref(Date.now());
    
    const newUsername = ref("");
    const message = ref("");
    const error = ref("");
    const selectedFile = ref(null);

    function updateUsernameHandler() {
      error.value = "";
      message.value = "";
      updateUsername(newUsername.value)
        .then(() => {
          currentUsername.value = newUsername.value;
          localStorage.setItem("username", newUsername.value);
          message.value = "Username updated successfully.";
          newUsername.value = "";
        })
        .catch(() => {
          error.value = "Username Already exists.";
        });
    }

    function handleFileChange(e) {
      const files = e.target.files;
      if (files && files[0]) {
        selectedFile.value = files[0];
      }
    }

    async function uploadPhotoHandler() {
      error.value = "";
      message.value = "";
      if (!selectedFile.value) {
        error.value = "Please select a file.";
        return;
      }
      const formData = new FormData();
      formData.append("photo", selectedFile.value);
      try {
        const response = await updatePhoto(formData);
        console.log("Photo update response:", response);
        
        if (response.data && response.data.photoUrl) {
          // Update the photo URL with a timestamp to prevent caching
          const photoUrl = response.data.photoUrl;
          const updatedUrl = photoUrl.includes('?') 
            ? `${photoUrl}&t=${Date.now()}` 
            : `${photoUrl}?t=${Date.now()}`;
          
          currentPhotoUrl.value = updatedUrl;
          localStorage.setItem("photoUrl", photoUrl); // Store without timestamp
          
          // Update image key to force re-rendering of the image
          imageKey.value = Date.now();
          
          message.value = "Profile photo updated successfully.";
          console.log("Updated photo URL:", updatedUrl);
        } else {
          error.value = "Invalid response format from server.";
          console.error("Invalid response format:", response);
        }
      } catch (err) {
        error.value = "Failed to update profile photo.";
        console.error("Photo update error:", err);
      }
    }
    
    // Function to refresh profile data periodically
    function refreshProfile() {
      const storedUsername = localStorage.getItem("username");
      const storedPhotoUrl = localStorage.getItem("photoUrl");
      
      if (storedUsername && storedUsername !== currentUsername.value) {
        currentUsername.value = storedUsername;
      }
      
      if (storedPhotoUrl && storedPhotoUrl !== currentPhotoUrl.value) {
        // Add timestamp to prevent browser caching
        const updatedUrl = storedPhotoUrl.includes('?') 
          ? `${storedPhotoUrl}&t=${Date.now()}` 
          : `${storedPhotoUrl}?t=${Date.now()}`;
        
        currentPhotoUrl.value = updatedUrl;
        imageKey.value = Date.now(); // Force image to re-render
      }
    }
    
    // Set up auto-refresh
    let refreshInterval;
    
    onMounted(() => {
      // Auto-refresh every 2 seconds
      refreshInterval = setInterval(refreshProfile, 2000);
    });
    
    onUnmounted(() => {
      if (refreshInterval) {
        clearInterval(refreshInterval);
      }
    });

    return {
      currentUsername,
      currentPhotoUrl,
      defaultPhoto,
      imageKey,
      newUsername,
      message,
      error,
      updateUsernameHandler,
      handleFileChange,
      uploadPhotoHandler,
    };
  },
};
</script>

<style scoped>
.profile-view {
  height: 100vh;
  background-color: #f8f9fa;
  padding: 20px;
}

.profile-container {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 1px;
  height: 100%;
  max-width: 1000px;
  margin: 0 auto;
  background-color: #e9ecef;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.profile-panel, .update-panel {
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.panel-header {
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
  background-color: #ffffff;
}

.panel-header h1, .panel-header h2 {
  margin: 0;
  color: #212529;
  font-size: 1.5rem;
  font-weight: 600;
}

.profile-content {
  padding: 30px;
  text-align: center;
}

.profile-image-container {
  margin-bottom: 20px;
}

.profile-image {
  width: 180px;
  height: 180px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #4dabf7;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.profile-username {
  font-size: 1.8rem;
  color: #212529;
  margin: 0;
}

.forms-container {
  padding: 20px;
  overflow-y: auto;
}

.update-form {
  background-color: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #495057;
  font-weight: 500;
}

.update-form input[type="text"] {
  width: 100%;
  padding: 12px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  font-size: 0.95rem;
  transition: border-color 0.2s ease;
}

.update-form input[type="text"]:focus {
  outline: none;
  border-color: #4dabf7;
  box-shadow: 0 0 0 3px rgba(77, 171, 247, 0.1);
}

.file-input {
  width: 100%;
  padding: 8px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background-color: #ffffff;
}

.update-form button {
  width: 100%;
  padding: 12px;
  background-color: #4dabf7;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.update-form button:hover {
  background-color: #3c99e6;
}

.message {
  padding: 12px;
  border-radius: 6px;
  margin-top: 15px;
  text-align: center;
}

.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}
</style>