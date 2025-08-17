import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  withCredentials: true,
  headers: {
    'Accept': 'application/json',
  },
});


api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      window.dispatchEvent(new CustomEvent('unauthenticated'));
    }

    return Promise.reject(error);
  }
);

export default api;