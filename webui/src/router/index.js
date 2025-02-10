import { createRouter, createWebHashHistory } from 'vue-router';
import Login from '../components/Login.vue';
import MyProfile from '../views/MyProfile.vue';
import UserList from '../views/UserList.vue';
import GroupView from '../views/GroupView.vue';
import ChatView from '../views/ChatView.vue';
import ChatApp from '../views/ChatApp.vue';

const routes = [
  { path: '/', component: Login },
  { path: '/myprofile', component: MyProfile },
  { path: '/users', component: UserList },
  { path: '/groups', component: GroupView },
  { path: '/chat', name: 'ChatView', component: ChatView },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
