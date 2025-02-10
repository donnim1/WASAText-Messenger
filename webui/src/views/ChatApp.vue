<template>
  <div class="chat-app">
    <!-- Sidebar: List of Conversations -->
    <div class="sidebar">
      <h2>Chats</h2>
      <ul>
        <li
          v-for="conv in conversations"
          :key="conv.id"
          :class="{ active: conv.id === selectedConversationId }"
          @click="selectConversation(conv.id)"
        >
          <div class="chat-summary">
            <div class="chat-name">{{ conv.name || 'Private Chat' }}</div>
            <div class="chat-timestamp">{{ formatTimestamp(conv.created_at) }}</div>
          </div>
        </li>
      </ul>
    </div>

    <!-- Chat Window -->
    <div class="chat-window" v-if="selectedConversationId">
      <div class="messages">
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="['message', { sent: msg.sender_id === currentUserId, received: msg.sender_id !== currentUserId }]"
        >
          <div class="message-content">{{ msg.content }}</div>
          <div class="message-timestamp">{{ formatTimestamp(msg.sent_at) }}</div>
        </div>
      </div>
      <form class="message-form" @submit.prevent="sendMessageHandler">
        <input v-model="newMessage" placeholder="Type a message" required />
        <button type="submit">Send</button>
      </form>
      <div v-if="chatError" class="error">{{ chatError }}</div>
    </div>

    <div class="chat-window empty" v-else>
      <p>Select a conversation to start chatting.</p>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue';
import { getMyConversations, getConversation, sendMessage } from '@/services/api.js';

export default {
  name: 'ChatApp',
  setup() {
    const conversations = ref([]);
    const selectedConversationId = ref(null);
    const messages = ref([]);
    const newMessage = ref("");
    const chatError = ref("");
    const currentUserId = localStorage.getItem("userID") || "";

    async function loadConversations() {
      try {
        const response = await getMyConversations();
        conversations.value = response.data.conversations;
      } catch (err) {
        console.error("Error loading conversations:", err);
        chatError.value = "Failed to load conversations";
      }
    }

    async function loadConversationMessages(convId) {
      try {
        const response = await getConversation(convId);
        messages.value = response.data.messages;
      } catch (err) {
        console.error("Error loading messages:", err);
        chatError.value = "Failed to load messages";
      }
    }

    function selectConversation(id) {
      selectedConversationId.value = id;
      loadConversationMessages(id);
    }

    async function sendMessageHandler() {
      chatError.value = "";
      if (!newMessage.value.trim()) return;

      // Determine conversation type.
      const conv = conversations.value.find(c => c.id === selectedConversationId.value);
      let receiverId = "";
      let isGroup = false;
      let groupId = "";
      if (conv && conv.is_group) {
        isGroup = true;
        groupId = conv.id;
      } else {
        // For private chats, assume receiverId is determined by your backend logic.
        // For demonstration, we leave receiverId empty.
      }

      try {
        const response = await sendMessage({
          receiverId,
          content: newMessage.value,
          isGroup,
          groupId,
        });
        messages.value.push({
          id: response.data.messageId,
          sender_id: currentUserId,
          content: newMessage.value,
          sent_at: new Date().toISOString(),
        });
        newMessage.value = "";
      } catch (err) {
        console.error(err);
        chatError.value = "Failed to send message";
      }
    }

    function formatTimestamp(ts) {
      return new Date(ts).toLocaleTimeString();
    }

    onMounted(() => {
      loadConversations();
    });

    watch(selectedConversationId, (newVal) => {
      if (newVal) loadConversationMessages(newVal);
    });

    return {
      conversations,
      selectedConversationId,
      messages,
      newMessage,
      chatError,
      selectConversation,
      sendMessageHandler,
      currentUserId,
      formatTimestamp,
    };
  },
};
</script>

<style scoped>
.chat-app {
  display: flex;
  height: 100vh;
  font-family: sans-serif;
}

/* Sidebar styling */
.sidebar {
  width: 30%;
  background-color: #f5f5f5;
  border-right: 1px solid #ccc;
  overflow-y: auto;
  padding: 10px;
}

.sidebar h2 {
  margin: 0 0 10px;
  padding: 0 10px;
}

.sidebar ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.sidebar li {
  padding: 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
}

.sidebar li.active {
  background-color: #ddd;
}

.chat-summary {
  display: flex;
  justify-content: space-between;
}

/* Chat window styling */
.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 10px;
  background-color: #e9ecef;
}

.chat-window.empty {
  display: flex;
  align-items: center;
  justify-content: center;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  background: #fff;
  border: 1px solid #ccc;
  margin-bottom: 10px;
}

.message {
  margin-bottom: 10px;
  padding: 5px 10px;
  border-radius: 8px;
  max-width: 70%;
}

.message.sent {
  align-self: flex-end;
  background-color: #dcf8c6;
}

.message.received {
  align-self: flex-start;
  background-color: #fff;
  border: 1px solid #ccc;
}

.message-timestamp {
  font-size: 0.75rem;
  color: #666;
  margin-top: 2px;
  text-align: right;
}

.message-form {
  display: flex;
  gap: 10px;
}

.message-form input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.message-form button {
  padding: 10px 20px;
  background-color: #27ae60;
  border: none;
  color: white;
  border-radius: 4px;
  cursor: pointer;
}

.error {
  color: red;
  margin-top: 10px;
}
</style>
