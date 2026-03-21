import { http } from '@/utils/request';
import type { DashboardSummaryResponse, RecentSessionItem, WeakPointItem } from '@/types/dashboard';

export const dashboardApi = {
  getSummary() {
    return http.get<DashboardSummaryResponse>('/dashboard/summary');
  },
  getRecentSessions() {
    return http.get<RecentSessionItem[]>('/dashboard/recent-sessions');
  },
  getWeakPoints() {
    return http.get<WeakPointItem[]>('/dashboard/weak-points');
  },
};