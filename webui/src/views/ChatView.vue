<template>
  <div class="chat-view">
    <header class="chat-header">
      <button class="back-button" @click="goBack">←</button>
      <h2>{{ conversationTitle }}</h2>
    </header>

    <!-- Messages Container -->
    <div class="chat-messages" ref="messagesContainer">
      <div
        v-for="msg in messages"
        :key="msg.ID"
        :class="['message-wrapper', { 'sent': msg.SenderID === currentUserId }]"
      >
        <div class="message-bubble">
          <!-- Heart overlay on top of message (like Instagram) -->
          <div class="heart-overlay" v-if="msg.reactions && msg.reactions.includes('❤️')">
            ❤️
          </div>
          <p class="message-content">{{ msg.Content }}</p>
          <span class="message-timestamp">{{ formatTimestamp(msg.SentAt) }}</span>

          <!-- Message Actions -->
          <div class="message-actions">
            <button @click="showForwardDialog(msg)">Forward</button>
            <!-- Heart button available for all messages -->
            <button class="heart-button" @click="toggleHeart(msg)">❤️</button>
            <!-- Only allow delete if it's your message -->
            <button v-if="msg.SenderID === currentUserId" @click="deleteMessage(msg.ID)">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="chat-input-container">
      <form @submit.prevent="sendMessageHandler" class="chat-input-form">
        <input v-model="newMessage" placeholder="Type a message..." required />
        <button type="submit">Send</button>
      </form>
    </div>

    <!-- Forward Modal -->
    <div v-if="showForwardModal" class="modal">
      <div class="modal-content">
        <h3>Forward Message</h3>
        <input v-model="forwardTargetConversationId" placeholder="Target Conversation ID" />
        <button @click="confirmForwardMessage">Forward</button>
        <button @click="closeForwardModal">Cancel</button>
      </div>
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
import {
  getConversation,
  sendMessage,
  getConversationByReceiver,
  forwardMessage as forwardMessageApi,
  commentMessage as commentMessageApi,
  uncommentMessage as uncommentMessageApi,
  deleteMessage as deleteMessageApi,
} from "@/services/api.js";

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
    const currentUserId = localStorage.getItem("userID") || "";
    const messagesContainer = ref(null);
    const loading = ref(false);

    // Modal and Message Actions State
    const showForwardModal = ref(false);
    const messageToForward = ref(null);
    const forwardSearchQuery = ref("");
    const conversations = ref([]);
    const selectedMessageId = ref(null);
    const forwardTargetConversationId = ref("");
    const defaultPhoto = "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";

    // For filtering conversation list
    const filteredConversations = computed(() => {
      if (!forwardSearchQuery.value) return conversations.value;
      const query = forwardSearchQuery.value.toLowerCase();
      return conversations.value.filter(conv =>
        conv.name.toLowerCase().includes(query)
      );
    });

    const conversationTitle = computed(() => receiverName.value || "Chat");

    async function toggleHeart(message) {
      try {
        const heart = "❤️";
        // Check if the message already has a heart reaction.
        const hasHeart = message.reactions && message.reactions.includes(heart);
        if (hasHeart) {
          // Remove the heart reaction.
          await uncommentMessageApi(message.ID);
          message.reactions = message.reactions.filter(r => r !== heart);
        } else {
          // Add the heart reaction.
          await commentMessageApi(message.ID, heart);
          if (!message.reactions) message.reactions = [];
          message.reactions.push(heart);
        }
        await loadConversationMessages(conversationId.value);
      } catch (err) {
        chatError.value = "Failed to toggle heart reaction";
        console.error(err);
      }
    }

    function showForwardDialog(message) {
      messageToForward.value = message;
      showForwardModal.value = true;
    }

    function closeForwardModal() {
      showForwardModal.value = false;
      forwardTargetConversationId.value = "";
    }

    async function forwardMessage(message, targetConversationId) {
      try {
        await forwardMessageApi(message.ID, targetConversationId);
        showForwardModal.value = false;
        messageToForward.value = null;
        chatError.value = "Message forwarded successfully";
      } catch (err) {
        chatError.value = "Failed to forward message";
        console.error(err);
      }
    }

    async function deleteMessage(messageId) {
      if (confirm("Are you sure you want to delete this message?")) {
        try {
          await deleteMessageApi(messageId);
          chatError.value = "";
          await loadConversationMessages(conversationId.value);
        } catch (error) {
          console.error("Delete Message Error:", error);
          chatError.value = "Failed to delete message";
        }
      }
    }

    async function loadConversationMessages(convId) {
      if (!convId) return;
      loading.value = true;
      chatError.value = "";
      try {
        const response = await getConversation(convId);
        messages.value = response.data.messages;
        await nextTick();
        scrollToBottom();
      } catch (err) {
        console.error("Error loading messages:", err);
        chatError.value = "Failed to load messages";
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
            name: 'ChatView',
            params: { conversationId: response.data.conversationId }
          });
          await loadConversationMessages(response.data.conversationId);
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
        if (!conversationId.value && response.data.conversationId) {
          conversationId.value = response.data.conversationId;
          router.replace({
            name: 'ChatView',
            params: { conversationId: response.data.conversationId }
          });
        }
        await loadConversationMessages(conversationId.value);
      } catch (err) {
        console.error("Error sending message:", err);
        chatError.value = "Failed to send message";
      }
    }

    async function initializeChat() {
      if (conversationId.value) {
        await loadConversationMessages(conversationId.value);
      } else if (receiverId.value) {
        await checkExistingConversation();
      }
    }

    function goBack() {
      router.back();
    }

    onMounted(() => {
      initializeChat();
    });

    watch(
      () => route.params.conversationId,
      async (newId) => {
        if (newId) {
          conversationId.value = newId;
          await loadConversationMessages(newId);
        }
      }
    );

    watch(
      () => route.query.receiverId,
      async (newReceiverId) => {
        if (newReceiverId) {
          receiverId.value = newReceiverId;
          receiverName.value = route.query.receiverName || "";
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

    const confirmForwardMessage = async () => {
      try {
        // Note: Adjust as needed if selectedMessageId changes in your logic.
        await forwardMessageApi(selectedMessageId.value, forwardTargetConversationId.value);
        chatError.value = "";
        showForwardModal.value = false;
        forwardTargetConversationId.value = '';
        await loadConversationMessages(conversationId.value);
      } catch (error) {
        console.error("Forward Message Error:", error);
        chatError.value = "Failed to forward message";
      }
    };

    return {
      conversationTitle,
      messages,
      newMessage,
      chatError,
      sendMessageHandler,
      formatTimestamp,
      currentUserId,
      goBack,
      messagesContainer,
      loading,
      showForwardModal,
      messageToForward,
      forwardSearchQuery,
      conversations,
      filteredConversations,
      defaultPhoto,
      showForwardDialog,
      closeForwardDialog: closeForwardModal,
      forwardMessage,
      deleteMessage,
      toggleHeart,
      confirmForwardMessage,
      closeForwardModal,
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
  font-size: 1.5rem;
  color: #495057;
}

.back-button:hover {
  background-color: #f8f9fa;
}

.chat-header h2 {
  margin: 0;
  font-size: 1.2rem;
  color: #212529;
}

.chat-messages {
  flex: 1;
  padding: 60px 20px 20px;
  overflow-y: auto;
  background-color: #ffffff;
  position: relative;
}

.message-wrapper {
  margin-bottom: 10px;
}

.message-bubble {
  position: relative;
  padding: 8px 12px;
  border-radius: 16px;
  max-width: 70%;
  word-wrap: break-word;
  cursor: pointer;
  transition: all 0.2s ease;
}

.sent .message-bubble {
  background-color: #dcf8c6;
  align-self: flex-end;
}

.received .message-bubble {
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
}

.message-actions {
  display: flex;
  gap: 5px;
  margin-top: 5px;
}

.message-actions button {
  padding: 5px 10px;
  font-size: 0.8rem;
  cursor: pointer;
}

/* Heart overlay (like Instagram) on top of the message */
.heart-overlay {
  position: absolute;
  top: 5px;
  right: 5px;
  font-size: 1.4rem;
  pointer-events: none;
}

/* Modal Styles */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 480px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 16px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.2rem;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #666;
  cursor: pointer;
}

.modal-body {
  padding: 16px;
  overflow-y: auto;
}

.search-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  margin-bottom: 16px;
}

.conversations-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.conversation-item:hover {
  background-color: #f5f5f5;
}

.conversation-avatar {
  width: 40px;
  height: 40px;
}

.avatar-image {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.conversation-info {
  flex: 1;
}

.conversation-info h4 {
  margin: 0;
  font-size: 1rem;
}

.last-message {
  margin: 4px 0 0;
  font-size: 0.9rem;
  color: #666;
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

.chat-input-form input {
  flex: 1;
  padding: 12px 15px;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  font-size: 0.95rem;
}

.chat-input-form button {
  padding: 12px 24px;
  background-color: #4dabf7;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  cursor: pointer;
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

/* Long Press Directive */
.message-bubble {
  touch-action: none;
}

/* First message container spacing */
.chat-messages > div:first-child .message-container {
  margin-top: 10px;
}
</style>
