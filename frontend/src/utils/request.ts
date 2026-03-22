import axios, { AxiosError, type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios';
import { ElMessage } from 'element-plus';

import { STORAGE_KEYS } from '@/constants/storage';
import type { ApiResponse } from '@/types/api';

const normalizeApiBaseURL = () => {
  const rawBaseURL = String(import.meta.env.VITE_API_BASE_URL ?? '').trim();

  if (!rawBaseURL) {
    return '/api';
  }

  return rawBaseURL.endsWith('/api')
    ? rawBaseURL
    : `${rawBaseURL.replace(/\/$/, '')}/api`;
};

const getStoredToken = () => {
  const token = localStorage.getItem(STORAGE_KEYS.token);

  if (!token || token === 'undefined' || token === 'null') {
    localStorage.removeItem(STORAGE_KEYS.token);
    return '';
  }

  return token;
};

const createRequestInstance = (): AxiosInstance => {
  const instance = axios.create({
    baseURL: normalizeApiBaseURL(),
    timeout: 15000,
  });

  instance.interceptors.request.use(
    (config) => {
      const token = getStoredToken();

      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }

      return config;
    },
    (error) => Promise.reject(error),
  );

  instance.interceptors.response.use(
    (response: AxiosResponse<ApiResponse<unknown>>) => {
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

      if (error.response?.status === 401) {
        localStorage.removeItem(STORAGE_KEYS.token);
        localStorage.removeItem(STORAGE_KEYS.user);
        window.location.href = '/login';
        return Promise.reject(error);
      }

      ElMessage.error(message);
      return Promise.reject(error);
    },
  );

  return instance;
};

export const request = createRequestInstance();

export const http = {
  async get<T>(url: string, config?: AxiosRequestConfig) {
    const response = await request.get<ApiResponse<T>>(url, config);
    return response.data as ApiResponse<T>;
  },
  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig) {
    const response = await request.post<ApiResponse<T>>(url, data, config);
    return response.data as ApiResponse<T>;
  },
  async put<T>(url: string, data?: unknown, config?: AxiosRequestConfig) {
    const response = await request.put<ApiResponse<T>>(url, data, config);
    return response.data as ApiResponse<T>;
  },
  async delete<T>(url: string, config?: AxiosRequestConfig) {
    const response = await request.delete<ApiResponse<T>>(url, config);
    return response.data as ApiResponse<T>;
  },
};
