import { createRouter, createWebHashHistory } from 'vue-router';
import MyProfile from '../views/MyProfile.vue';
import UserList from '../views/UserList.vue';
import ChatView from '../views/ChatView.vue';
import GroupView from '../views/GroupView.vue';


const routes = [
  { path: '/Myprofile', component: MyProfile },
  { path: '/users', component: UserList }, // New route for listing users
  { path: '/chat/:conversationId', component: ChatView },
  { path: '/groups', component: GroupView },
  { path: '/users', component: UserList },
];



const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
