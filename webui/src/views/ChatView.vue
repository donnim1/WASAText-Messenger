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
        <div class="message-bubble" :data-message-id="msg.ID">
          <!-- New: Display sender username for group chat messages (if not your own) -->
          <div v-if="conversationIsGroup && msg.SenderID !== currentUserId" class="sender-name">
            {{ getSenderName(msg.SenderID) }}
          </div>
          
          <!-- Only show reply preview when message was sent as a reply -->
          <div v-if="msg.ReplyTo && msg.ReplyTo !== ''" class="inline-reply-preview">
            <small>In reply to: {{ getReplyContent(msg.ReplyTo) }}</small>
          </div>
          
          <!-- Render message content -->
          <div v-if="isHtmlContent(msg.Content)" v-html="msg.Content"></div>
          <div v-else-if="isForwardedImage(msg.Content)" class="forwarded-image">
            <div class="forward-caption">Forwarded from you:</div>
            <img :src="getForwardedImageSrc(msg.Content)" alt="Image message" class="sent-image" />
          </div>
          <div v-else-if="isImage(msg.Content)">
            <img :src="msg.Content" alt="Image message" class="sent-image" />
          </div>
          <div v-else>
            <p class="message-content">{{ msg.Content }}</p>
          </div>
          
          <!-- Display reactions if available -->
          <div v-if="msg.reactions && msg.reactions.length" class="message-reactions">
            <template v-for="(group, idx) in groupReactions(msg.reactions)" :key="idx">
              <span class="reaction">
                {{ group.reaction }}<span v-if="group.count > 1"> ({{ group.count }})</span>
                <small v-if="group.userNames.length" title="Reacted by: {{ group.userNames.join(', ') }}">
                  {{ group.userNames.join(', ') }}
                </small>
              </span>
            </template>
          </div>

          <span class="message-timestamp">{{ formatTimestamp(msg.SentAt) }}</span>
          
          <!-- Checkmarks for sent messages (only for messages you sent) -->
          <template v-if="msg.SenderID === currentUserId">
            <span class="message-status">
              <template v-if="msg.status === 'pending'">
                <!-- Pending (clock icon) -->
                <i class="fas fa-clock"></i>
              </template>
              <template v-else-if="msg.status === 'sent'">
                <!-- Single checkmark -->
                <i class="fas fa-check"></i>
              </template>
              <template v-else-if="msg.status === 'delivered'">
                <!-- Double checkmark for delivered -->
                <i class="fas fa-check-double"></i>
              </template>
              <template v-else-if="msg.status === 'read'">
                <!-- Double checkmark colored for read -->
                <i class="fas fa-check-double read"></i>
              </template>
            </span>
          </template>
          
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
      <div class="forward-modal-content">
        <div class="forward-modal-header">
          <h3>Forward Message</h3>
        </div>
        <div class="forward-modal-body">
          <input v-model="forwardSearchQuery" class="forward-search-input" placeholder="Search chats or users" />
          <ul class="conversation-list">
            <li v-if="filteredForwardTargets.length === 0" class="no-results">
              No conversations found
            </li>
            <li v-for="conv in filteredForwardTargets" 
                :key="conv.id"
                @click="selectForwardTarget(conv)"
                :class="['conversation-item', { selected: forwardTargetConversation && ((conv.id || conv.ID) === (forwardTargetConversation.id || forwardTargetConversation.ID)) }]">
              <span class="conversation-name">{{ conv.name || 'Unnamed Chat' }}</span>
              <span v-if="conv.is_group" class="group-badge">Group</span>
            </li>
          </ul>
          <div v-if="forwardTargetConversation" class="selected-info">
            Selected: <strong>{{ forwardTargetConversation.name }}</strong>
          </div>
        </div>
        <div class="forward-modal-footer">
          <button 
            @click="confirmForwardMessage" 
            :disabled="!isForwardEnabled" 
            class="btn forward-btn"
            style="cursor: pointer !important; position: relative !important; z-index: 9999 !important;">
            Forward
          </button>
          <button @click="closeForwardModal" class="btn cancel-btn">
            Cancel
          </button>
        </div>
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
  forwardMessageApi,
  commentMessage as commentMessageApi,
  uncommentMessage as uncommentMessageApi,
  deleteMessage as deleteMessageApi,
  uploadImage,
  getMyConversations,
  listUsers,
  updateMessageStatus
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
    const contacts = ref([]);
    const forwardTargetConversation = ref(null);
    const defaultPhoto = "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg";

    // For filtering conversation list
    const filteredForwardTargets = computed(() => {
      if (!forwardSearchQuery.value) return allForwardTargets.value;
      const query = forwardSearchQuery.value.toLowerCase();
      return allForwardTargets.value.filter(target =>
        target.name.toLowerCase().includes(query)
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

    const conversation = ref({});

    const conversationIsGroup = computed(() => {
      if (conversation.value && conversation.value.is_group !== undefined) {
        console.log("[conversationIsGroup] from conversation:", conversation.value.is_group);
        return conversation.value.is_group;
      }
      // Fallback to route query if conversation isn't loaded yet.
      const flag = route.query.isGroup === "true";
      console.log("[conversationIsGroup] from route query:", flag);
      return flag;
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
      console.log("Message to forward object:", message);
      console.log("Message ID exists:", message.id !== undefined);
      console.log("Message ID uppercase exists:", message.ID !== undefined);
      
      messageToForward.value = message;
      showForwardModal.value = true;
      forwardSearchQuery.value = "";
      forwardTargetConversation.value = null;
      loadForwardTargets();
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
        messages.value = response.data.messages || [];
        conversation.value = response.data.conversation || {};
        await nextTick();
        scrollToBottom();
        markVisibleMessagesAsRead();
      } catch (err) {
        console.error("Error loading messages:", err);
        chatError.value = "Failed to load messages";
        messages.value = []; // Fallback to empty array on error.
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
      // Poll every half second (500ms)
      messagePollingInterval = setInterval(async () => {
        if (conversationId.value) {
          await loadConversationMessages(conversationId.value);
        }
      }, 500); // Changed from 1000 to 500ms
    }

    function stopMessagePolling() {
      clearInterval(messagePollingInterval);
    }

    onMounted(() => {
      loadContacts();
      initializeChat();
      startMessagePolling();
    });

    onUnmounted(() => {
      stopMessagePolling();
    });

    watch(
      () => route.params.conversationId,
      async (newId, oldId) => {
        if (newId && newId !== oldId) {
          console.log("Loading conversation:", newId);
          await loadConversationMessages(newId);
        }
      },
      { immediate: true }
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
        chatError.value = "Please select a conversation to forward this message to";
        return;
      }
    
      // Determine the target conversation ID.
      let targetConversationId = forwardTargetConversation.value.id || forwardTargetConversation.value.ID;
      let newConversationCreated = false;
    
      // If the selected target is a contact (i.e. no existing conversation)
      if (forwardTargetConversation.value.isContact) {
        try {
          const response = await getConversationByReceiver(forwardTargetConversation.value.id);
          if (response.data && response.data.conversationId) {
            targetConversationId = response.data.conversationId;
          }
        } catch (error) {
          // If no conversation exists (404), create one by sending a message.
          if (error.response && error.response.status === 404) {
            try {
              // Prepend "forwarded from you:" to the original message content.
              const forwardedContent = `forwarded from you: ${messageToForward.value.Content || messageToForward.value.content}`;
              const payload = {
                conversationId: "", // empty to trigger new conversation creation
                receiverId: forwardTargetConversation.value.id,
                content: forwardedContent,
                isGroup: false,
                groupId: "",
                replyTo: ""
              };
              const res = await sendMessage(payload);
              if (res.data && res.data.conversationId) {
                targetConversationId = res.data.conversationId;
                newConversationCreated = true;
              } else {
                chatError.value = "Failed to create a new conversation for forwarding";
                return;
              }
            } catch (err) {
              chatError.value = "Failed to create a conversation for forwarding: " +
                                 (err.response?.data || err.message);
              return;
            }
          } else {
            chatError.value = "Error checking conversation with contact; please try again later.";
            return;
          }
        }
      }
    
      // If a new conversation was created, the forward message has been sent as the first message.
      if (newConversationCreated) {
        chatError.value = `Message forwarded to ${forwardTargetConversation.value.name}!`;
        router.push({
          name: "ChatView",
          params: { conversationId: targetConversationId },
          query: {
            receiverName: forwardTargetConversation.value.name,
            isGroup: forwardTargetConversation.value.isGroup,
          },
        });
        return;
      }
    
      // Otherwise, forward using the existing conversation.
      const messageId = messageToForward.value.ID || messageToForward.value.id;
      if (!messageId || !targetConversationId) {
        chatError.value = "Invalid message or conversation";
        return;
      }
    
      try {
        console.log(`Forwarding message ${messageId} to conversation ${targetConversationId}`);
        await forwardMessageApi(messageId, targetConversationId);
        chatError.value = `Message forwarded to ${forwardTargetConversation.value.name}!`;
        router.push({
          name: "ChatView",
          params: { conversationId: targetConversationId },
          query: {
            receiverName: forwardTargetConversation.value.name,
            isGroup: forwardTargetConversation.value.isGroup,
          },
        });
      } catch (error) {
        console.error("Forward Message Error:", error);
        chatError.value = "Failed to forward message: " + (error.response?.data || error.message);
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
      console.log("Selected conversation full object:", conv);
      console.log("ID property exists:", conv.id !== undefined);
      console.log("ID property value:", conv.id);
      console.log("ID property uppercase exists:", conv.ID !== undefined);
      console.log("ID property uppercase value:", conv.ID);
      
      forwardTargetConversation.value = conv;
      forwardSearchQuery.value = conv.name || "Selected conversation"; 
    }

    async function loadConversations() {
      try {
        const response = await getMyConversations();
        console.log("Raw conversations response:", response);
        
        if (response.data && Array.isArray(response.data.conversations)) {
          conversations.value = response.data.conversations;
          console.log("First conversation:", conversations.value[0]);
        } else if (response.data && Array.isArray(response.data)) {
          conversations.value = response.data;
          console.log("First conversation:", conversations.value[0]);
        } else {
          console.error("Unexpected conversations API response structure:", response);
          conversations.value = [];
        }
      } catch (err) {
        console.error("Error loading conversations:", err);
        conversations.value = [];
      }
    }

    async function loadContacts() {
      try {
        const response = await listUsers();
        console.log("[loadContacts] API response:", response.data);
        if (response.data && Array.isArray(response.data.users)) {
          contacts.value = response.data.users;
        } else {
          contacts.value = [];
          console.warn("[loadContacts] No users array in response.data:", response.data);
        }
        console.log("[loadContacts] Loaded contacts:", contacts.value);
      } catch (err) {
        console.error("[loadContacts] Error loading contacts:", err);
        contacts.value = [];
      }
    }

    async function loadForwardTargets() {
      await Promise.all([loadConversations(), loadContacts()]);
    }

    const allForwardTargets = computed(() => {
      // Map conversations to a target object.
      const convTargets = conversations.value.map(c => ({
        id: c.id || c.ID,
        name: c.name || "Unnamed Chat",
        isGroup: c.is_group // Existing conversation flag.
      }));
      // Map contacts to a target object.
      const contactTargets = contacts.value.map(u => ({
        id: u.id,
        name: u.username,
        isGroup: false,
        isContact: true
      }));
      // Optionally, you could remove duplicates if a conversation already exists with that contact.
      return convTargets.concat(contactTargets);
    });

    async function markMessagesAsRead() {
      const unread = messages.value.filter(
        (msg) => msg.SenderID !== currentUserId && msg.status !== 'read'
      );
      for (const msg of unread) {
        try {
          await updateMessageStatus(msg.ID, 'read');
          msg.status = 'read';
          // Optionally update msg.readAt if you return it from the API.
        } catch (err) {
          console.error(`Failed to update status for message ${msg.ID}:`, err);
        }
      }
    }

    async function markVisibleMessagesAsRead() {
      if (!messagesContainer.value) return;

      const observer = new IntersectionObserver(
        async (entries) => {
          for (const entry of entries) {
            if (entry.isIntersecting) {
              const msgId = entry.target.dataset.messageId;
              const msg = messages.value.find(m => m.ID === msgId);
              // Only call update if it wasn't already marked as read
              if (msg && msg.SenderID !== currentUserId && msg.status !== 'read') {
                try {
                  // Call the API which should update the 'read' status only once all users have seen the message.
                  const res = await updateMessageStatus(msg.ID, 'read');
                  // Check the response—for example, if res.data.allSeen is true then update the local status.
                  if (res.data && res.data.allSeen) {
                    msg.status = 'read';
                    // Optionally, update other properties returned by the API (like msg.readAt)
                  }
                } catch (err) {
                  console.error(`Failed to mark message ${msg.ID} as read:`, err);
                }
              }
            }
          }
        },
        {
          root: messagesContainer.value,
          threshold: 0.5,
        }
      );

      document.querySelectorAll(".message-bubble").forEach(el => {
        observer.observe(el);
      });
    }

    const groupReactions = (reactions) => {
      const groups = {};
      reactions.forEach(r => {
        const key = r.reaction;
        if (!groups[key]) {
          groups[key] = { reaction: key, count: 0, userNames: [] };
        }
        groups[key].count++;
        groups[key].userNames.push(r.userName);
      });
      return Object.values(groups);
    };

    const isForwardEnabled = computed(() => {
      const msgId = messageToForward.value ? (messageToForward.value.ID || messageToForward.value.id) : null;
      const targetId = forwardTargetConversation.value ? (forwardTargetConversation.value.id || forwardTargetConversation.value.ID) : null;
      return !!(msgId && targetId);
    });

    function isForwardedImage(content) {
      if (!content) return false;
      return content.startsWith("FORWARDED_IMAGE:");
    }

    function getForwardedImageSrc(content) {
      // Remove the marker "FORWARDED_IMAGE:" from the content.
      return content.substring("FORWARDED_IMAGE:".length);
    }

    function isHtmlContent(content) {
      return content.trim().startsWith('<div class="forward-caption">');
    }

    function getSenderName(senderId) {
      console.log("[getSenderName] Called for senderId:", senderId);
      if (senderId === currentUserId) {
        console.log("[getSenderName] Sender is current user. Returning empty string.");
        return "";
      }
      console.log("[getSenderName] Contacts:", contacts.value);
      // Try matching using both 'id' and 'ID'
      const foundContact = contacts.value.find(contact => (contact.id || contact.ID) === senderId);
      if (foundContact) {
        console.log("[getSenderName] Found contact:", foundContact, "Available keys:", Object.keys(foundContact));
        return foundContact.username || foundContact.userName || "Unknown";
      } else {
        console.warn("[getSenderName] No contact found for senderId:", senderId, "Available contacts:", contacts.value);
        return "Unknown";
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
      contacts,
      allForwardTargets, // merged array
      filteredForwardTargets, // list filtered by search query
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
      selectForwardTarget, // <-- Added here
      groupReactions, // <-- Added here
      isForwardEnabled, // <-- Added here
      loadForwardTargets, // <-- Added here
      markVisibleMessagesAsRead, // <-- Added here
      isForwardedImage, // <-- Added here
      getForwardedImageSrc, // <-- Added here
      isHtmlContent, // <-- Added here
      conversationIsGroup, // <-- Added here
      getSenderName, // <-- Added here
      loadContacts
    };
  },
};
</script>

<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.2.0/css/all.min.css" />

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

/* Updated message bubble styling */
.message-bubble {
  position: relative;
  padding: 12px 16px;
  border-radius: 20px;
  max-width: 70%;
  word-wrap: break-word;
  background-color: #ffffff; /* a clean white for received messages */
  border: 1px solid #e0e0e0; /* subtle border */
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08); /* softer shadow */
}

/* For sent messages, use a complementary background */
.message-wrapper.sent .message-bubble {
  background-color: #e7f3fe; /* softer blue */
  border: 1px solid #c2e0f4;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* Sender name styling */
.sender-name {
  font-size: 0.85rem;
  font-weight: bold;
  color: #007bff;
  margin-bottom: 4px;
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
  width: 100%;
}

.chat-input-form {
  display: flex;
  gap: 10px;
  align-items: center;
}

.chat-input-form input {
  flex: 1;                   /* Grow to fill available space */
  width: 100%;               /* Ensure full width */
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

/* Forward Modal Layout Changes */
.forward-modal-content {
  background: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  padding: 20px;
}
.forward-modal-header h3 {
  margin: 0 0 10px 0;
  text-align: center;
}
.forward-modal-body {
  flex: 1;
  overflow-y: auto;
}
.forward-search-input {
  width: 100%;
  padding: 10px;
  margin-bottom: 15px;
  border: 1px solid #ccc;
  border-radius: 8px;
}
.conversation-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.conversation-item {
  padding: 10px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.conversation-item:hover,
.conversation-item.selected {
  background-color: #f1f2f5;
}
.group-badge {
  background-color: #007bff;
  color: #fff;
  font-size: 0.75rem;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 10px;
}
.selected-info {
  margin-top: 15px;
  font-size: 1rem;
  text-align: center;
}
.forward-modal-footer {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}
.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}
.forward-btn {
  background-color: #007bff;
  color: #fff;
  z-index: 1001; /* This is correct */
}
.forward-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.cancel-btn {
  background-color: #ccc;
  color: #333;
}

.sent-image {
  max-width: 100%; /* Prevents image from overflowing message container */
  max-height: 250px; /* Sets a reasonable maximum height */
  border-radius: 8px; /* Rounds corners to match message style */
  object-fit: contain; /* Maintains aspect ratio */
  border: 1px solid rgba(0,0,0,0.1);
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  transition: opacity 0.3s ease;
  display: block;
  margin-left: auto;
  margin-right: auto;
}

.message-status {
  display: inline-block;
  margin-left: 6px;
  font-size: 0.8rem;
  color: #888;
}

/* Updated checkmark styling for neat look */
.message-status i {
  font-size: 1rem;            /* Slightly larger for clarity */
  margin-left: 2px;           /* Small margin */
  background: transparent;    /* Remove background & border */
  padding: 0;
  color: #a0a0a0;             /* Default grey for sent/delivered status */
}

.message-status .read {
  color: #34B7F1 !important;
}

.forward-caption {
  font-size: 0.9rem;
  color: #555;
  margin-bottom: 4px;
  text-align: center;
  width: 100%;
}

.forwarded-image {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.forwarded-image .sent-image {
  max-height: 150px;  /* Adjust height as needed */
  max-width: 80%;     /* Optionally, limit width for forwarded images */
}
</style>