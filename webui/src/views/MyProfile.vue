<!-- filepath: /c:/Users/HP/OneDrive/Desktop/2ND SEM/WASAText/webui/src/views/MyProfile.vue -->
<template>
  <div class="profile-card">
    <!-- Sidebar with profile image and username -->
    <div class="profile-sidebar">
      <img :src="currentPhotoUrl || defaultPhoto" alt="Profile Photo" class="profile-image" />
      <h2 class="profile-username">{{ currentUsername }}</h2>
    </div>

    <!-- Main section with forms -->
    <div class="profile-main">
      <h1>Update Profile</h1>
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
      <form @submit.prevent="uploadPhotoHandler" class="update-form">
        <label for="photo-upload">Update Profile Photo</label>
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
    const currentUsername = ref(localStorage.getItem("username") || "Anonymous");
    const currentPhotoUrl = ref(localStorage.getItem("photoUrl") || "");
    const defaultPhoto = ref("https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg");

    const newUsername = ref("");
    const selectedPhotoFile = ref(null);
    const message = ref("");
    const error = ref("");

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

    function handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        selectedPhotoFile.value = file;
      }
    }

    async function uploadPhotoHandler() {
      message.value = "";
      error.value = "";
      if (!selectedPhotoFile.value) {
        error.value = "No photo selected";
        return;
      }
      try {
        const formData = new FormData();
        formData.append("photo", selectedPhotoFile.value);
        const response = await updatePhoto(formData);
        message.value = "Profile photo updated successfully";
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
  max-width: 600px;
  margin: 40px auto;
  display: flex;
  background-color: #2c2c2c;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.5);
  color: #e0e0e0;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.profile-sidebar {
  flex: 0 0 220px;
  background-color: #1a1a1a;
  padding: 30px;
  text-align: center;
}

.profile-sidebar .profile-image {
  width: 140px;
  height: 140px;
  object-fit: cover;
  border-radius: 50%;
  border: 3px solid #4caf50;
  margin-bottom: 15px;
}

.profile-username {
  font-size: 1.8rem;
  margin-top: 10px;
}

.profile-main {
  flex: 1;
  padding: 30px;
  background-color: #333;
}

.profile-main h1 {
  margin-bottom: 20px;
  font-size: 1.8rem;
  color: #4caf50;
}

.update-form {
  margin-bottom: 25px;
  display: flex;
  flex-direction: column;
}

.update-form label {
  margin-bottom: 8px;
  font-weight: bold;
  color: #bbb;
}

.update-form input {
  padding: 12px;
  border: 1px solid #555;
  border-radius: 6px;
  background-color: #444;
  color: #e0e0e0;
  margin-bottom: 15px;
  font-size: 1rem;
}

.update-form button {
  padding: 12px;
  background-color: #4caf50;
  border: none;
  border-radius: 6px;
  color: #fff;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.update-form button:hover {
  background-color: #43a047;
}

.success-message {
  color: #4caf50;
  text-align: center;
  margin-top: 15px;
}

.error-message {
  color: #f44336;
  text-align: center;
  margin-top: 15px;
}
</style>