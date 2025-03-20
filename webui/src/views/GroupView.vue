<template>
  <div class="group-management">
    <div class="groups-container">
      <!-- Left Panel: Create and List Groups -->
      <div class="groups-panel">
        <div class="panel-header">
          <h1>Your Groups</h1>
          <button @click="showCreateGroupModal = true" class="create-button">
            Create New Group
          </button>
        </div>
        <div class="groups-list">
          <div v-if="groups.length === 0" class="empty-state">
            <div class="empty-icon">👥</div>
            <h3>No Groups Yet</h3>
            <p>Create a new group to get started</p>
          </div>
          <div v-else class="group-items">
            <div v-for="group in groups" :key="group.id" class="group-item">
              <div class="group-avatar">
                <img
                  :src="group.photoUrl ? group.photoUrl : defaultGroupPhoto"
                  :alt="group.name"
                  class="avatar-image"
                />
              </div>
              <div class="group-content">
                <div class="group-header">
                  <h3 class="group-name">{{ group.name }}</h3>
                  <div class="group-actions">
                    <button class="action-button edit" @click="openUpdateGroupModal(group)">
                      Edit
                    </button>
                    <button class="action-button add" @click="openAddUserModal(group.id)">
                      Add
                    </button>
                    <button class="action-button leave" @click="leaveGroupHandler(group.id)">
                      Leave
                    </button>
                    <button class="action-button members" @click="openMembersModal(group)">
                      Members
                    </button>
                  </div>
                </div>
              </div>
            </div> <!-- end group-item -->
          </div> <!-- end group-items -->
        </div> <!-- end groups-list -->
      </div> <!-- end groups-panel -->
    </div> <!-- end groups-container -->

    <!-- Create Group Modal -->
    <div class="modal" v-if="showCreateGroupModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Create New Group</h2>
          <button class="close-button" @click="showCreateGroupModal = false">&times;</button>
        </div>
        <form @submit.prevent="createGroupHandler" class="modal-form">
          <div class="form-group">
            <label for="group-name">Group Name</label>
            <input
              id="group-name"
              v-model="groupName"
              type="text"
              placeholder="Enter group name"
              required
            />
          </div>
          <div class="form-group">
            <!-- Changed from text input to file input -->
            <label for="group-photo-upload">Group Photo</label>
            <input
              id="group-photo-upload"
              type="file"
              accept="image/*"
              @change="handleGroupPhotoUpload"
            />
          </div>
          <div class="modal-actions">
            <button type="button" class="cancel-button" @click="showCreateGroupModal = false">
              Cancel
            </button>
            <button type="submit" class="submit-button">
              Create Group
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Update Group Modal -->
    <div class="modal" v-if="showUpdateModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Update Group</h2>
          <button class="close-button" @click="closeUpdateModal">&times;</button>
        </div>
        <form @submit.prevent="updateGroupHandler" class="modal-form">
          <div class="form-group">
            <label for="update-name">Group Name</label>
            <input
              id="update-name"
              v-model="updateGroupName"
              type="text"
              placeholder="Enter new group name"
              required
            />
          </div>
          <div class="form-group">
            <!-- Replace text input with a file input -->
            <label for="update-photo-upload">Group Photo</label>
            <input
              id="update-photo-upload"
              type="file"
              accept="image/*"
              @change="handleUpdateGroupPhotoUpload"
            />
            <!-- Optionally display the current photo URL -->
            <div v-if="updateGroupPhoto">
              <img :src="updateGroupPhoto" alt="Group Photo" class="avatar-image" />
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="cancel-button" @click="closeUpdateModal">
              Cancel
            </button>
            <button type="submit" class="submit-button">
              Update Group
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Add User Modal -->
    <div class="modal" v-if="showAddUserModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Add User to Group</h2>
          <button class="close-button" @click="closeAddUserModal">&times;</button>
        </div>
        <div class="modal-body">
          <input
            type="text"
            v-model="userSearchQuery"
            placeholder="Search users..."
            class="search-input"
          />
          <ul>
            <li v-for="user in filteredAvailableUsers" :key="user.id" @click="addUser(user)">
              {{ user.username }}
            </li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Members Modal -->
    <div class="modal" v-if="showMembersModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Group Members</h2>
          <button class="close-button" @click="showMembersModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <ul>
            <li v-for="member in groupMembers" :key="member.id">
              {{ member.username }}
            </li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Notifications -->
    <div v-if="message" class="notification success">{{ message }}</div>
    <div v-if="error" class="notification error">{{ error }}</div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from "vue";
