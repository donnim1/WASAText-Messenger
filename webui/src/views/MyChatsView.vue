<!-- filepath: /c:/Users/HP/OneDrive/Desktop/2ND SEM/WASAText/webui/src/views/MyChatsView.vue -->
<template>
  <div class="chats-view">
    <h1>Your Conversations</h1>
    <!-- Existing Conversations Section -->
    <div v-if="conversationsError" class="error">{{ conversationsError }}</div>
    <ul v-if="conversations.length" class="conversation-list">
      <li
        v-for="conv in conversations"
        :key="conv.id"
        class="conversation-item"
        @click="openConversation(conv.id)"
      >
        <div class="conversation-info">
          <div class="conversation-header">
            <strong>{{ conv.name || 'Private Chat' }}</strong>
            <span class="timestamp">
              {{ formatTimestamp(conv.lastMessage && conv.lastMessage.sent_at ? conv.lastMessage.sent_at : conv.created_at) }}
            </span>
          </div>
          <div class="conversation-preview">
            <span class="last-message">
              {{ conv.lastMessage ? conv.lastMessage.content : 'No messages yet.' }}
            </span>
          </div>
        </div>
      </li>
    </ul>
    <div v-else class="empty-message">
      You have no conversations yet.
    </div>

    <hr />

    <!-- New Chat / Search Contacts Section -->
    <h2>Start a Chat</h2>
    <div class="search-container">
      <input
        v-model="userSearchQuery"
        placeholder="Search contacts..."
        class="search-input"
      />
    </div>
    <ul class="user-search-list">
      <li
        v-for="user in filteredUsers"
        :key="user.id"
        class="user-search-item"
        @click="openChatWithUser(user)"
      >
        <img
          :src="user.photoUrl || defaultPhoto"
          alt="User Photo"
          class="user-photo"
        />
        <span class="user-name">{{ user.username }}</span>
      </li>
    </ul>
    <div v-if="usersError" class="error">{{ usersError }}</div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from "vue";
import { getMyConversations, listUsers } from "@/services/api.js";
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
    const defaultPhoto = "path/to/default-photo.png"; // Update with your default photo URL

    // Load existing conversations for the logged-in user.
    async function loadConversations() {
      conversationsError.value = "";
      try {
        const response = await getMyConversations();
        // Ensure that conversations is always an array
        conversations.value = response.data.conversations || [];
      } catch (err) {
        conversationsError.value = "Failed to load conversations";
        console.error(err);
      }
    }

    // Load all users and filter out the current user.
    async function loadUsers() {
      usersError.value = "";
      try {
        const response = await listUsers();
        users.value = response.data.users.filter(
          (u) => u.id !== currentUserID
        );
      } catch (err) {
        usersError.value = "Failed to load contacts";
        console.error(err);
      }
    }

    // Computed property to filter users based on search query.
    const filteredUsers = computed(() => {
      if (!userSearchQuery.value) return users.value;
      const query = userSearchQuery.value.toLowerCase();
      return users.value.filter((user) =>
        user.username.toLowerCase().includes(query)
      );
    });

    // Navigate to the ChatView for an existing conversation.
    function openConversation(conversationId) {
      router.push({ name: "ChatView", params: { conversationId } });
    }

    // Navigate to ChatView with receiver details (to start a new chat).
    function openChatWithUser(user) {
      router.push({
        name: "ChatView",
        params: { conversationId: "" },
        query: { receiverId: user.id, receiverName: user.username }
      });
    }

    // Format timestamp for display.
    function formatTimestamp(ts) {
      return new Date(ts).toLocaleTimeString();
    }

    onMounted(() => {
      loadConversations();
      loadUsers();
    });

    return {
      conversations,
      conversationsError,
      users,
      usersError,
      userSearchQuery,
      filteredUsers,
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
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
  background-color: #f1f3f5;
  font-family: Arial, sans-serif;
}

.conversation-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.conversation-item {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  margin-bottom: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: transform 0.1s ease, box-shadow 0.1s ease;
}

.conversation-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.conversation-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.conversation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #dee2e6;
  padding-bottom: 4px;
  margin-bottom: 4px;
}

.timestamp {
  font-size: 0.8rem;
  color: #6c757d;
}

.conversation-preview {
  color: #495057;
  font-size: 0.9rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-message {
  margin-top: 20px;
  text-align: center;
  font-style: italic;
  color: #6c757d;
}

.search-container {
  margin: 20px 0;
}

.search-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 1rem;
}

.user-search-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.user-search-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #dee2e6;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.user-search-item:hover {
  background-color: #e9ecef;
}

.user-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.user-name {
  font-size: 1.1rem;
  color: #212529;
}

.error {
  color: #dc3545;
  margin-top: 10px;
}
</style>