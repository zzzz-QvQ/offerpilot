import { http } from '@/utils/request';
import type { DashboardStatItem, RecentSessionItem, RecommendTaskItem, WeakPointItem } from '@/types/dashboard';

export interface DashboardOverview {
  statCards: DashboardStatItem[];
  recentSessions: RecentSessionItem[];
  weakPoints: WeakPointItem[];
  recommendTasks: RecommendTaskItem[];
}

export const dashboardApi = {
  getOverview() {
    return http.get<DashboardOverview>('/dashboard/overview');
  },
};