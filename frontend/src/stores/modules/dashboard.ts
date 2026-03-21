import { defineStore } from 'pinia';
import { ref } from 'vue';

import { dashboardApi } from '@/api/dashboard';
import type {
  DashboardStatItem,
  RecentSessionItem,
  RecommendTaskItem,
  WeakPointItem,
} from '@/types/dashboard';

export const useDashboardStore = defineStore('dashboard', () => {
  const loading = ref(false);
  const error = ref('');
  const statCards = ref<DashboardStatItem[]>([]);
  const recentSessions = ref<RecentSessionItem[]>([]);
  const weakPoints = ref<WeakPointItem[]>([]);
  const recommendTasks = ref<RecommendTaskItem[]>([]);

  const fetchDashboardData = async () => {
    loading.value = true;
    error.value = '';

    try {
      const [summaryRes, recentRes, weakPointsRes] = await Promise.all([
        dashboardApi.getSummary(),
        dashboardApi.getRecentSessions(),
        dashboardApi.getWeakPoints(),
      ]);

      statCards.value = summaryRes.data.statCards;
      recommendTasks.value = summaryRes.data.recommendTasks;
      recentSessions.value = recentRes.data;
      weakPoints.value = weakPointsRes.data;
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载 Dashboard 数据失败';
      statCards.value = [];
      recommendTasks.value = [];
      recentSessions.value = [];
      weakPoints.value = [];
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    error,
    statCards,
    recentSessions,
    weakPoints,
    recommendTasks,
    fetchDashboardData,
  };
});