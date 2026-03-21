export interface DashboardStatItem {
  label: string;
  value: string;
  hint: string;
}

export interface RecentSessionItem {
  id: string;
  title: string;
  score: number;
  durationText: string;
  finishedAt: string;
  tag: string;
}

export interface WeakPointItem {
  name: string;
  value: number;
}

export interface RecommendTaskItem {
  id: string;
  title: string;
  description: string;
  level: '高优先级' | '中优先级' | '巩固建议';
}

export interface DashboardSummaryResponse {
  statCards: DashboardStatItem[];
  recommendTasks: RecommendTaskItem[];
}