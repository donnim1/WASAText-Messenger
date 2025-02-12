<template>
  <div class="chat-view">
    <div class="chat-header">
      <button class="back-button" @click="goBack">←</button>
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

    // conversationId is read from route params if available.
    const conversationId = ref(route.params.conversationId || "");
    // Also, if coming from the user directory, receiverId and receiverName can be in the query.
    const receiverId = ref(route.query.receiverId || "");
    const messages = ref([]);
    const newMessage = ref("");
    const chatError = ref("");
    const currentUserId = localStorage.getItem("userID") || "";

    // Compute the conversation title:
    // If a receiverName is provided (for a new conversation) use that, otherwise default to "Conversation".
    const conversationTitle = computed(() => {
      if (route.query.receiverName && !conversationId.value) {
        return route.query.receiverName;
      }
      return "Conversation";
    });

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

    async function sendMessageHandler() {
      chatError.value = "";
      if (!newMessage.value.trim()) return;
      try {
        const payload = {
          receiverId: conversationId.value ? "" : receiverId.value, // If conversation exists, backend already knows the receiver.
          content: newMessage.value,
          isGroup: false,
          groupId: ""
        };
        const response = await sendMessage(payload);
        // If a new conversation was created, update conversationId and update the route.
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
      conversationId
    };
  }
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