import {
  createGroup,
  uploadGroupImage,
  leaveGroup,
  listUserGroups,
  listUsers,
  addUserToGroupByUsername,
  setGroupName,
  setGroupPhoto
} from "@/services/api.js";

export default {
  name: "GroupManagement",
  setup() {
    // Existing reactive state
    const groupName = ref("");
    const groups = ref([]);
    const groupPhotoFile = ref(null);
    const message = ref("");
    const error = ref("");
    const defaultGroupPhoto = ref("https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg");
    const showCreateGroupModal = ref(false);
    const showUpdateModal = ref(false);
    const showAddUserModal = ref(false);
    const availableUsers = ref([]);
    const userSearchQuery = ref("");
    const loggedInUserID = localStorage.getItem("userID") || "";
    const selectedGroup = ref(null);
    const updateGroupName = ref("");
    const updateGroupPhoto = ref("");

    // New state for displaying group members
    const showMembersModal = ref(false);
    const groupMembers = ref([]);

    let refreshInterval = null;

    const filteredAvailableUsers = computed(() => {
      return availableUsers.value.filter(user => 
        user.id !== loggedInUserID &&
        user.username.toLowerCase().includes(userSearchQuery.value.toLowerCase())
      );
    });

    async function createGroupHandler() {
      message.value = "";
      error.value = "";
      try {
        const groupResponse = await createGroup({
          groupName: groupName.value,
          groupPhoto: ""
        });
        message.value = "Group created successfully";
        showCreateGroupModal.value = false;
        const newGroupId = groupResponse.data.groupId || groupResponse.data.id;
        if (groupPhotoFile.value && newGroupId) {
          const formData = new FormData();
          formData.append("photo", groupPhotoFile.value);
          try {
            await uploadGroupImage(newGroupId, formData);
            message.value += " and group photo updated successfully";
          } catch (uploadErr) {
            error.value = "Group created but failed to upload photo";
            console.error(uploadErr);
          }
        }
        groupName.value = "";
        groupPhotoFile.value = null;
        await refreshGroups();
      } catch (err) {
        error.value = "Failed to create group";
        console.error(err);
      }
    }

    async function refreshGroups() {
      try {
        const response = await listUserGroups();
        groups.value = (response.data.groups || []).map(group => {
          return {
            id: group.id,
            name: group.name,
            is_group: group.is_group,
            created_at: group.created_at,
            group_photo: group.group_photo,
            photoUrl: group.group_photo || "",
            last_message_content: group.last_message_content,
            last_message_sent_at: group.last_message_sent_at,
            members: (group.members || []).map(member => ({
              id: member.ID, 
              username: member.Username, 
              photo_url:
                member.PhotoUrl && member.PhotoUrl.Valid
                  ? member.PhotoUrl.String
                  : ""
            }))
          };
        });
      } catch (err) {
        error.value = "Failed to load groups";
        console.error(err);
      }
    }

    async function leaveGroupHandler(groupId) {
      if (!confirm("Are you sure you want to leave this group?")) return;
      message.value = "";
      error.value = "";
      try {
        await leaveGroup(groupId);
        message.value = "Left group successfully";
        await refreshGroups();
      } catch (err) {
        error.value = "Failed to leave group";
        console.error(err);
      }
    }

    async function openAddUserModal(groupId) {
      message.value = "";
      error.value = "";
      selectedGroup.value = groups.value.find(g => g.id === groupId);
      try {
        const response = await listUsers();
        availableUsers.value = response.data.users.filter(u => u.id !== loggedInUserID) || [];
        showAddUserModal.value = true;
      } catch (err) {
        error.value = "Failed to load users";
        console.error(err);
      }
    }

    async function addUser(user) {
      message.value = "";
      error.value = "";
      try {
        await addUserToGroupByUsername(selectedGroup.value.id, user.username);
        message.value = `User ${user.username} added successfully`;
        showAddUserModal.value = false;
        await refreshGroups();
      } catch (err) {
        error.value = "Failed to add user to group";
        console.error(err);
      }
    }

    function closeAddUserModal() {
      showAddUserModal.value = false;
    }

    function openUpdateGroupModal(group) {
      selectedGroup.value = group;
      updateGroupName.value = group.name;
      updateGroupPhoto.value = group.photoUrl || "";
      showUpdateModal.value = true;
    }

    async function handleUpdateGroupPhotoUpload(event) {
      const file = event.target.files[0];
      if (!file || !selectedGroup.value || !selectedGroup.value.id) return;
      const formData = new FormData();
      formData.append("photo", file);
      try {
        const response = await uploadGroupImage(selectedGroup.value.id, formData);
        updateGroupPhoto.value = response.data.photoUrl;
        message.value = "Group photo updated successfully";
      } catch (err) {
        error.value = "Failed to upload group photo";
        console.error(err);
      }
    }

    async function updateGroupHandler() {
      message.value = "";
      error.value = "";
      try {
        await setGroupName(selectedGroup.value.id, { newName: updateGroupName.value });
        if (updateGroupPhoto.value) {
          await setGroupPhoto(selectedGroup.value.id, { photoUrl: updateGroupPhoto.value });
        }
        message.value = "Group updated successfully";
        showUpdateModal.value = false;
        await refreshGroups();
      } catch (err) {
        error.value = "Failed to update group";
        console.error(err);
      }
    }

    function closeUpdateModal() {
      showUpdateModal.value = false;
    }

    function handleGroupPhotoUpload(event) {
      const file = event.target.files[0];
      if (file) {
        groupPhotoFile.value = file;
      }
    }

    // New: Open members modal using selected group's members
    function openMembersModal(group) {
      // Use members array from the mapped group object.
      groupMembers.value = group.members || [];
      // Optionally set selectedGroup if you need extra detail in modal header.
      showMembersModal.value = true;
    }

    // Auto-refresh groups periodically
    onMounted(() => {
      refreshGroups();
      refreshInterval = setInterval(refreshGroups, 500);
    });

    onUnmounted(() => {
      if (refreshInterval) {
        clearInterval(refreshInterval);
      }
    });

    return {
      groupName,
      groups,
      message,
      error,
      defaultGroupPhoto,
      showCreateGroupModal,
      createGroupHandler,
      leaveGroupHandler,
      openAddUserModal,
      refreshGroups,
      showUpdateModal,
      updateGroupName,
      updateGroupPhoto,
      openUpdateGroupModal,
      updateGroupHandler,
      closeUpdateModal,
      showAddUserModal,
      availableUsers,
      userSearchQuery,
      filteredAvailableUsers,
      addUser,
      closeAddUserModal,
      handleUpdateGroupPhotoUpload,
      handleGroupPhotoUpload,
      // Expose members modal properties
      showMembersModal,
      groupMembers,
      openMembersModal,
    };
  },
};
</script>

