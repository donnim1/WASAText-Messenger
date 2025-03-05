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
        :class="['message-wrapper', { sent: msg.SenderID === currentUserId }]"
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
            <button @click="replyTo(msg)">Reply</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="chat-input-container">
      <form @submit.prevent="sendMessageHandler" class="chat-input-form">
        <label for="image-upload" class="image-upload-button">
          <i class="fas fa-image"></i>
        </label>
        <input 
          id="image-upload" 
          type="file" 
          accept="image/*" 
          style="display: none" 
          @change="handleImageUpload"
        />
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

    <!-- Add reply UI -->
    <div v-if="replyingTo" class="reply-container">
      <div class="reply-preview">
        <p>Replying to: {{ replyingTo.Content }}</p>
        <button @click="cancelReply">✕</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch, computed, nextTick, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  getConversation,
  sendMessage,
  getConversationByReceiver,
  forwardMessage as forwardMessageApi,
  commentMessage as commentMessageApi,
  uncommentMessage as uncommentMessageApi,
  deleteMessage as deleteMessageApi,
  uploadImage
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

    const replyingTo = ref(null);

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
      
      // Add reply reference if replying
      if (replyingTo.value) {
        payload.replyTo = replyingTo.value.ID;
      }
      
      try {
        const response = await sendMessage(payload);
        newMessage.value = "";
        replyingTo.value = null; // Clear reply state
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

    async function handleImageUpload(event) {
      const file = event.target.files[0];
      if (!file) return;
      
      const formData = new FormData();
      formData.append('image', file);
      
      try {
        const response = await uploadImage(formData);
        const imageUrl = response.data.url;
        
        // Send as special message with image URL
        const payload = {
          conversationId: conversationId.value,
          receiverId: receiverId.value,
          content: `[Image](${imageUrl})`,
          isImage: true
        };
        
        await sendMessage(payload);
        await loadConversationMessages(conversationId.value);
      } catch (err) {
        chatError.value = "Failed to upload image";
        console.error(err);
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

    function replyTo(message) {
      replyingTo.value = message;
    }

    function cancelReply() {
      replyingTo.value = null;
    }

    let messagePollingInterval;

    function startMessagePolling() {
      // Poll every 5 seconds
      messagePollingInterval = setInterval(async () => {
        if (conversationId.value) {
          await loadConversationMessages(conversationId.value);
        }
      }, 5000);
    }

    function stopMessagePolling() {
      clearInterval(messagePollingInterval);
    }

    onMounted(() => {
      initializeChat();
      startMessagePolling();
    });

    onUnmounted(() => {
      stopMessagePolling();
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
      if (!messageToForward.value || !forwardTargetConversationId.value) {
        chatError.value = "Missing required information";
        return;
      }
      
      try {
        await forwardMessageApi(
          messageToForward.value.ID, 
          forwardTargetConversationId.value
        );
        showForwardModal.value = false;
        forwardTargetConversationId.value = '';
        messageToForward.value = null;
        chatError.value = "Message forwarded successfully";
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
      handleImageUpload,
      replyingTo,
      replyTo,
      cancelReply
    };
  },
};
</script>

<style scoped>
/* Main chat view container */
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f1f2f5;
}

/* Header styling */
.chat-header {
  background-color: #fff;
  border-bottom: 1px solid #e0e0e0;
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
  background-color: #f1f2f5;
}

.chat-header h2 {
  margin: 0;
  font-size: 1.3rem;
  color: #212529;
}

/* Messages container styling */
.chat-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f1f2f5;
}

/* Message wrapper with alignment based on sender */
.message-wrapper {
  display: flex;
  margin-bottom: 10px;
}

.message-wrapper.sent {
  justify-content: flex-end;
}

.message-wrapper:not(.sent) {
  justify-content: flex-start;
}

/* Message bubble styling for a modern look */
.message-bubble {
  position: relative;
  padding: 12px 16px;
  border-radius: 20px;
  max-width: 70%;
  word-wrap: break-word;
  background-color: #fff;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

/* Different background for sent messages */
.message-wrapper.sent .message-bubble {
  background-color: #dcf8c6;
}

/* Message content and timestamp */
.message-content {
  margin: 0;
  font-size: 1rem;
  line-height: 1.4;
  color: #333;
}

.message-timestamp {
  display: block;
  font-size: 0.75rem;
  color: #888;
  margin-top: 6px;
  text-align: right;
}

/* Message actions styling */
.message-actions {
  display: flex;
  gap: 10px;
  margin-top: 8px;
}

.message-actions button {
  padding: 4px 8px;
  font-size: 0.8rem;
  cursor: pointer;
  background: none;
  border: none;
  color: #007bff;
}

.message-actions button:hover {
  text-decoration: underline;
}

/* Heart overlay on top of the message bubble */
.heart-overlay {
  position: absolute;
  top: -5px;
  right: -5px;
  background: #fff;
  border-radius: 50%;
  padding: 3px;
  font-size: 1.2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  pointer-events: none;
}

/* Input area styling */
.chat-input-container {
  padding: 15px 20px;
  background-color: #fff;
  border-top: 1px solid #e0e0e0;
}

.chat-input-form {
  display: flex;
  gap: 10px;
}

.chat-input-form input {
  flex: 1;
  padding: 12px 15px;
  border: 1px solid #ccc;
  border-radius: 20px;
  font-size: 1rem;
}

.chat-input-form button {
  padding: 12px 24px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 20px;
  font-size: 1rem;
  cursor: pointer;
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
  background: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 480px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  padding: 20px;
}

/* Chat error message */
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
</style>
