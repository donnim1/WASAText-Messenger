// src/router/index.js
import { createRouter, createWebHashHistory } from 'vue-router';
import MyProfile from '../views/MyProfile.vue';
import UserList from '../views/UserList.vue';
import GroupView from '../views/GroupView.vue';
import ChatView from '../views/ChatView.vue';
import MyChatsView from '../views/MyChatsView.vue';
import Login from '../components/Login.vue';

const routes = [
  {path:  '/login' , component:  Login },
  { path: '/myprofile', component: MyProfile },
  { path: '/users', component: UserList },
  { path: '/groups', component: GroupView },
  { path: '/chats', name: 'MyChatsView', component: MyChatsView },
  // conversationId is optional here – if not provided, ChatView will use query parameters.
  { path: '/chat/:conversationId?', name: 'ChatView', component: ChatView },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
