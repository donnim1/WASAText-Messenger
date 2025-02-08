<!-- src/views/GroupManagement.vue -->
<template>
    <div class="group-management">
      <h1>Group Management</h1>
      
      <!-- Create Group Section -->
      <section class="create-group">
        <h2>Create Group</h2>
        <form @submit.prevent="createGroupHandler">
          <input v-model="groupName" placeholder="Enter group name" required />
          <input v-model="groupPhoto" placeholder="Enter group photo URL (optional)" />
          <button type="submit">Create Group</button>
        </form>
      </section>
  
      <!-- List of Groups with Leave and Add User Options -->
      <section class="group-list">
        <h2>Your Groups</h2>
        <div v-if="!groups">You are not a member of any groups.</div>
        <ul>
          <li v-for="group in groups" :key="group.id" class="group-item">
            <span>{{ group.name || "Private Chat" }}</span>
            <!-- Leave Group Button -->
            <button @click="leaveGroupHandler(group.id)">Leave Group</button>
            <!-- Add User Button -->
            <button @click="promptAddUser(group.id)">Add User</button>
          </li>
        </ul>
      </section>
  
      <div v-if="message" class="message">{{ message }}</div>
      <div v-if="error" class="error">{{ error }}</div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { createGroup, addUserToGroup, leaveGroup, listUserGroups } from '@/services/api.js';

  
  export default {
    name: 'GroupManagement',
    setup() {
      const groupName = ref("");
      const groupPhoto = ref("");
      const groups = ref([]);
      const message = ref("");
      const error = ref("");
  
      async function createGroupHandler() {
        message.value = "";
        error.value = "";
        try {
          const response = await createGroup({ groupName: groupName.value, groupPhoto: groupPhoto.value });
          message.value = "Group created with ID: " + response.data.groupId;
          // Refresh the list of groups
          refreshGroups();
        } catch (err) {
          error.value = "Failed to create group";
          console.error(err);
        }
      }
  
      async function refreshGroups() {
        try {
          const response = await listUserGroups();
          groups.value = response.data.groups; // Assumes backend returns a list of groups in { groups: [...] }
        } catch (err) {
          error.value = "Failed to load groups";
          console.error(err);
        }
      }
  
      async function leaveGroupHandler(groupId) {
        message.value = "";
        error.value = "";
        try {
          const response = await leaveGroup(groupId);
          message.value = response.data.message;
          refreshGroups();
        } catch (err) {
          error.value = "Failed to leave group";
          console.error(err);
        }
      }
  
      async function promptAddUser(groupId) {
        const targetUserId = prompt("Enter the user ID to add:");
        if (targetUserId) {
          try {
            // Call API to add another user to the group.
            const response = await addUserToGroup(groupId, targetUserId);
            message.value = response.data.message;
            refreshGroups();
          } catch (err) {
            error.value = "Failed to add user to group";
            console.error(err);
          }
        }
      }
  
      onMounted(() => {
        refreshGroups();
      });
  
      return { groupName, groupPhoto, groups, message, error, createGroupHandler, leaveGroupHandler, promptAddUser };
    },
  };
  </script>
  
  <style scoped>
  .group-management {
    padding: 20px;
    max-width: 600px;
    margin: 20px auto;
    background: #e8f5e9;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }
  section {
    margin-bottom: 20px;
  }
  input {
    width: 100%;
    padding: 10px;
    margin: 10px 0;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 1rem;
  }
  button {
    background-color: #43a047;
    color: white;
    padding: 10px 15px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    margin-right: 5px;
  }
  button:hover {
    background-color: #388e3c;
  }
  .message {
    color: green;
    margin-top: 10px;
    font-weight: bold;
  }
  .error {
    color: red;
    margin-top: 10px;
    font-weight: bold;
  }
  .group-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px;
    border-bottom: 1px solid #ccc;
  }
  </style>
  