<style scoped>
.group-management {
  height: 100%;
  background-color: #f8f9fa;
  overflow: hidden;
}

.groups-container {
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #ffffff;
}

.groups-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.panel-header {
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-header h1 {
  margin: 0;
  color: #212529;
  font-size: 1.5rem;
  font-weight: 600;
}

.create-button {
  padding: 8px 16px;
  background-color: #4dabf7;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.2s;
}

.create-button:hover {
  background-color: #3c99e6;
}

.groups-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.group-items {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

.group-item {
  background: #ffffff;
  border: 1px solid #e9ecef;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  gap: 16px;
  transition: all 0.2s ease;
}

.group-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.group-avatar {
  flex-shrink: 0;
}

.avatar-image {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  object-fit: cover;
}

.group-content {
  flex: 1;
  min-width: 0;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.group-name {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #212529;
}



.group-actions {
  display: flex;
  gap: 8px;
}

.action-button {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.action-button.edit {
  background-color: #4dabf7;
  color: white;
}

.action-button.add {
  background-color: #40c057;
  color: white;
}

.action-button.leave {
  background-color: #ff6b6b;
  color: white;
}

.action-button.members {
  background-color: #f59f00;
  color: white;
}

.action-button:hover {
  opacity: 0.9;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #6c757d;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.empty-state h3 {
  margin: 0 0 8px;
  color: #495057;
}

empty-state p {
  margin: 0;
  font-size: 0.9rem;
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
  max-width: 500px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.25rem;
  color: #212529;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #6c757d;
  cursor: pointer;
}

.modal-form {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 0.9rem;
  color: #495057;
}

.form-group input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 0.9rem;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}

.cancel-button,
.submit-button {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  cursor: pointer;
}

.cancel-button {
  background-color: #e9ecef;
  color: #495057;
}

.submit-button {
  background-color: #4dabf7;
  color: white;
}

/* Notification Styles */
.notification {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 12px 20px;
  border-radius: 8px;
  font-size: 0.9rem;
  z-index: 1000;
  animation: slideIn 0.3s ease-out;
}

.notification.success {
  background-color: #40c057;
  color: white;
}

.notification.error {
  background-color: #ff6b6b;
  color: white;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}
</style>