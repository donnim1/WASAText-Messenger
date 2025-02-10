<template>
    <div class="chat-view">
      <div class="sidebar">
        <!-- Optional: You could show a sidebar with conversation list here -->
        <h2>Chat</h2>
      </div>
      <div class="chat-window">
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
          <input v-model="newMessage" placeholder="Type your message..." required />
          <button type="submit">Send</button>
        </form>
        <div v-if="chatError" class="error">{{ chatError }}</div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { sendMessage, getConversation } from '@/services/api.js';
  import { useRoute, useRouter } from 'vue-router';
  
  export default {
    name: 'ChatView',
    setup() {
      const route = useRoute();
      const router = useRouter();
  
      // If a conversation ID is provided as a route param, use it; otherwise, we rely on the receiverId query param.
      const conversationId = ref(route.params.conversationId || "");
      const receiverId = ref(route.query.receiverId || "");
      const messages = ref([]);
      const newMessage = ref("");
      const chatError = ref("");
      const currentUserId = localStorage.getItem("userID") || "";
  
      // This function loads conversation messages if conversationId exists.
      async function loadConversation() {
        if (conversationId.value) {
          try {
            const response = await getConversation(conversationId.value);
            messages.value = response.data.messages;
          } catch (err) {
            chatError.value = "Failed to load conversation";
            console.error(err);
          }
        }
      }
  
      // Function to send a message.
      async function sendMessageHandler() {
        chatError.value = "";
        if (!newMessage.value.trim()) return;
  
        try {
          // For a private conversation, if conversationId is empty, use receiverId.
          const payload = {
            receiverId: conversationId.value ? "" : receiverId.value, // if conversationId exists, assume this conversation is already open
            content: newMessage.value,
            isGroup: false,
            groupId: ""
          };
          const response = await sendMessage(payload);
          // If conversationId is not set, assume the backend created one and (optionally) return it.
          // You might modify your backend to return conversationId along with messageId.
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
          chatError.value = "Failed to send message";
          console.error(err);
        }
      }
  
      // Format timestamp
      function formatTimestamp(ts) {
        return new Date(ts).toLocaleTimeString();
      }
  
      onMounted(() => {
        loadConversation();
      });
  
      return {
        messages,
        newMessage,
        chatError,
        sendMessageHandler,
        conversationId,
        receiverId,
        currentUserId,
        formatTimestamp
      };
    }
  };
  </script>
  
  <style scoped>
  .chat-view {
    display: flex;
    height: 100vh;
    font-family: sans-serif;
  }
  
  .sidebar {
    width: 25%;
    background-color: #f5f5f5;
    border-right: 1px solid #ccc;
    padding: 10px;
    overflow-y: auto;
  }
  
  .chat-window {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 10px;
    background-color: #e9ecef;
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
  