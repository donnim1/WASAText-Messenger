<template>
  <div id="app">
    <Login v-if="!userID" @loggedIn="handleLogin" />
    <div v-else>
      <header class="navbar navbar-dark bg-dark p-0 shadow">
        <a class="navbar-brand px-3 fs-6" >WASATEXT</a>
        <nav>
          
          <router-link to="/myprofile">My Profile</router-link>
          <router-link to="/users">User Directory</router-link>
          <router-link to="/groups">Groups</router-link>
          <router-link to="/chat">Chats</router-link>
        </nav>
        <button class="btn btn-sm btn-outline-light" @click="logout">Logout</button>
      </header>
      <router-view />
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
  router.push("/");
}
</script>
