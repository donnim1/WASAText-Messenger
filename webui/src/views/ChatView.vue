<!-- src/views/ChatView.vue -->
<template>
    <div class="chat-view">
      <h1>Conversation: {{ conversationId }}</h1>
      <button @click="refreshMessages">Refresh</button>
      <ul>
        <li v-for="msg in messages" :key="msg.id" class="message-item">
          <strong>{{ msg.sender_id }}:</strong> {{ msg.content }} <small>({{ msg.sent_at }})</small>
          <!-- Buttons for forwarding, commenting, deleting can be added here -->
          <button @click="deleteMsg(msg.id)">Delete</button>
        </li>
      </ul>
      <form @submit.prevent="send">
        <input v-model="newMessage" placeholder="Type your message" required />
        <button type="submit">Send</button>
      </form>
      <div v-if="error" class="error">{{ error }}</div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { sendMessage, deleteMessage } from '@/services/api.js';
  import { useRoute } from 'vue-router';
  
  export default {
    name: 'ChatView',
    setup() {
      const route = useRoute();
      const conversationId = route.params.conversationId;
      const messages = ref([]);
      const newMessage = ref("");
      const error = ref("");
  
      // Dummy function: replace with an API call to get conversation messages if available
      async function refreshMessages() {
        error.value = "";
        // Example: load messages from an API endpoint (e.g., getConversation)
        // Here we simply simulate:
        messages.value = []; // Fill with actual messages from API
      }
  
      async function send() {
        error.value = "";
        try {
          // For private messages, receiverId must be provided. For group messages, groupId is required.
          // Here we assume a private message for simplicity.
          await sendMessage({
            receiverId: "",  // Fill in with actual receiverId or leave empty for group messages
            content: newMessage.value,
            isGroup: false,
            groupId: ""
          });
          newMessage.value = "";
          refreshMessages();
        } catch (err) {
          error.value = "Failed to send message";
          console.error(err);
        }
      }
  
      async function deleteMsg(messageId) {
        error.value = "";
        try {
          await deleteMessage(messageId);
          refreshMessages();
        } catch (err) {
          error.value = "Failed to delete message";
          console.error(err);
        }
      }
  
      onMounted(() => {
        refreshMessages();
      });
  
      return { conversationId, messages, newMessage, error, refreshMessages, send, deleteMsg };
    },
  };
  </script>
  
  <style scoped>
  .chat-view {
    padding: 20px;
  }
  .message-item {
    border-bottom: 1px solid #ccc;
    padding: 10px 0;
  }
  .error {
    color: red;
    margin-top: 10px;
  }
  </style>
  