import { http } from '@/utils/request';
import type { LoginParams, LoginResponse, UserProfile } from '@/types/user';

export const authApi = {
  login(data: LoginParams) {
    return http.post<LoginResponse>('/auth/login', data);
  },
  logout() {
    return http.post<void>('/auth/logout');
  },
  getProfile() {
    return http.get<UserProfile>('/auth/me');
  },
};