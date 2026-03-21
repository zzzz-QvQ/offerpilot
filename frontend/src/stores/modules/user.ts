import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { STORAGE_KEYS } from '@/constants/storage';
import { authApi } from '@/services/modules/auth';
import type { LoginParams, UserProfile } from '@/types/user';

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
  const token = ref<string>(localStorage.getItem(STORAGE_KEYS.token) ?? '');
  const userInfo = ref<UserProfile | null>(readUserFromStorage());
  const isLoggedIn = computed(() => Boolean(token.value));

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

  const login = async (params: LoginParams) => {
    const result = await authApi.login(params);

    setToken(result.token);
    setUserInfo(result.user);

    return result;
  };

  const logout = () => {
    token.value = '';
    userInfo.value = null;
    localStorage.removeItem(STORAGE_KEYS.token);
    localStorage.removeItem(STORAGE_KEYS.user);
  };

  return {
    token,
    userInfo,
    isLoggedIn,
    setToken,
    setUserInfo,
    login,
    logout,
  };
});