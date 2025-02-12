<!-- filepath: /c:/Users/HP/OneDrive/Desktop/2ND SEM/WASAText/webui/src/App.vue -->
<template>
  <div id="app" class="app-container">
    <div v-if="!userID">
      <Login @loggedIn="handleLogin" />
    </div>
    <div v-else class="authenticated-container">
      <aside class="side-nav">
        <div class="logo">WASATEXT</div>
        <nav>
          <router-link to="/myprofile" class="nav-link">My Profile</router-link>
          <router-link to="/users" class="nav-link">Contacts</router-link>
          <router-link to="/groups" class="nav-link">Groups</router-link>
          <router-link to="/chats" class="nav-link">Conversations</router-link>
        </nav>
        <button class="logout-btn" @click="logout">Logout</button>
      </aside>
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
}

.authenticated-container {
  flex: 1;
  display: flex;
}

.side-nav {
  width: 220px;
  background-color: #343a40;
  color: #fff;
  padding: 1rem;
  display: flex;
  flex-direction: column;
}

.side-nav .logo {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  text-align: center;
}

.side-nav nav {
  flex-grow: 1;
}

.nav-link {
  display: block;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
  color: #ccc;
  text-decoration: none;
}

.nav-link:hover {
  background-color: #495057;
  color: #fff;
}

.logout-btn {
  background-color: transparent;
  color: #fff;
  border: 1px solid #fff;
  padding: 0.5rem;
  cursor: pointer;
  margin-top: auto;
}

.main-content {
  flex-grow: 1;
  padding: 1rem;
  overflow-y: auto;
}
</style>