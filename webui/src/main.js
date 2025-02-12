import { createApp, reactive } from 'vue';
import App from './App.vue';
import router from './router';
import axios from './services/axios.js'; // Use your custom axios instance
import ErrorMsg from './components/ErrorMsg.vue';
import LoadingSpinner from './components/LoadingSpinner.vue';
import Login from './components/Login.vue';

import './assets/dashboard.css';
import './assets/main.css';

const app = createApp(App);
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Login", Login);

app.use(router);
app.mount('#app');
