<template>
  <div class="chats-view">
    <h1>Your Conversations</h1>
    <div v-if="error" class="error">{{ error }}</div>
    <ul class="conversation-list">
      <li
        v-for="conv in conversations"
        :key="conv.id"
        class="conversation-item"
        @click="openConversation(conv.id)"
      >
        <div class="conversation-info">
          <strong>{{ conv.name || 'Private Chat' }}</strong>
          <span class="timestamp">{{ formatTimestamp(conv.created_at) }}</span>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { getMyConversations } from "@/services/api.js";
import { useRouter } from "vue-router";

export default {
  name: "MyChatsView",
  setup() {
    const conversations = ref([]);
    const error = ref("");
    const router = useRouter();

    async function loadConversations() {
      error.value = "";
      try {
        const response = await getMyConversations();
        conversations.value = response.data.conversations;
      } catch (err) {
        error.value = "Failed to load conversations";
        console.error(err);
      }
    }

    function openConversation(conversationId) {
      // Navigate to ChatView with the conversationId as a route parameter.
      router.push({ name: "ChatView", params: { conversationId } });
    }

    function formatTimestamp(ts) {
      return new Date(ts).toLocaleTimeString();
    }

    onMounted(() => {
      loadConversations();
    });

    return {
      conversations,
      error,
      openConversation,
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

.error {
  color: red;
  margin-top: 10px;
}
</style>
