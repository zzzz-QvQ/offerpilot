import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { authApi } from '@/api/auth';
import { STORAGE_KEYS } from '@/constants/storage';
import type { LoginParams, UserProfile } from '@/types/user';

const readTokenFromStorage = () => {
  const value = localStorage.getItem(STORAGE_KEYS.token);

  if (!value || value === 'undefined' || value === 'null') {
    localStorage.removeItem(STORAGE_KEYS.token);
    return '';
  }

  return value;
};

const readUserFromStorage = (): UserProfile | null => {
  const raw = localStorage.getItem(STORAGE_KEYS.user);

  if (!raw) {
    return null;
  }

  try {
    return JSON.parse(raw) as UserProfile;
  } catch {
    localStorage.removeItem(STORAGE_KEYS.user);
    return null;
  }
};

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(readTokenFromStorage());
  const userInfo = ref<UserProfile | null>(readUserFromStorage());
  const initializing = ref(false);
  const loggingOut = ref(false);
  const isLoggedIn = computed(() => Boolean(token.value) && Boolean(userInfo.value));

  const clearAuthState = () => {
    token.value = '';
    userInfo.value = null;
    localStorage.removeItem(STORAGE_KEYS.token);
    localStorage.removeItem(STORAGE_KEYS.user);
  };

  const setToken = (value: string) => {
    token.value = value;
    localStorage.setItem(STORAGE_KEYS.token, value);
  };

  const setUserInfo = (value: UserProfile | null) => {
    userInfo.value = value;

    if (value) {
      localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(value));
      return;
    }

    localStorage.removeItem(STORAGE_KEYS.user);
  };

  const fetchCurrentUser = async () => {
    const result = await authApi.getProfile();
    setUserInfo(result.data);
    return result.data;
  };

  const restoreSession = async () => {
    if (!token.value) {
      clearAuthState();
      return null;
    }

    initializing.value = true;

    try {
      return await fetchCurrentUser();
    } catch {
      clearAuthState();
      return null;
    } finally {
      initializing.value = false;
    }
  };

  const login = async (params: LoginParams) => {
    const result = await authApi.login(params);

    setToken(result.data.token);
    setUserInfo(result.data.user);

    return result.data;
  };

  const logout = async () => {
    loggingOut.value = true;

    try {
      if (token.value) {
        await authApi.logout();
      }
    } finally {
      clearAuthState();
      loggingOut.value = false;
    }
  };

  return {
    token,
    userInfo,
    initializing,
    loggingOut,
    isLoggedIn,
    setToken,
    setUserInfo,
    fetchCurrentUser,
    restoreSession,
    login,
    logout,
    clearAuthState,
  };
});
