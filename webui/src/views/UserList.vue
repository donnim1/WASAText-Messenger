<template>
  <div class="user-list-view">
    <h1>Contacts</h1>
    <!-- Search Field -->
    <div class="search-container">
      <input
        v-model="searchQuery"
        placeholder="Search contacts..."
        class="search-input"
      />
    </div>
    <!-- Refresh Button -->
    <button @click="refreshUsers" class="refresh-button">Refresh Contacts</button>
    <!-- Error Message -->
    <div v-if="error" class="error">{{ error }}</div>
    <!-- Contacts List -->
    <ul class="contact-list">
      <li v-for="user in filteredUsers" :key="user.id" class="contact-item">
        <img
          :src="user.photoUrl || defaultPhoto"
          alt="User Photo"
          class="contact-photo"
        />
        <div class="contact-info">
          <span class="contact-name">{{ user.username }}</span>
          <button class="chat-button" @click="openChatWithUser(user)">
            Chat
          </button>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import { ref, computed, onMounted } from "vue";
import { listUsers } from "@/services/api.js";
import { useRouter } from "vue-router";

export default {
  name: "UserList",
  setup() {
    const users = ref([]);
    const error = ref("");
    const searchQuery = ref("");
    const defaultPhoto =
      "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";
    const router = useRouter();
    const currentUserID = localStorage.getItem("userID");

    // Fetch users from the backend and filter out the currently logged-in user.
    async function refreshUsers() {
      error.value = "";
      try {
        const response = await listUsers();
        // Filter out the logged-in user.
        users.value = response.data.users.filter(
          (u) => u.id !== currentUserID
        );
      } catch (err) {
        error.value = "Failed to load contacts.";
        console.error(err);
      }
    }

    // Computed property for filtering based on the search query.
    const filteredUsers = computed(() => {
      if (!searchQuery.value) return users.value;
      const query = searchQuery.value.toLowerCase();
      return users.value.filter((user) =>
        user.username.toLowerCase().includes(query)
      );
    });

    function openChatWithUser(user) {
  // Use user.conversationId if it exists; otherwise, pass an empty string.
  const conversationId = user.conversationId || "";
  router.push({
    name: 'ChatView',
    params: { conversationId },
    query: { receiverId: user.id, receiverName: user.username }
  });
}



    onMounted(() => {
      refreshUsers();
    });

    return {
      users,
      error,
      defaultPhoto,
      refreshUsers,
      openChatWithUser,
      searchQuery,
      filteredUsers,
    };
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

/* Search field styling */
.search-container {
  margin-bottom: 10px;
}

.search-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

/* Refresh button */
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

/* Contact list styles */
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
