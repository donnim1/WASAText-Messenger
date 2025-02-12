<template>
  <div class="chat-view">
    <div class="chat-header">
      <button class="back-button" @click="goBack">←</button>
      <!-- Show the partner's name if provided; otherwise, default to "Conversation" -->
      <h2>{{ conversationTitle }}</h2>
    </div>
    <div class="chat-messages">
      <div
        v-for="msg in messages"
        :key="msg.ID"
        :class="['chat-message', { sent: msg.SenderID === currentUserId, received: msg.SenderID !== currentUserId }]"
      >
        <p class="message-content">{{ msg.Content }}</p>
        <span class="message-timestamp">{{ formatTimestamp(msg.SentAt) }}</span>
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

    // Read conversationId from route parameters (if present)
    const conversationId = ref(route.params.conversationId || "");
    // Get the receiver info from query parameters; this is used to send a message.
    const receiverId = ref(route.query.receiverId || "");
    const receiverName = ref(route.query.receiverName || "");
    const messages = ref([]);
    const newMessage = ref("");
    const chatError = ref("");
    const currentUserId = localStorage.getItem("userID") || "";

    // Compute conversation title. If conversationId is not present, use the receiverName.
    const conversationTitle = computed(() => {
      if (!conversationId.value && receiverName.value) {
        return receiverName.value;
      }
      return "Conversation";
    });

    // Load messages if a conversation already exists.
    async function loadConversationMessages() {
      if (!conversationId.value) return;
      try {
        const response = await getConversation(conversationId.value);
        messages.value = response.data.messages;
      } catch (err) {
        console.error("Error loading messages:", err);
        chatError.value = "Failed to load messages";
      }
    }

    // Send a message. Always include the receiverId (from the query) in the payload.
    async function sendMessageHandler() {
      chatError.value = "";
      if (!newMessage.value.trim()) return;
      try {
        const payload = {
          receiverId: receiverId.value, // Always provide the receiver's id
          content: newMessage.value,
          isGroup: false,
          groupId: ""
        };
        const response = await sendMessage(payload);
        // If a new conversation was created (backend returns conversationId), update conversationId and the route.
        if (!conversationId.value && response.data.conversationId) {
          conversationId.value = response.data.conversationId;
          router.replace({ name: "ChatView", params: { conversationId: conversationId.value } });
        }
        messages.value.push({
          ID: response.data.messageId,
          SenderID: currentUserId,
          Content: newMessage.value,
          SentAt: new Date().toISOString()
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
      // If conversationId is missing (i.e. a new chat), there are no messages to load.
      if (conversationId.value) {
        loadConversationMessages();
      }
    });

    // Reload messages whenever conversationId changes.
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
      conversationId,
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
