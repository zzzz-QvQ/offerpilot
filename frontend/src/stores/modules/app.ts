import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false);

  const setSidebarCollapsed = (value: boolean) => {
    sidebarCollapsed.value = value;
  };

  const reset = () => {
    sidebarCollapsed.value = false;
  };

  return {
    sidebarCollapsed,
    setSidebarCollapsed,
    reset,
  };
});
