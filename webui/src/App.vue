<template>
  <div id="app">
    <!-- If user is not logged in, show the Login component -->
    <Login v-if="!userID" @loggedIn="handleLogin" />

    <!-- If user is logged in, show the main app along with a logout button -->
    <div v-else>
      <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
        <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASATEXT</a>
        <!-- Logout button -->
        <button class="btn btn-sm btn-outline-light" @click="logout">Logout</button>
      </header>

      <router-view />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Login from './components/Login.vue'

// Retrieve the userID from localStorage on load
const userID = ref(localStorage.getItem('userID') || '')
const router = useRouter()

// When a user logs in, store their ID and update our reactive variable
function handleLogin(id) {
  localStorage.setItem('userID', id)
  userID.value = id
}

// Logout clears the stored userID and navigates to the login route
function logout() {
  localStorage.removeItem('userID')
  userID.value = ''
  router.push('/login')
}
</script>

<style>
/* Global styles (if any) */
</style>
