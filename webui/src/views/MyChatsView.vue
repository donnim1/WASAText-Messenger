<template>
  <div class="chats-view">
    <div class="chats-container">
      <!-- Left Panel: Conversations List -->
      <div class="conversations-panel">
        <div class="panel-header">
          <h2>Conversations</h2>
        </div>

        <div v-if="conversationsError" class="error">{{ conversationsError }}</div>
        
        <div class="conversations-list">
          <div
            v-for="conv in conversations"
            :key="conv.id"
            class="conversation-item"
            @click="openConversation(conv)"
          >
            <div class="conversation-avatar">
              <img
                :src="conv.group_photo || defaultPhoto"
                alt="Avatar"
                class="avatar-image"
              />
            </div>
            <div class="conversation-content">
              <div class="conversation-header">
                <h3 class="conversation-name">{{ conv.name }}</h3>
                <span class="timestamp">{{ conv.created_at }}</span>
              </div>
              <p class="last-message">
                {{ conv.last_message_content || 'No messages yet.' }}
              </p>
              <span class="timestamp">
                {{ formatTimestamp(conv.last_message_sent_at) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Panel: Contacts List -->
      <div class="contacts-panel">
        <div class="panel-header">
          <h2>Contacts</h2>
          <div class="search-container">
            <input
              v-model="userSearchQuery"
              placeholder="Search contacts..."
              class="search-input"
            />
          </div>
        </div>

        <div v-if="usersError" class="error">{{ usersError }}</div>

        <div class="contacts-list">
          <div
            v-for="user in filteredUsers"
            :key="user.id"
            class="contact-item"
            @click="openChatWithUser(user)"
          >
            <div class="contact-avatar">
              <img
                :src="user.photoUrl || defaultPhoto"
                :alt="user.username"
                class="avatar-image"
              />
              <span class="online-indicator" :class="{ 'online': user.isOnline }"></span>
            </div>
            <div class="contact-info">
              <h3 class="contact-name">{{ user.username }}</h3>
              <p class="contact-status">{{ user.status || 'Available' }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { getMyConversations, listUsers, getConversationByReceiver } from "@/services/api.js";
import { useRouter } from "vue-router";

export default {
  name: "MyChatsView",
  setup() {
    const conversations = ref([]);
    const conversationsError = ref("");
    const users = ref([]);
    const usersError = ref("");
    const userSearchQuery = ref("");
    const router = useRouter();
    const currentUserID = localStorage.getItem("userID");
    const defaultPhoto = "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";
    let refreshInterval = null;

    async function loadConversations() {
      conversationsError.value = "";
      try {
        const response = await getMyConversations();
        conversations.value = response.data.conversations || [];
      } catch (err) {
        conversationsError.value = "Failed to load conversations";
        console.error(err);
      }
    }

    async function loadUsers() {
      usersError.value = "";
      try {
        const response = await listUsers();
        users.value = response.data.users
          .filter((u) => u.id !== currentUserID)
          .map((user) => ({
            ...user,
            isOnline: false // Default to offline; replace with real online status when available.
          }));
      } catch (err) {
        usersError.value = "Failed to load contacts";
        console.error(err);
      }
    }

    // Computed property to sort conversations by latest message timestamp (fallback to created_at)
    const sortedConversations = computed(() => {
      return conversations.value.slice().sort((a, b) => {
        const tsA = new Date(a.last_message_sent_at || a.created_at);
        const tsB = new Date(b.last_message_sent_at || b.created_at);
        return tsB - tsA;
      });
    });

    function openConversation(conv) {
      router.push({ name: "ChatView", params: { conversationId: conv.id } });
    }

    function openChatWithUser(user) {
      getConversationByReceiver(user.id)
        .then(response => {
          // Open chat if conversation exists
          router.push({
            name: "ChatView",
            params: { conversationId: response.data.conversationId }
          });
        })
        .catch(error => {
          if (error.response && error.response.status === 404) {
            // If conversation doesn't exist, navigate with receiver's ID to start a new conversation
            router.push({
              name: "ChatView",
              params: { conversationId: "" },
              query: { receiverId: user.id, receiverName: user.username }
            });
          } else {
            console.error("Error checking conversation:", error);
          }
        });
    }

    function formatTimestamp(ts) {
      if (!ts) return '';
      const date = new Date(ts);
      const now = new Date();
      const diff = now - date;
      const days = Math.floor(diff / (1000 * 60 * 60 * 24));
      
      if (days === 0) {
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } else if (days === 1) {
        return 'Yesterday';
      } else if (days < 7) {
        return date.toLocaleDateString([], { weekday: 'long' });
      } else {
        return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
      }
    }

    onMounted(() => {
      // Initial data load
      loadConversations();
      loadUsers();

      // Auto-refresh every 1 second
      refreshInterval = setInterval(() => {
        loadConversations();
        loadUsers();
      }, 1000);
    });

    onUnmounted(() => {
      if (refreshInterval) {
        clearInterval(refreshInterval);
      }
    });

    return {
      conversations: sortedConversations, // use sortedConversations here
      conversationsError,
      users,
      usersError,
      userSearchQuery,
      openConversation,
      openChatWithUser,
      formatTimestamp,
      defaultPhoto
    };
  },
};
</script>

<style scoped>
.chats-view {
  height: 100%;
  background-color: #f8f9fa;
  overflow: hidden;
}

.chats-container {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 1px;
  height: 100%;
  background-color: #e9ecef;
  max-width: 1400px;
  margin: 0 auto;
}

.conversations-panel,
.contacts-panel {
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

.panel-header h1,
.panel-header h2 {
  margin: 0 0 15px 0;
  color: #212529;
  font-size: 1.5rem;
  font-weight: 600;
}

.search-container {
  position: relative;
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

.conversations-list,
.contacts-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.conversation-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.conversation-item:hover {
  background-color: #f8f9fa;
}

.conversation-avatar,
.contact-avatar {
  position: relative;
  margin-right: 15px;
}

.avatar-image {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.online-indicator {
  position: absolute;
  bottom: 2px;
  right: 2px;
  width: 12px;
  height: 12px;
  background-color: #adb5bd; /* Default offline color */
  border: 2px solid #ffffff;
  border-radius: 50%;
}

.online-indicator.online {
  background-color: #40c057; /* Online color */
}

.conversation-content {
  flex: 1;
  min-width: 0;
}

.conversation-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 4px;
}

.conversation-name {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #212529;
}

.timestamp {
  font-size: 0.75rem;
  color: #868e96;
}

.last-message {
  margin: 0;
  font-size: 0.875rem;
  color: #495057;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.last-message.no-message {
  color: #adb5bd;
  font-style: italic;
}

.contact-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.contact-item:hover {
  background-color: #f8f9fa;
}

.contact-info {
  flex: 1;
}

.contact-name {
  margin: 0;
  font-size: 0.9rem;
  font-weight: 500;
  color: #212529;
}

.contact-status {
  margin: 0;
  font-size: 0.8rem;
  color: #868e96;
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

.error {
  color: #dc3545;
  padding: 10px;
  margin: 10px;
  background-color: #fff5f5;
  border-radius: 4px;
  font-size: 0.9rem;
}

/* Scrollbar Styling */
.conversations-list::-webkit-scrollbar,
.contacts-list::-webkit-scrollbar {
  width: 6px;
}

.conversations-list::-webkit-scrollbar-track,
.contacts-list::-webkit-scrollbar-track {
  background: #f1f3f5;
}

.conversations-list::-webkit-scrollbar-thumb,
.contacts-list::-webkit-scrollbar-thumb {
  background-color: #ced4da;
  border-radius: 3px;
}

.conversations-list::-webkit-scrollbar-thumb:hover,
.contacts-list::-webkit-scrollbar-thumb:hover {
  background-color: #adb5bd;
}
</style>