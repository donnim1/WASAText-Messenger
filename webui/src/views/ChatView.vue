<template>
    <div class="chat-view">
      <div class="chat-header">
        <h2>Chat with {{ receiverId || "Unknown" }}</h2>
      </div>
      <div class="chat-messages">
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="['message', { sent: msg.sender_id === currentUserId, received: msg.sender_id !== currentUserId }]"
        >
          <div class="message-content">{{ msg.content }}</div>
          <div class="message-timestamp">{{ formatTimestamp(msg.sent_at) }}</div>
        </div>
      </div>
      <form class="chat-input" @submit.prevent="handleSendMessage">
        <input v-model="newMessage" placeholder="Type your message..." required />
        <button type="submit">Send</button>
      </form>
      <div v-if="chatError" class="error">{{ chatError }}</div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { sendMessage } from '@/services/api.js';
  import { useRoute } from 'vue-router';
  
  export default {
    name: 'ChatView',
    setup() {
      const route = useRoute();
      // Get the receiver's ID from the query parameter (set by UserList.vue).
      const receiverId = ref(route.query.receiverId || "");
      const currentUserId = localStorage.getItem("userID") || "";
      const messages = ref([]);
      const newMessage = ref("");
      const chatError = ref("");
  
      // In a real app, you might load existing conversation messages by calling a "getConversation" API.
      // For simplicity, we start with an empty conversation.
      async function handleSendMessage() {
        chatError.value = "";
        if (!newMessage.value.trim()) return;
        try {
          // For private messages, we send with isGroup set to false.
          const response = await sendMessage({
            receiverId: receiverId.value,
            content: newMessage.value,
            isGroup: false,
            groupId: ""
          });
          // Append the new message to the messages array.
          messages.value.push({
            id: response.data.messageId,
            sender_id: currentUserId,
            content: newMessage.value,
            sent_at: new Date().toISOString()
          });
          newMessage.value = "";
        } catch (err) {
          console.error(err);
          chatError.value = "Failed to send message";
        }
      }
  
      function formatTimestamp(ts) {
        return new Date(ts).toLocaleTimeString();
      }
  
      // Optionally, load existing messages if your backend supports it.
      // onMounted(() => {
      //   loadConversationMessages();
      // });
  
      return {
        receiverId,
        currentUserId,
        messages,
        newMessage,
        chatError,
        handleSendMessage,
        formatTimestamp
      };
    }
  };
  </script>
  
  <style scoped>
  .chat-view {
    display: flex;
    flex-direction: column;
    height: 100vh;
    font-family: sans-serif;
  }
  .chat-header {
    padding: 10px;
    background-color: #27ae60;
    color: white;
    text-align: center;
  }
  .chat-messages {
    flex: 1;
    padding: 10px;
    background-color: #e9ecef;
    overflow-y: auto;
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
  .chat-input {
    display: flex;
    padding: 10px;
    background-color: #fff;
    border-top: 1px solid #ccc;
  }
  .chat-input input {
    flex: 1;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  .chat-input button {
    margin-left: 10px;
    padding: 10px 20px;
    background-color: #27ae60;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  .error {
    color: red;
    margin-top: 10px;
  }
  </style>
  