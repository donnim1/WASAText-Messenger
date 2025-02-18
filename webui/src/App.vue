<template>
  <div id="app" class="app-container">
    <div v-if="!userID">
      <Login @loggedIn="handleLogin" />
    </div>
    <div v-else class="authenticated-container">
      <!-- Modern Sidebar -->
      <aside class="side-nav">
        <div class="nav-header">
          <div class="logo">
            <span class="logo-text">WASATEXT</span>
          </div>
        </div>
        
        <nav class="nav-links">
          <router-link to="/myprofile" class="nav-link">
            <span class="nav-icon">👤</span>
            <span class="nav-text">My Profile</span>
          </router-link>
          <router-link to="/users" class="nav-link">
            <span class="nav-icon">👥</span>
            <span class="nav-text">Contacts</span>
          </router-link>
          <router-link to="/groups" class="nav-link">
            <span class="nav-icon">🗣️</span>
            <span class="nav-text">Groups</span>
          </router-link>
          <router-link to="/chats" class="nav-link">
            <span class="nav-icon">💬</span>
            <span class="nav-text">Conversations</span>
          </router-link>
        </nav>
        
        <div class="nav-footer">
          <button class="logout-btn" @click="logout">
            <span class="nav-icon">🚪</span>
            <span class="nav-text">Logout</span>
          </button>
        </div>
      </aside>

      <!-- Main Content Area -->
      <main class="main-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import Login from "./components/Login.vue";

const userID = ref(localStorage.getItem("userID") || "");
const router = useRouter();

function handleLogin(id, name, photo) {
  localStorage.setItem("userID", id);
  localStorage.setItem("username", name);
  localStorage.setItem("photoUrl", photo || "");
  userID.value = id;
  router.push("/myprofile");
}

function logout() {
  localStorage.removeItem("userID");
  localStorage.removeItem("username");
  localStorage.removeItem("photoUrl");
  userID.value = "";
  router.push("/login");
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f8f9fa;
}

.authenticated-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* Sidebar Styles */
.side-nav {
  width: 260px;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #e9ecef;
  transition: width 0.3s ease;
}

.nav-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e9ecef;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-size: 1.5rem;
  font-weight: 700;
  color: #4dabf7;
  letter-spacing: 0.5px;
}

/* Navigation Links */
.nav-links {
  flex: 1;
  padding: 1rem 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.nav-link {
  display: flex;
  align-items: center;
  padding: 0.875rem 1rem;
  color: #495057;
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.nav-link:hover {
  background-color: #f8f9fa;
  color: #4dabf7;
}

.router-link-active {
  background-color: #e7f5ff;
  color: #4dabf7;
  font-weight: 500;
}

.nav-icon {
  font-size: 1.25rem;
  margin-right: 0.75rem;
  width: 24px;
  text-align: center;
}

.nav-text {
  font-size: 0.95rem;
}

/* Footer Section */
.nav-footer {
  padding: 1rem 0.75rem;
  border-top: 1px solid #e9ecef;
}

.logout-btn {
  width: 100%;
  display: flex;
  align-items: center;
  padding: 0.875rem 1rem;
  background-color: transparent;
  border: none;
  border-radius: 8px;
  color: #dc3545;
  cursor: pointer;
  transition: all 0.2s ease;
}

.logout-btn:hover {
  background-color: #fff5f5;
}

.logout-btn .nav-icon {
  color: #dc3545;
}

/* Main Content Area */
.main-content {
  flex: 1;
  overflow-y: auto;
  background-color: #f8f9fa;
  position: relative;
}

/* Responsive Design */
@media (max-width: 768px) {
  .side-nav {
    width: 80px;
  }

  .nav-text {
    display: none;
  }

  .logo-text {
    display: none;
  }

  .nav-link, .logout-btn {
    justify-content: center;
    padding: 0.875rem;
  }

  .nav-icon {
    margin-right: 0;
  }
}

@media (max-width: 480px) {
  .side-nav {
    width: 60px;
  }

  .nav-link, .logout-btn {
    padding: 0.75rem;
  }

  .nav-icon {
    font-size: 1.1rem;
  }
}
</style>