<!-- src/views/MyChatsView.vue -->
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
          <!-- For private chats, conv.name should be set by the backend to the partner's name -->
          <strong>{{ conv.name || 'Private Chat' }}</strong>
          <span class="timestamp">{{ formatTimestamp(conv.created_at) }}</span>
        </div>
      </li>
    </ul>
    <div v-else class="empty-message">
      You have no conversations yet.
    </div>

    <hr />

    <!-- New Chat / Search Contacts Section -->
    <h2>Start a New Chat</h2>
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

    // Load existing conversations for the logged-in user.
    async function loadConversations() {
      conversationsError.value = "";
      try {
        const response = await getMyConversations();
        conversations.value = response.data.conversations;
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

    // Computed property for filtering users by search query.
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

    // Navigate to ChatView with receiver details (for starting a new chat).
    function openChatWithUser(user) {
      router.push({
        name: "ChatView",
        params: { conversationId: "" },
        query: { receiverId: user.id, receiverName: user.username }
      });
    }

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
    };
  },
};
</script>

<style scoped>
.chats-view {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
  background-color: #f8f9fa;
}
.conversation-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.conversation-item {
  padding: 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
  transition: background-color 0.2s ease;
}
.conversation-item:hover {
  background-color: #e9ecef;
}
.conversation-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.timestamp {
  font-size: 0.8rem;
  color: #666;
}
.empty-message {
  margin-top: 20px;
  text-align: center;
  font-style: italic;
  color: #666;
}
.search-container {
  margin: 10px 0;
}
.search-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.user-search-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.user-search-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
  transition: background-color 0.2s ease;
}
.user-search-item:hover {
  background-color: #e9ecef;
}
.user-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}
.user-name {
  font-size: 1rem;
}
.error {
  color: red;
  margin-top: 10px;
}
</style>
