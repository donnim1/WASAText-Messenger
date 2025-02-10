<template>
  <div id="app">
    <Login v-if="!userID" @loggedIn="handleLogin" />
    <div v-else>
      <header class="navbar navbar-dark bg-dark p-0 shadow">
        <a class="navbar-brand px-3 fs-6" href="#">WASATEXT</a>
        <nav>
          <router-link to="/myprofile">My Profile</router-link>
          <router-link to="/users">User Directory</router-link>
          <router-link to="/groups">Groups</router-link>
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
const username = ref(localStorage.getItem("username") || "Anonymous");
const photoUrl = ref(localStorage.getItem("photoUrl") || "");

const defaultPhoto = "https://via.placeholder.com/50?text=Photo";
const router = useRouter();

function handleLogin(id, name, photo) {
  localStorage.setItem("userID", id);
  localStorage.setItem("username", name);
  localStorage.setItem("photoUrl", photo || "");

  userID.value = id;
  username.value = name;
  photoUrl.value = photo || "";

  router.push("/myprofile");
}

function logout() {
  localStorage.removeItem("userID");
  localStorage.removeItem("username");
  localStorage.removeItem("photoUrl");

  userID.value = "";
  username.value = "Anonymous";
  photoUrl.value = defaultPhoto;

  router.push("/");
}
</script>
