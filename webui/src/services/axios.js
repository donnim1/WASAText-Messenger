import axios from "axios";

const instance = axios.create({
  baseURL: __API_URL__,
  timeout: 1000 * 5,
});

// Add a request interceptor to attach the Authorization header.
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("userID");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export default instance;
