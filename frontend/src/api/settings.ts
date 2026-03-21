import { http } from '@/utils/request';
import type { SettingsState } from '@/types/settings';

export const settingsApi = {
  getSettings() {
    return http.get<SettingsState>('/settings');
  },
  updateSettings(data: Partial<SettingsState>) {
    return http.put<SettingsState>('/settings', data);
  },
};