<template>
  <div class="chat-view">
    <header class="chat-header">
      <button class="back-button" @click="$router.back()">
        <i class="back-icon fas fa-arrow-left"></i>
      </button>
      <div class="chat-info">
        <img
          :src="conversation.PhotoUrl || defaultPhoto"
          alt="Avatar"
          class="chat-avatar"
        />
        <h2>{{ conversation.Name }}</h2>
      </div>
    </header>

    <!-- Messages Container -->
    <div class="chat-messages" ref="messagesContainer">
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Loading messages...</p>
      </div>

      <div v-else-if="messages.length === 0" class="empty-chat">
        <div class="empty-icon">💭</div>
        <h3>No Messages Yet</h3>
        <p>Start the conversation by sending a message</p>
      </div>
      
      <div
        v-else
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
import { ref, onMounted, watch, computed, nextTick } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getConversation, sendMessage, getConversationByReceiver } from "@/services/api.js";

export default {
  name: "ChatView",
  setup() {
    const route = useRoute();
    const router = useRouter();
    const conversationId = ref(route.params.conversationId || "");
    const receiverId = ref(route.query.receiverId || "");
    const receiverName = ref(route.query.receiverName || "");
    const messages = ref([]);
    const newMessage = ref("");
    const chatError = ref("");
    const loading = ref(false);
    const currentUserId = localStorage.getItem("userID") || "";
    const messagesContainer = ref(null);
    // Holds conversation details (name, photo, etc.)
    const conversation = ref({});
    const defaultPhoto =
      "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";

    const conversationTitle = computed(() => {
      return receiverName.value || "Chat";
    });

    async function loadConversationMessages() {
      loading.value = true;
      chatError.value = "";
      try {
        const response = await getConversation(conversationId.value);
        // Assuming your API returns the conversation object with details
        // For example:
        // { conversation: { id, name, photo_url, ... }, messages: [ ... ] }
        messages.value = response.data.messages || [];
        if(response.data.conversation) {
          conversation.value = response.data.conversation;
        }
      } catch (err) {
        chatError.value = "Failed to load messages.";
      } finally {
        loading.value = false;
      }
    }

    async function checkExistingConversation() {
      if (!receiverId.value) return;
      loading.value = true;
      try {
        const response = await getConversationByReceiver(receiverId.value);
        if (response.data.conversationId) {
          conversationId.value = response.data.conversationId;
          router.replace({
            name: "ChatView",
            params: { conversationId: conversationId.value }
          });
        }
      } catch (err) {
        console.error("No existing conversation found:", err);
      } finally {
        loading.value = false;
      }
    }

    function scrollToBottom() {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
      }
    }

    async function sendMessageHandler() {
      if (!newMessage.value.trim()) return;
      chatError.value = "";
      
      const payload = {
        conversationId: conversationId.value,
        receiverId: receiverId.value,
        content: newMessage.value,
      };

      try {
        const response = await sendMessage(payload);
        newMessage.value = "";
        
        // If no conversation existed before and the backend created one, update the conversationId and the route
        if (!conversationId.value && response.data.conversationId) {
          conversationId.value = response.data.conversationId;
          router.replace({
            name: "ChatView",
            params: { conversationId: conversationId.value }
          });
        }
        
        await loadConversationMessages();
        await nextTick();
        scrollToBottom();
      } catch (err) {
        console.error("Error sending message:", err);
        chatError.value = "Failed to send message";
      }
    }

    async function initializeChat() {
      if (conversationId.value) {
        await loadConversationMessages();
      } else if (receiverId.value) {
        await checkExistingConversation();
      }
    }

    onMounted(() => {
      initializeChat();
    });

    // Watch for route changes (if applicable)
    watch(
      () => route.params.conversationId,
      async (newId) => {
        if (newId) {
          conversationId.value = newId;
          await loadConversationMessages();
        }
      }
    );

    watch(
      () => route.query.receiverId,
      async (newReceiverId) => {
        if (newReceiverId) {
          receiverId.value = newReceiverId;
          await checkExistingConversation();
        }
      }
    );

    function formatTimestamp(ts) {
      if (!ts) return '';
      const date = new Date(ts);
      return date.toLocaleTimeString([], { 
        hour: '2-digit', 
        minute: '2-digit'
      });
    }

    function goBack() {
      router.back();
    }

    return {
      conversationTitle,
      messages,
      newMessage,
      chatError,
      loading,
      sendMessageHandler,
      formatTimestamp,
      currentUserId,
      goBack,
      messagesContainer,
      conversation,
      defaultPhoto
    };
  },
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
  display: flex;
  align-items: center;
  gap: 10px;
}

.chat-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
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

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #868e96;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #4dabf7;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
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