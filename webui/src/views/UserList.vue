<template>
  <div class="user-list-view">
    <h1>Contacts</h1>
    <button @click="refreshUsers" class="refresh-button">Refresh Contacts</button>
    <div v-if="error" class="error">{{ error }}</div>
    <ul class="contact-list">
      <li v-for="user in users" :key="user.id" class="contact-item">
        <img :src="user.photoUrl || defaultPhoto" alt="User Photo" class="contact-photo" />
        <div class="contact-info">
          <span class="contact-name">{{ user.username }}</span>
          <button class="chat-button" @click="openChatWithUser(user)">Chat</button>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { listUsers } from '@/services/api.js';
import { useRouter } from 'vue-router';

export default {
  name: 'UserList',
  setup() {
    const users = ref([]);
    const error = ref("");
    const defaultPhoto = "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";
    const router = useRouter();

    async function refreshUsers() {
      error.value = "";
      try {
        const response = await listUsers();
        users.value = response.data.users;
      } catch (err) {
        error.value = "Failed to load contacts.";
        console.error(err);
      }
    }

    function openChatWithUser(user) {
      // Navigate to ChatView (ensure your router has a named route 'ChatView')
      router.push({ name: 'ChatView', query: { receiverId: user.id } });
    }

    onMounted(() => {
      refreshUsers();
    });

    return { users, error, defaultPhoto, refreshUsers, openChatWithUser };
  },
};
</script>

<style scoped>
.user-list-view {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
  background-color: #f8f9fa;
}

.refresh-button {
  background-color: #27ae60;
  color: white;
  padding: 8px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-bottom: 10px;
}

.refresh-button:hover {
  background-color: #219150;
}

.contact-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px;
  border-bottom: 1px solid #ddd;
}

.contact-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.contact-info {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.contact-name {
  font-size: 1.2rem;
  font-weight: 500;
}

.chat-button {
  background-color: #0275d8;
  color: white;
  border: none;
  padding: 6px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.chat-button:hover {
  background-color: #025aa5;
}

.error {
  color: red;
  margin-top: 10px;
}
</style>
