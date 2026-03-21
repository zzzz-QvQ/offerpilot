import { defineStore } from 'pinia';
import { ref } from 'vue';

import type { SettingsState } from '@/types/settings';

const defaultSettings = (): SettingsState => ({
  defaultInterviewMode: '综合模拟',
  followUpIntensity: '中',
  scoringStrictness: '标准',
  modelName: 'gpt-4o-mini',
  streamOutput: true,
});

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<SettingsState>(defaultSettings());

  const updateSettings = (payload: Partial<SettingsState>) => {
    settings.value = {
      ...settings.value,
      ...payload,
    };
  };

  const resetSettings = () => {
    settings.value = defaultSettings();
  };

  return {
    settings,
    updateSettings,
    resetSettings,
  };
});