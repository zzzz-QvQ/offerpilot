import axios, { AxiosError, type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios';
import { ElMessage } from 'element-plus';

import { STORAGE_KEYS } from '@/constants/storage';
import type { ApiResponse } from '@/types/api';

const createRequestInstance = (): AxiosInstance => {
  const instance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 15000,
  });

  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem(STORAGE_KEYS.token);

      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }

      return config;
    },
    (error) => Promise.reject(error),
  );

  instance.interceptors.response.use(
    (response: AxiosResponse<ApiResponse>) => {
      const payload = response.data;

      if (typeof payload?.code === 'number' && payload.code !== 0) {
        ElMessage.error(payload.message || '请求失败');
        return Promise.reject(new Error(payload.message || 'Request failed'));
      }

      return response;
    },
    (error: AxiosError) => {
      const message = error.response?.data && typeof error.response.data === 'object' && 'message' in error.response.data
        ? String(error.response.data.message)
        : error.message || '网络异常，请稍后重试';

      ElMessage.error(message);
      return Promise.reject(error);
    },
  );

  return instance;
};

export const request = createRequestInstance();

export const http = {
  get<T>(url: string, config?: AxiosRequestConfig) {
    return request.get<ApiResponse<T>, ApiResponse<T>>(url, config);
  },
  post<T>(url: string, data?: unknown, config?: AxiosRequestConfig) {
    return request.post<ApiResponse<T>, ApiResponse<T>>(url, data, config);
  },
  put<T>(url: string, data?: unknown, config?: AxiosRequestConfig) {
    return request.put<ApiResponse<T>, ApiResponse<T>>(url, data, config);
  },
  delete<T>(url: string, config?: AxiosRequestConfig) {
    return request.delete<ApiResponse<T>, ApiResponse<T>>(url, config);
  },
};