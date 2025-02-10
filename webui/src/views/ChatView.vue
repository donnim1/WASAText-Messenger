<template>
  <div class="chat-view">
    <h1>Conversation</h1>
    <ul>
      <li v-for="msg in messages" :key="msg.id" class="message-item">
        <strong>{{ msg.sender }}:</strong> {{ msg.content }} <small>({{ msg.sent_at }})</small>
      </li>
    </ul>
    <form @submit.prevent="sendMessageHandler" class="message-form">
      <input v-model="messageContent" placeholder="Type your message..." required />
      <button type="submit">Send</button>
    </form>
    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { sendMessage } from '@/services/api.js';

export default {
  name: 'ChatView',
  setup() {
    const messages = ref([]);
    const messageContent = ref("");
    const error = ref("");

    // For demonstration, these values are hardcoded.
    // In a real application, these would come from the conversation context.
    const receiverId = "RECEIVER_UUID"; // Set this to the receiver's user ID
    const isGroup = false;              // Change to true for group messages
    const groupId = "";                 // Set the group ID if applicable

    async function sendMessageHandler() {
      error.value = "";
      try {
        const response = await sendMessage({
          receiverId,
          content: messageContent.value,
          isGroup,
          groupId,
        });
        // Append the new message to the messages array.
        messages.value.push({
          id: response.data.messageId,
          sender: "You",
          content: messageContent.value,
          sent_at: new Date().toISOString(),
        });
        messageContent.value = "";
      } catch (err) {
        error.value = "Failed to send message";
        console.error(err);
      }
    }

    onMounted(() => {
      // Optionally, fetch existing messages for the conversation here.
      // For now, this example just initializes an empty conversation.
      messages.value = [];
    });

    return {
      messages,
      messageContent,
      error,
      sendMessageHandler,
    };
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
.message-form {
  margin-top: 20px;
}
.message-form input {
  width: 80%;
  padding: 10px;
  margin-right: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.message-form button {
  padding: 10px;
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
