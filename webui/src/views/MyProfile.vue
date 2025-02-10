<template>
  <div class="profile-card">
    <!-- Sidebar: Displays current user's photo and username -->
    <div class="profile-sidebar">
      <img :src="currentPhotoUrl || defaultPhoto" alt="Profile Photo" class="profile-image" />
      <h2 class="profile-username">{{ currentUsername }}</h2>
    </div>

    <!-- Main Section: Forms to update username and upload a new profile photo -->
    <div class="profile-main">
      <h1>Update Profile</h1>

      <!-- Update Username Form -->
      <form @submit.prevent="updateUsernameHandler" class="update-form">
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
        <button type="submit">Update Username</button>
      </form>

      <!-- Update Profile Photo Form -->
      <form @submit.prevent="uploadPhotoHandler" class="update-form">
        <label for="photo-upload">Update Profile Photo</label>
        <!-- Note: Add a name attribute "photo" for proper file field recognition -->
        <input id="photo-upload" name="photo" type="file" accept="image/*" @change="handleFileChange" required />
        <button type="submit">Upload Photo</button>
      </form>

      <div v-if="message" class="success-message">{{ message }}</div>
      <div v-if="error" class="error-message">{{ error }}</div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { updateUsername, updatePhoto } from "@/services/api.js";

export default {
  name: "MyProfileView",
  setup() {
    // Load current profile info from localStorage.
    const currentUsername = ref(localStorage.getItem("username") || "Anonymous");
    const currentPhotoUrl = ref(localStorage.getItem("photoUrl") || "");
    const defaultPhoto = ref("https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg");

    // Update form fields.
    const newUsername = ref("");
    const selectedPhotoFile = ref(null);
    const message = ref("");
    const error = ref("");

    // Update username handler.
    async function updateUsernameHandler() {
      message.value = "";
      error.value = "";
      try {
        const response = await updateUsername(newUsername.value);
        message.value = response.data.message;
        currentUsername.value = newUsername.value;
        localStorage.setItem("username", newUsername.value);
      } catch (err) {
        error.value = "Failed to update username";
        console.error(err);
      }
    }

    // File change event: store the selected file.
    function handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        selectedPhotoFile.value = file;
      }
    }

    // Upload photo handler: create FormData and call updatePhoto API.
    async function uploadPhotoHandler() {
      message.value = "";
      error.value = "";
      if (!selectedPhotoFile.value) {
        error.value = "No photo selected";
        return;
      }
      try {
        // Wrap the file in FormData.
        const formData = new FormData();
        formData.append("photo", selectedPhotoFile.value);

        const response = await updatePhoto(formData);
        message.value = "Profile photo updated successfully";

        // Append a timestamp to avoid caching issues.
        const updatedPhotoUrl = response.data.photoUrl + `?t=${Date.now()}`;
        currentPhotoUrl.value = updatedPhotoUrl;
        localStorage.setItem("photoUrl", updatedPhotoUrl);
      } catch (err) {
        error.value = "Failed to update photo";
        console.error(err);
      }
    }

    onMounted(() => {
      currentUsername.value = localStorage.getItem("username") || "Anonymous";
      currentPhotoUrl.value = localStorage.getItem("photoUrl") || "";
    });

    return {
      currentUsername,
      currentPhotoUrl,
      defaultPhoto,
      newUsername,
      selectedPhotoFile,
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
  background: #66bb6a;
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
  background: #f9fbe7;
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
