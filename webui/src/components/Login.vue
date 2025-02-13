<template>
  <div class="login-view">
    <div class="login-container">
      <div class="login-card">
        <!-- Logo and Title -->
        <div class="brand-section">
          <div class="logo-container">
            <img src="@/assets/green-phone.jpg" alt="WASA TEXT Logo" class="logo" />
          </div>
          <h1 class="brand-title">WASA TEXT</h1>
          <p class="brand-subtitle">Connect with your friends effortlessly</p>
        </div>

        <!-- Login Form -->
        <form @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label for="username">Username</label>
            <input
              id="username"
              v-model="username"
              type="text"
              placeholder="Enter your username"
              required
              minlength="3"
              maxlength="16"
              class="form-input"
            />
          </div>
          
          <button type="submit" class="login-button">
            Login
          </button>
        </form>

        <!-- Error Message -->
        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <!-- Features -->
        <div class="features-section">
          <div class="feature-item">
            <span class="feature-icon">💬</span>
            <p>Instant Messaging</p>
          </div>
          <div class="feature-item">
            <span class="feature-icon">🔒</span>
            <p>Secure Chat</p>
          </div>
          <div class="feature-item">
            <span class="feature-icon">👥</span>
            <p>User Friendly</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, defineEmits } from "vue";
import { login } from "@/services/api.js";

export default {
  name: "Login",
  setup(_, { emit }) {
    const username = ref("");
    const errorMessage = ref("");

    async function handleLogin() {
      errorMessage.value = "";
      try {
        const response = await login(username.value);
        // Save user data into localStorage
        localStorage.setItem("userID", response.identifier);
        localStorage.setItem("username", response.username);
        localStorage.setItem("photoUrl", response.photoUrl || "");
        // Emit the loggedIn event so App.vue can update its state and route accordingly.
        emit("loggedIn", response.identifier, response.username, response.photoUrl);
      } catch (error) {
        errorMessage.value = typeof error === "string"
          ? error
          : "Login failed. Please try again.";
      }
    }

    return {
      username,
      errorMessage,
      handleLogin,
    };
  }
};
</script>

<style scoped>
.login-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #4dabf7 0%, #3c99e6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 440px;
}

.login-card {
  background-color: #ffffff;
  border-radius: 16px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.brand-section {
  text-align: center;
  margin-bottom: 40px;
}

.logo-container {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background-color: #f8f9fa;
  border-radius: 20px;
  padding: 15px;
}

.logo {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.brand-title {
  font-size: 2.4rem;
  font-weight: 700;
  color: #212529;
  margin: 0 0 10px 0;
}

.brand-subtitle {
  color: #868e96;
  font-size: 1rem;
  margin: 0;
}

.login-form {
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #495057;
  font-weight: 500;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: #4dabf7;
  box-shadow: 0 0 0 3px rgba(77, 171, 247, 0.1);
}

.login-button {
  width: 100%;
  padding: 14px;
  background-color: #4dabf7;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.login-button:hover {
  background-color: #3c99e6;
}

.error-message {
  text-align: center;
  color: #dc3545;
  background-color: #f8d7da;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 0.9rem;
}

.features-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-top: 30px;
  padding-top: 30px;
  border-top: 1px solid #e9ecef;
}

.feature-item {
  text-align: center;
}

.feature-icon {
  font-size: 2rem;
  margin-bottom: 8px;
  display: block;
}

.feature-item p {
  margin: 0;
  color: #495057;
  font-size: 0.9rem;
}
</style>