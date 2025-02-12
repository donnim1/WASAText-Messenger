<!-- filepath: /c:/Users/HP/OneDrive/Desktop/2ND SEM/WASAText/webui/src/views/GroupView.vue -->
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

    <!-- List of Groups with options -->
    <section class="group-list">
      <h2>Your Groups</h2>
      <div v-if="groups.length === 0">
        You are not a member of any groups.
      </div>
      <ul>
        <li v-for="group in groups" :key="group.id" class="group-item">
          <div class="group-info">
            <img :src="group.photo || defaultGroupPhoto" alt="Group Photo" class="group-photo" />
            <span class="group-name">{{ group.name }}</span>
          </div>
          <div class="group-actions">
            <button @click="openUpdateGroupModal(group)">Update Group</button>
            <button @click="promptAddUser(group.id)">Add User</button>
            <button @click="leaveGroupHandler(group.id)">Leave Group</button>
          </div>
        </li>
      </ul>
    </section>

    <div v-if="message" class="message">{{ message }}</div>
    <div v-if="error" class="error">{{ error }}</div>

    <!-- Modal for updating group name and photo -->
    <div class="modal" v-if="showUpdateModal">
      <div class="modal-content">
        <h2>Update Group</h2>
        <form @submit.prevent="updateGroupHandler">
          <input type="text" v-model="updateGroupName" placeholder="Enter new group name" required />
          <input type="text" v-model="updateGroupPhoto" placeholder="Enter new group photo URL (optional)" />
          <div class="modal-actions">
            <button type="submit">Update</button>
            <button type="button" @click="closeUpdateModal">Cancel</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import {
  createGroup,
  addUserToGroupByUsername,
  leaveGroup,
  listUserGroups,
  setGroupName,
  setGroupPhoto,
} from "@/services/api.js";

export default {
  name: "GroupManagement",
  setup() {
    const groupName = ref("");
    const groupPhoto = ref("");
    const groups = ref([]);
    const message = ref("");
    const error = ref("");
    // Update the default group photo URL:
    const defaultGroupPhoto = ref("https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg");

    // Fields for update modal
    const showUpdateModal = ref(false);
    const updateGroupName = ref("");
    const updateGroupPhoto = ref("");
    const selectedGroup = ref(null);

    async function createGroupHandler() {
      message.value = "";
      error.value = "";
      try {
        const response = await createGroup({
          groupName: groupName.value,
          groupPhoto: groupPhoto.value,
        });
        message.value = "Group created with ID: " + response.data.groupId;
        groupName.value = "";
        groupPhoto.value = "";
        refreshGroups();
      } catch (err) {
        error.value = "Failed to create group";
        console.error(err);
      }
    }

    async function refreshGroups() {
      try {
        const response = await listUserGroups();
        groups.value = response.data.groups || [];
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
      message.value = "";
      error.value = "";
      const targetUsername = prompt("Enter the username to add:");
      if (targetUsername) {
        try {
          const response = await addUserToGroupByUsername(groupId, targetUsername);
          message.value = response.data.message;
          refreshGroups();
        } catch (err) {
          error.value = "Failed to add user to group";
          console.error(err);
        }
      }
    }

    function openUpdateGroupModal(group) {
      selectedGroup.value = group;
      updateGroupName.value = group.name;
      updateGroupPhoto.value = group.photo || "";
      showUpdateModal.value = true;
    }

    async function updateGroupHandler() {
      message.value = "";
      error.value = "";
      try {
        // Payload now contains { newName: ... }
        await setGroupName(selectedGroup.value.id, { newName: updateGroupName.value });
        // Update group photo if provided; payload now contains { photoUrl: ... }
        if (updateGroupPhoto.value) {
          await setGroupPhoto(selectedGroup.value.id, { photoUrl: updateGroupPhoto.value });
        }
        message.value = "Group updated successfully";
        showUpdateModal.value = false;
        refreshGroups();
      } catch (err) {
        error.value = "Failed to update group";
        console.error(err);
      }
    }

    function closeUpdateModal() {
      showUpdateModal.value = false;
    }

    onMounted(() => {
      refreshGroups();
    });

    return {
      groupName,
      groupPhoto,
      groups,
      message,
      error,
      defaultGroupPhoto,
      createGroupHandler,
      leaveGroupHandler,
      promptAddUser,
      refreshGroups,
      showUpdateModal,
      updateGroupName,
      updateGroupPhoto,
      openUpdateGroupModal,
      updateGroupHandler,
      closeUpdateModal,
    };
  },
};
</script>

<style scoped>
.group-management {
  padding: 20px;
  max-width: 700px;
  margin: 20px auto;
  background: #2c2c2c;
  color: #e0e0e0;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.5);
}

h1,
h2 {
  color: #4caf50;
}

section {
  margin-bottom: 20px;
}

input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #555;
  border-radius: 4px;
  background-color: #444;
  color: #e0e0e0;
  font-size: 1rem;
}

button {
  background-color: #4caf50;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  margin: 5px 0;
}

button:hover {
  background-color: #43a047;
}

.group-list ul {
  list-style: none;
  padding: 0;
}

.group-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #555;
  padding: 10px 0;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.group-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.group-actions button {
  margin-right: 5px;
}

.message {
  color: #4caf50;
  margin-top: 10px;
  font-weight: bold;
}

.error {
  color: #f44336;
  margin-top: 10px;
  font-weight: bold;
}

/* Modal Styling */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: #333;
  padding: 20px;
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
}

.modal-content h2 {
  margin-top: 0;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>