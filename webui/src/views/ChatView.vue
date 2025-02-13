<template>
  <div class="chat-view">
    <!-- Chat Header -->
    <div class="chat-header">
      <div class="header-content">
        <button class="back-button" @click="goBack">
          <span class="back-icon">←</span>
        </button>
        <div class="chat-info">
          <h2>{{ conversationTitle }}</h2>
          <span class="status">Online</span>
        </div>
      </div>
    </div>

    <!-- Messages Container -->
    <div class="chat-messages" ref="messagesContainer">
      <div v-if="messages.length === 0" class="empty-chat">
        <div class="empty-icon">💭</div>
        <h3>No Messages Yet</h3>
        <p>Start the conversation by sending a message</p>
      </div>
      
      <div
        v-for="msg in messages"
        :key="msg.ID"
        :class="['message-wrapper', { 'sent': msg.SenderID === currentUserId }]"
      >
        <div class="message-bubble">
          <p class="message-content">{{ msg.Content }}</p>
          <span class="message-timestamp">{{ formatTimestamp(msg.SentAt) }}</span>
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="chat-input-container">
      <form @submit.prevent="sendMessageHandler" class="chat-input-form">
        <input
          v-model="newMessage"
          placeholder="Type a message..."
          required
          class="message-input"
        />
        <button type="submit" class="send-button">
          Send
        </button>
      </form>
    </div>

    <!-- Error Message -->
    <div v-if="chatError" class="chat-error">
      {{ chatError }}
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getConversation, sendMessage, getConversationByReceiver } from "@/services/api.js";

export default {
  name: "ChatView",
  setup() {
    const route = useRoute();
    const router = useRouter();

    // Get conversation and user details
    const conversationId = ref(route.params.conversationId || "");
    const receiverId = ref(route.query.receiverId || "");
    const receiverName = ref(route.query.receiverName || "");
    const messages = ref([]);
    const newMessage = ref("");
    const chatError = ref("");
    const currentUserId = localStorage.getItem("userID") || "";

    // Default title when no conversation is available
    const conversationTitle = ref(receiverName.value || "Chat");

    // Function to load messages
    async function loadConversationMessages() {
      try {
        if (conversationId.value) {
          const response = await getConversation(conversationId.value);
          messages.value = response.data.messages || [];
        } else if (receiverId.value) {
          // Try to load conversation by receiver if conversationId is missing
          const response = await getConversationByReceiver(receiverId.value);
          if (response.data.conversationId) {
            conversationId.value = response.data.conversationId;
            messages.value = response.data.messages || [];
            // Update URL to include conversation ID
            router.replace({ name: "ChatView", params: { conversationId: conversationId.value } });
          }
        }
      } catch (err) {
        console.error("Error loading messages:", err);
        chatError.value = "Failed to load messages";
      }
    }

    async function sendMessageHandler() {
      if (!newMessage.value.trim()) return;
      try {
        const response = await sendMessage({
          conversationId: conversationId.value,
          receiverId: receiverId.value,
          content: newMessage.value
        });
        if (response.data.conversationId && !conversationId.value) {
          conversationId.value = response.data.conversationId;
          router.replace({ name: "ChatView", params: { conversationId: conversationId.value } });
        }
        newMessage.value = "";
        loadConversationMessages();
      } catch (err) {
        console.error("Error sending message:", err);
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
      loadConversationMessages();
    });

    watch(
      () => route.params.conversationId,
      (newId) => {
        if (newId) {
          conversationId.value = newId;
          loadConversationMessages();
        }
      }
    );

    return {
      conversationTitle,
      messages,
      newMessage,
      chatError,
      sendMessageHandler,
      formatTimestamp,
      goBack,
      currentUserId
    };
  }
};
</script>

<style scoped>
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f8f9fa;
}

.chat-header {
  background-color: #ffffff;
  border-bottom: 1px solid #e9ecef;
  padding: 15px 20px;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.back-button {
  background: transparent;
  border: none;
  padding: 8px;
  cursor: pointer;
  border-radius: 50%;
  transition: background-color 0.2s ease;
}

.back-button:hover {
  background-color: #f8f9fa;
}

.back-icon {
  font-size: 1.5rem;
  color: #495057;
}

.chat-info {
  flex: 1;
}

.chat-info h2 {
  margin: 0;
  font-size: 1.2rem;
  color: #212529;
}

.status {
  font-size: 0.85rem;
  color: #40c057;
}

.chat-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #ffffff;
}

.empty-chat {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #868e96;
  text-align: center;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

.empty-chat h3 {
  margin: 0 0 10px 0;
  color: #495057;
}

.empty-chat p {
  margin: 0;
  font-size: 0.9rem;
}

.message-wrapper {
  display: flex;
  margin-bottom: 10px;
}

.message-wrapper.sent {
  justify-content: flex-end;
}

.message-bubble {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 12px;
  background-color: #f8f9fa;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.message-wrapper.sent .message-bubble {
  background-color: #4dabf7;
  color: white;
}

.message-content {
  margin: 0;
  font-size: 0.95rem;
  line-height: 1.4;
}

.message-timestamp {
  display: block;
  font-size: 0.75rem;
  margin-top: 4px;
  opacity: 0.8;
}

.chat-input-container {
  padding: 15px 20px;
  background-color: #ffffff;
  border-top: 1px solid #e9ecef;
}

.chat-input-form {
  display: flex;
  gap: 10px;
}

.message-input {
  flex: 1;
  padding: 12px 15px;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  font-size: 0.95rem;
  transition: all 0.2s ease;
}

.message-input:focus {
  outline: none;
  border-color: #4dabf7;
  box-shadow: 0 0 0 3px rgba(77, 171, 247, 0.1);
}

.send-button {
  padding: 12px 24px;
  background-color: #4dabf7;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.send-button:hover {
  background-color: #3c99e6;
}

.chat-error {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background-color: #f8d7da;
  color: #721c24;
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 0.9rem;
}

/* Scrollbar Styling */
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-track {
  background: #f1f3f5;
}

.chat-messages::-webkit-scrollbar-thumb {
  background-color: #ced4da;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background-color: #adb5bd;
}
</style>