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
          <!-- Only show reply preview when message was sent as a reply -->
          <div v-if="msg.ReplyTo && msg.ReplyTo !== ''" class="inline-reply-preview">
            <small>In reply to: {{ getReplyContent(msg.ReplyTo) }}</small>
          </div>
          
          <!-- Render message content -->
          <div v-if="isImage(msg.Content)">
            <img :src="msg.Content" alt="Image message" class="sent-image" />
          </div>
          <div v-else>
            <p class="message-content">{{ msg.Content }}</p>
          </div>
          
          <!-- Display reactions if available -->
          <div v-if="msg.reactions && msg.reactions.length" class="message-reactions">
            <span 
              v-for="(reaction, index) in msg.reactions" 
              :key="index" 
              class="reaction"
            >
              {{ reaction.Reaction }} by <small>{{ reaction.UserName }}</small>
            </span>
          </div>

          <span class="message-timestamp">{{ formatTimestamp(msg.SentAt) }}</span>
          
          <!-- Message Actions -->
          <div class="message-actions">
            <button @click="replyTo(msg)">Reply</button>
            <button @click="showForwardDialog(msg)">Forward</button>
            <!-- Heart button available for all messages -->
            <button class="heart-button" @click="toggleHeart(msg)">❤️</button>
            <!-- Conditionally show "Remove Reaction" button if user already reacted -->
            <button v-if="msg.reactions && msg.reactions.some(r => r.Reaction === '❤️' && r.UserID === currentUserId)"
                    @click="removeReaction(msg)">
              Remove Reaction
            </button>
            <!-- Only allow delete if it's your message -->
            <button v-if="msg.SenderID === currentUserId" @click="deleteMessage(msg.ID)">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Reply Preview Section -->
    <div v-if="replyingTo" class="reply-preview">
      <p>Replying to: {{ replyingTo.Content }}</p>
      <button class="cancel-reply" @click="cancelReply">Cancel Reply</button>
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
        <!-- Search input for conversations -->
        <input v-model="forwardSearchQuery" placeholder="Search chats or users" />
        <!-- List suggestions filtered by the search query -->
        <ul v-if="forwardSearchQuery && filteredConversations.length" class="conversation-list">
          <li v-for="conv in filteredConversations" :key="conv.id" @click="selectForwardTarget(conv)">
            {{ conv.name }}
          </li>
        </ul>
        <!-- Show selected target if any -->
        <p v-if="forwardTargetConversation">Selected: {{ forwardTargetConversation.name }}</p>
        <button @click="confirmForwardMessage" :disabled="!forwardTargetConversation">
          Forward
        </button>
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
  uploadImage,
  getMyConversations // Updated import—use getMyConversations instead of getConversations.
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
    const forwardTargetConversation = ref(null);
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

    const messagesMap = computed(() => {
      const map = {};
      messages.value.forEach(m => {
        map[m.ID] = m;
      });
      return map;
    });

    // Helper to get reply content from reply_to ID.
    function getReplyContent(replyID) {
      const originalMsg = messagesMap.value[replyID];
      if (!originalMsg) return "[Message deleted]";
      return isImage(originalMsg.Content) ? "Image" : originalMsg.Content;
    }

    async function toggleHeart(message) {
      try {
        const heart = "❤️";
        // Get current user's ID
        const currentUserId = localStorage.getItem("userID");
        
        console.log("Toggle heart for message:", message.ID);
        console.log("Current user ID:", currentUserId);
        
        if (message.reactions) {
          console.log("All reactions on this message:", message.reactions);
        }
        
        // Fix: Use lowercase property names to match the actual data structure
        const hasHeart = message.reactions && 
                        message.reactions.some(r => {
                          console.log("Checking reaction:", r);
                          console.log("r.userID:", r.userID, "currentUserId:", currentUserId);
                          console.log("Are they equal?", r.userID === currentUserId);
                          console.log("r.reaction:", r.reaction, "heart:", heart);
                          console.log("Are they equal?", r.reaction === heart);
                          return r.userID === currentUserId && r.reaction === heart;
                        });
        
        console.log("Has heart reaction?", hasHeart);
        
        if (hasHeart) {
          console.log("REMOVING reaction for message:", message.ID);
          const result = await uncommentMessageApi(message.ID);
          console.log("Uncomment result:", result);
        } else {
          console.log("ADDING reaction for message:", message.ID);
          const result = await commentMessageApi(message.ID, heart);
          console.log("Comment result:", result);
        }
        
        // Refresh messages from backend
        await loadConversationMessages(conversationId.value);
      } catch (err) {
        chatError.value = "Failed to toggle heart reaction";
        console.error("Toggle heart error:", err);
      }
    }

    function showForwardDialog(message) {
      messageToForward.value = message;
      showForwardModal.value = true;
      forwardSearchQuery.value = "";
      forwardTargetConversation.value = null;
      loadConversations();
    }

    function closeForwardModal() {
      showForwardModal.value = false;
      forwardSearchQuery.value = "";
      forwardTargetConversation.value = null;
      messageToForward.value = null;
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
      const container = messagesContainer.value;
      if (!container) return;

      // Define a threshold in pixels. Only auto-scroll if the user is within 50px of the bottom.
      const threshold = 10;
      const distanceFromBottom = container.scrollHeight - container.scrollTop - container.clientHeight;

      if (distanceFromBottom < threshold) {
        container.scrollTop = container.scrollHeight;
      }
    }

    async function sendMessageHandler() {
      if (!newMessage.value.trim()) return;

      const payload = {
        conversationId: conversationId.value,
        receiverId: receiverId.value,
        content: newMessage.value,
        isGroup: false,
        groupId: "",
        replyTo: replyingTo.value ? replyingTo.value.ID : ""
      };

      try {
        const response = await sendMessage(payload);
        
        // If this is a new conversation, update the conversationId
        if (!conversationId.value && response.data && response.data.conversationId) {
          conversationId.value = response.data.conversationId;
          
          // Update URL without reloading the page
          router.replace({
            name: 'ChatView',
            params: { conversationId: response.data.conversationId }
          }, { replace: true });
        }
        
        newMessage.value = "";
        replyingTo.value = null; // clear reply state after sending
        
        // Load messages with the updated conversationId
        await loadConversationMessages(conversationId.value);
      } catch (err) {
        console.error("Error sending message:", err);
        chatError.value = "Failed to send message";
      }
    }

    async function handleImageUpload(event) {
      const file = event.target.files[0];
      if (!file) return;

      const reader = new FileReader();
      reader.onload = async function(e) {
        const imageUrl = e.target.result; // Base64 encoded data URL

        const payload = {
          conversationId: conversationId.value,
          receiverId: receiverId.value,
          content: imageUrl,
          isGroup: false, // Adjust based on chat type if necessary
          groupId: ""     // Adjust if it's a group chat
        };

        try {
          await sendMessage(payload);
          await loadConversationMessages(conversationId.value);
        } catch (err) {
          chatError.value = "Failed to send image message";
          console.error("Image message error:", err);
        }
      };

      reader.readAsDataURL(file);
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
      }, 1000);
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
      if (!messageToForward.value || !forwardTargetConversation.value) {
        chatError.value = "Missing required information";
        return;
      }
      // Use conv.id if available; otherwise, try conv.ID.
      const targetId =
        forwardTargetConversation.value.id ||
        forwardTargetConversation.value.ID;
      if (!targetId) {
        chatError.value = "Selected conversation has no valid ID.";
        return;
      }
      
      try {
        await forwardMessageApi(messageToForward.value.ID, targetId);
        showForwardModal.value = false;
        forwardTargetConversation.value = null;
        messageToForward.value = null;
        chatError.value = "Message forwarded successfully";
        await loadConversationMessages(conversationId.value);
      } catch (error) {
        console.error("Forward Message Error:", error);
        chatError.value = "Failed to forward message";
      }
    };

    const isImage = (content) => {
      if (!content) return false;
      // Check for Base64 URL starting with "data:image/"
      if (content.startsWith("data:image/")) return true;

      // Otherwise, test for normal HTTP/HTTPS image URL based on common extensions.
      const imageExtensions = ["jpg", "jpeg", "png", "gif", "webp"];
      const pattern = new RegExp(`https?://.*\\.(${imageExtensions.join("|")})(\\?.*)?$`, "i");
      return pattern.test(content);
    };

    async function removeReaction(message) {
      try {
        // Call the uncomment endpoint to remove the reaction for this message.
        await uncommentMessageApi(message.ID);
        // Refresh messages after removing the reaction.
        await loadConversationMessages(conversationId.value);
      } catch (err) {
        chatError.value = "Failed to remove reaction";
        console.error("Remove Reaction error:", err);
      }
    }

    function selectForwardTarget(conv) {
      forwardTargetConversation.value = conv;
      forwardSearchQuery.value = conv.name; // Shows the selected conversation
    }

    async function loadConversations() {
      try {
        const response = await getMyConversations();
        // Assume API returns an array of conversations under response.data.conversations
        conversations.value = response.data.conversations;
      } catch (err) {
        console.error("Error loading conversations:", err);
      }
    }

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
      closeForwardModal,
      forwardMessage,
      deleteMessage,
      toggleHeart,
      confirmForwardMessage,
      handleImageUpload,
      replyingTo,
      replyTo,
      cancelReply,
      isImage, // <-- Added here
      messagesMap, // <-- Added here
      getReplyContent, // <-- Added here
      removeReaction, // <-- Added here
      selectForwardTarget // <-- Added here
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

