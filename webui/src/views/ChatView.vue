<template>
  <div class="chat-view">
    <div class="chat-header">
      <button class="back-button" @click="goBack">←</button>
      <h2>{{ conversationTitle }}</h2>
    </div>
    <div class="chat-messages">
      <div
        v-for="msg in messages"
        :key="msg.id"
        :class="['chat-message', { sent: msg.sender_id === currentUserId, received: msg.sender_id !== currentUserId }]"
      >
        <p class="message-content">{{ msg.content }}</p>
        <span class="message-timestamp">{{ formatTimestamp(msg.sent_at) }}</span>
      </div>
    </div>
    <div class="chat-input-container">
      <form @submit.prevent="sendMessageHandler" class="chat-input-form">
        <input v-model="newMessage" placeholder="Type a message..." required />
        <button type="submit">Send</button>
      </form>
    </div>
    <div v-if="chatError" class="chat-error">{{ chatError }}</div>
  </div>
</template>

<script>
import { ref, onMounted, watch, computed } from "vue";
import { getConversation, sendMessage } from "@/services/api.js";
import { useRoute, useRouter } from "vue-router";

export default {
  name: "ChatView",
  setup() {
    const route = useRoute();
    const router = useRouter();

    // conversationId is read from route params (it may be empty for a new conversation)
    const conversationId = ref(route.params.conversationId || "");
    // We also store the full conversation details when loaded
    const conversation = ref(null);
    const messages = ref([]);
    const newMessage = ref("");
    const chatError = ref("");
    const currentUserId = localStorage.getItem("userID") || "";

    // Compute the conversation title.
    // If a conversation exists and its name is nonempty, use that.
    // Otherwise, if a receiverName query parameter is provided, use that.
    // Otherwise, fall back to "Private Chat".
    const conversationTitle = computed(() => {
      if (conversation.value) {
        if (conversation.value.name && conversation.value.name.trim() !== "") {
          return conversation.value.name;
        } else if (route.query.receiverName) {
          return route.query.receiverName;
        } else {
          return "Private Chat";
        }
      }
      return route.query.receiverName || "New Conversation";
    });

    // Load conversation details and messages
    async function loadConversationMessages() {
      if (!conversationId.value) return;
      try {
        const response = await getConversation(conversationId.value);
        conversation.value = response.data.conversation;
        messages.value = response.data.messages;
      } catch (err) {
        console.error("Error loading conversation:", err);
        chatError.value = "Failed to load conversation";
      }
    }

    // Send a message; if the conversation is new, the backend should return the conversationId.
    async function sendMessageHandler() {
      chatError.value = "";
      if (!newMessage.value.trim()) return;

      // Build the payload.
      // For private chats: if conversationId is empty, we assume a new conversation will be created.
      const payload = {
        receiverId: conversationId.value ? "" : (route.query.receiverId || ""),
        content: newMessage.value,
        isGroup: false, // adjust this if you support group messaging
        groupId: ""     // for private chats, leave empty
      };

      try {
        const response = await sendMessage(payload);
        // If the conversation is new, update conversationId if returned by backend
        if (!conversationId.value && response.data.conversationId) {
          conversationId.value = response.data.conversationId;
        }
        messages.value.push({
          id: response.data.messageId,
          sender_id: currentUserId,
          content: newMessage.value,
          sent_at: new Date().toISOString()
        });
        newMessage.value = "";
      } catch (err) {
        console.error("Send message error:", err);
        chatError.value = "Failed to send message";
      }
    }

    function formatTimestamp(ts) {
      return new Date(ts).toLocaleTimeString();
    }

    function goBack() {
      router.back();
    }

    onMounted(() => {
      if (conversationId.value) {
        loadConversationMessages();
      }
    });

    watch(conversationId, (newVal) => {
      if (newVal) {
        loadConversationMessages();
      }
    });

    return {
      conversationTitle,
      messages,
      newMessage,
      chatError,
      sendMessageHandler,
      formatTimestamp,
      currentUserId,
      goBack,
    };
  },
};
</script>

<style scoped>
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #ece5dd;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #075e54;
  color: white;
}

.back-button {
  background: transparent;
  border: none;
  color: white;
  font-size: 1.5rem;
  margin-right: 10px;
  cursor: pointer;
}

.chat-messages {
  flex: 1;
  padding: 15px;
  overflow-y: auto;
  background-color: #ffffff;
}

.chat-message {
  margin-bottom: 10px;
  padding: 8px 12px;
  border-radius: 10px;
  max-width: 70%;
  word-wrap: break-word;
}

.chat-message.sent {
  align-self: flex-end;
  background-color: #dcf8c6;
}

.chat-message.received {
  align-self: flex-start;
  background-color: #ffffff;
  border: 1px solid #ccc;
}

.message-content {
  margin: 0;
}

.message-timestamp {
  display: block;
  font-size: 0.75rem;
  color: #999;
  margin-top: 4px;
  text-align: right;
}

.chat-input-container {
  padding: 10px;
  background-color: #f0f0f0;
  border-top: 1px solid #ccc;
}

.chat-input-form {
  display: flex;
  gap: 10px;
}

.chat-input-form input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.chat-input-form button {
  margin-left: 10px;
  padding: 10px 20px;
  background-color: #128c7e;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.chat-error {
  color: red;
  text-align: center;
  margin-top: 10px;
}
</style>
