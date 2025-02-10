<!-- src/views/UserList.vue -->
<template>
    <div class="user-list-view">
      <h1>User Directory</h1>
      <button @click="refreshUsers">Refresh</button>
      <div v-if="error" class="error">{{ error }}</div>
      <ul>
        <li v-for="user in users" :key="user.id" class="user-item">
          <img :src="user.photoUrl || defaultPhoto" alt="User Photo" class="user-photo" />
          <span class="username">{{ user.username }}</span>
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { listUsers } from '@/services/api.js';
  
  export default {
    name: 'UserList',
    setup() {
      const users = ref([]);
      const error = ref("");
      const defaultPhoto = "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";
  
      async function refreshUsers() {
        error.value = "";
        try {
          const response = await listUsers();
          users.value = response.data.users;
        } catch (err) {
          error.value = "Failed to load users.";
          console.error(err);
        }
      }
  
      onMounted(() => {
        refreshUsers();
      });
  
      return { users, error, defaultPhoto, refreshUsers };
    },
  };
  </script>
  
  <style scoped>
  .user-list-view {
    padding: 20px;
    max-width: 800px;
    margin: 0 auto;
  }
  .user-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    border-bottom: 1px solid #ccc;
  }
  .user-photo {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    object-fit: cover;
  }
  .username {
    font-size: 1.2rem;
  }
  .error {
    color: red;
  }
  </style>
  