.image-upload-button {
  display: inline-block;
  font-size: 1.8rem; /* Increase font size for the icon */
  width: 48px;       /* Increase width */
  height: 48px;      /* Increase height */
  line-height: 48px; /* Center the icon vertically */
  text-align: center;
  border: 2px solid #007bff;
  border-radius: 50%;
  color: #007bff;
  cursor: pointer;
  transition: background-color 0.2s ease, color 0.2s ease;
}

.image-upload-button:hover {
  background-color: #007bff;
  color: #fff;
}

.reply-preview {
  background-color: #f1f1f1;
  padding: 8px;
  border-left: 4px solid #007bff;
  margin-bottom: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cancel-reply {
  background: none;
  border: none;
  color: #007bff;
  cursor: pointer;
}

.chat-input-container {
  display: flex;
  gap: 8px;
  align-items: center;
}

.chat-input {
  flex: 1;
  padding: 8px;
}

.send-button {
  padding: 8px 16px;
}

.inline-reply-preview {
  background-color: #e9f5ff;
  border-left: 3px solid #007bff;
  padding: 4px 8px;
  margin-bottom: 6px;
  border-radius: 4px;
  font-size: 0.85rem;
  color: #555;
}

.message-reactions {
  margin-top: 4px;
  display: flex;
  gap: 4px;
  font-size: 1.2rem;
}
.reaction {
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 50%;
  padding: 2px 6px;
}

.conversation-list {
  list-style: none;
  padding: 0;
  max-height: 200px;
  overflow-y: auto;
  margin: 8px 0;
}

.conversation-list li {
  padding: 8px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
}

.conversation-list li:hover {
  background: #f1f2f5;
}
</style>
