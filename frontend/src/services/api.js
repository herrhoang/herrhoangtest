import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080'; // 根据实际后端地址修改

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 添加请求拦截器
api.interceptors.request.use(
  (config) => {
    // 这里可以添加认证token等
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 添加响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    // 这里可以统一处理错误
    return Promise.reject(error);
  }
);

// API 方法
export const apiService = {
  // 示例方法
  getData: () => api.get('/api/data'),
  // 这里可以添加更多API方法
};

export default api;
