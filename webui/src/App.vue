<template>
  <div id="app">
    <Login v-if="!userID" @loggedIn="handleLogin" />
    <div v-else>
      <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
        <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASATEXT</a>
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
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import Login from './components/Login.vue';

const userID = ref(localStorage.getItem('userID') || '');
const router = useRouter();

function handleLogin(id) {
  localStorage.setItem('userID', id);
  // Optionally, also clear profile details so the new login refreshes them.
  localStorage.removeItem('username');
  localStorage.removeItem('photoUrl');
  userID.value = id;
  router.push('/');
}

function logout() {
  localStorage.removeItem('userID');
  localStorage.removeItem('username');
  localStorage.removeItem('photoUrl');
  userID.value = '';
  router.push('/login');
}
</script>

<style>
/* Global styles (if any) */
</style>  
