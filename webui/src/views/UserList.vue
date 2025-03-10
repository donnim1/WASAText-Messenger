<template>
  <div class="users-view">
    <div class="users-container">
      <!-- Left Panel: User List -->
      <div class="users-panel">
        <div class="panel-header">
          <h1>Contacts</h1>
          <div class="search-container">
            <input
              v-model="searchQuery"
              placeholder="Search contacts..."
              class="search-input"
            />
          </div>
          <button @click="refreshUsers" class="refresh-button">
            Refresh Contacts
          </button>
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <div class="users-list">
          <div v-if="filteredUsers.length" class="user-items">
            <div
              v-for="user in filteredUsers"
              :key="user.id"
              class="user-item"
              @click="openChatWithUser(user)"
            >
              <div class="user-avatar">
                <img
                  :src="user.photoUrl || defaultPhoto"
                  :alt="user.username"
                  class="avatar-image"
                />
              </div>
              <div class="user-content">
                <h3 class="user-name">{{ user.username }}</h3>
                <p class="user-status">Available</p>
              </div>
              <button class="chat-button">Chat</button>
            </div>
          </div>
          <div v-else class="empty-state">
            <div class="empty-icon">👥</div>
            <h3>No Contacts Found</h3>
            <p>Try searching with a different name</p>
          </div>
        </div>
      </div>

      <!-- Right Panel: Quick Info -->
      <div class="info-panel">
        <div class="panel-header">
          <h2>Quick Info</h2>
        </div>
        <div class="info-content">
          <div class="info-card">
            <h3>Start Chatting</h3>
            <p>Click on any contact to start a conversation instantly.</p>
          </div>
          <div class="info-card">
            <h3>Search Contacts</h3>
            <p>Use the search bar to quickly find specific contacts.</p>
          </div>
          <div class="info-card">
            <h3>Stay Connected</h3>
            <p>Keep in touch with your friends and colleagues.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { listUsers, getConversationByReceiver } from "@/services/api.js";
import { useRouter } from "vue-router";

export default {
  name: "UserList",
  setup() {
    const users = ref([]);
    const searchQuery = ref("");
    const error = ref("");
    const router = useRouter();
    const currentUserID = localStorage.getItem("userID");
    const defaultPhoto = "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";

    const filteredUsers = computed(() => {
      if (!searchQuery.value) return users.value;
      const query = searchQuery.value.toLowerCase();
      return users.value.filter(user =>
        user.username.toLowerCase().includes(query)
      );
    });

    async function refreshUsers() {
      error.value = "";
      try {
        const response = await listUsers();
        users.value = response.data.users.filter(u => u.id !== currentUserID);
      } catch (err) {
        error.value = "Failed to load users";
        console.error(err);
      }
    }

    function openChatWithUser(user) {
      getConversationByReceiver(user.id)
        .then(response => {
          // Conversation exists
          router.push({
            name: "ChatView",
            params: { conversationId: response.data.conversationId }
          });
        })
        .catch(e => {
          // If 404, it means no conversation exists
          if (e.response && e.response.status === 404) {
            router.push({
              name: "ChatView",
              params: { conversationId: "" },
              query: { receiverId: user.id, receiverName: user.username }
            });
          } else {
            console.error("Error checking conversation:", e);
          }
        });
    }

    onMounted(() => {
      refreshUsers();
      
      // Add auto-refresh every half second
      const refreshInterval = setInterval(() => {
        refreshUsers();
      }, 500);
      
      // Clean up interval when component unmounts
      onUnmounted(() => {
        if (refreshInterval) {
          clearInterval(refreshInterval);
        }
      });
    });

    return {
      searchQuery,
      filteredUsers,
      error,
      refreshUsers,
      openChatWithUser,
      defaultPhoto
    };
  }
};
</script>

<style scoped>
.users-view {
  height: 100vh;
  background-color: #f8f9fa;
  overflow: hidden;
}

.users-container {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 1px;
  height: 100%;
  background-color: #e9ecef;
  max-width: 1400px;
  margin: 0 auto;
}

.users-panel, .info-panel {
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.panel-header {
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
  background-color: #ffffff;
}

.panel-header h1, .panel-header h2 {
  margin: 0 0 15px 0;
  color: #212529;
  font-size: 1.5rem;
  font-weight: 600;
}

.search-container {
  margin-bottom: 15px;
}

.search-input {
  width: 100%;
  padding: 10px 15px;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  font-size: 0.9rem;
  background-color: #f8f9fa;
  transition: all 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #4dabf7;
  background-color: #ffffff;
  box-shadow: 0 0 0 3px rgba(77, 171, 247, 0.1);
}

.refresh-button {
  width: 100%;
  padding: 10px;
  background-color: #4dabf7;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.2s ease;
}

.refresh-button:hover {
  background-color: #3c99e6;
}

.users-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.user-item:hover {
  background-color: #f8f9fa;
}

.user-avatar {
  margin-right: 15px;
}

.avatar-image {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.user-content {
  flex: 1;
  min-width: 0;
}

.user-name {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #212529;
}

.user-status {
  margin: 4px 0 0 0;
  font-size: 0.85rem;
  color: #868e96;
}

.chat-button {
  padding: 8px 16px;
  background-color: #4dabf7;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.chat-button:hover {
  background-color: #3c99e6;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 20px;
  text-align: center;
  color: #868e96;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

.empty-state h3 {
  margin: 0 0 10px 0;
  color: #495057;
}

.empty-state p {
  margin: 0;
  font-size: 0.9rem;
}

.error-message {
  color: #dc3545;
  padding: 10px;
  margin: 10px;
  background-color: #f8d7da;
  border-radius: 6px;
  font-size: 0.9rem;
}

.info-content {
  padding: 20px;
}

.info-card {
  background-color: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 15px;
}

.info-card h3 {
  margin: 0 0 10px 0;
  color: #212529;
  font-size: 1.1rem;
}

.info-card p {
  margin: 0;
  color: #495057;
  font-size: 0.9rem;
  line-height: 1.5;
}

/* Scrollbar Styling */
.users-list::-webkit-scrollbar {
  width: 6px;
}

.users-list::-webkit-scrollbar-track {
  background: #f1f3f5;
}

.users-list::-webkit-scrollbar-thumb {
  background-color: #ced4da;
  border-radius: 3px;
}

.users-list::-webkit-scrollbar-thumb:hover {
  background-color: #adb5bd;
}
</style>