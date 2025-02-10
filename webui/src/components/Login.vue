<template>
  <div class="login-container">
    <div class="login-header">
      <img src="@/assets/green-phone.jpg" alt="WASA TEXT Logo" class="logo" />
      <h1 class="title">WASA TEXT</h1>
      <p class="subtitle">Connect with your friends effortlessly and stay in touch!</p>
    </div>
    <form @submit.prevent="handleLogin" class="login-form">
      <input v-model="username" placeholder="Enter your username" required minlength="3" maxlength="16" class="form-control" />
      <button type="submit" class="btn btn-primary">Login</button>
    </form>
    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
  </div>
</template>

<script>
import { login } from "../services/api.js";

export default {
  data() {
    return {
      username: "",
      errorMessage: "",
    };
  },
  methods: {
    async handleLogin() {
      try {
        const { identifier, username, photoUrl } = await login(this.username);

        localStorage.setItem("userID", identifier);
        localStorage.setItem("username", username);
        localStorage.setItem("photoUrl", photoUrl || "");

        this.$emit("loggedIn", identifier, username, photoUrl);
      } catch (error) {
        this.errorMessage = error || "Login failed. Please try again.";
      }
    },
  },
};
</script>


<style scoped>
.login-container {
  max-width: 400px;
  margin: 80px auto;
  padding: 20px;
  border: 1px solid #eaeaea;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  background-color: #ffffff;
  text-align: center;
}
.login-header {
  margin-bottom: 30px;
}
.logo {
  width: 60px;
  height: 60px;
  margin-bottom: 10px;
}
.title {
  font-size: 2.2rem;
  color: #2c3e50;
  margin-bottom: 5px;
}
.subtitle {
  font-size: 1rem;
  color: #7f8c8d;
  margin-bottom: 20px;
}
.login-form input.form-control {
  width: 100%;
  padding: 10px 15px;
  margin-bottom: 15px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}
.login-form button.btn {
  width: 100%;
  padding: 10px;
  font-size: 1rem;
  border-radius: 4px;
  background-color: #27ae60;
  border: none;
  color: white;
  cursor: pointer;
}
.login-form button.btn:hover {
  background-color: #219150;
}
.error {
  color: red;
  margin-top: 15px;
}
</style>
