import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import type {
  DashboardStatItem,
  RecentSessionItem,
  RecommendTaskItem,
  WeakPointItem,
} from '@/types/dashboard';

export const useDashboardStore = defineStore('dashboard', () => {
  const totalSessions = ref(128);
  const averageScore = ref(82);
  const weeklySessions = ref(9);
  const weakestSkill = ref('性能优化');

  const recentSessions = ref<RecentSessionItem[]>([
    {
      id: 'session-1',
      title: 'Vue 3 组件通信模拟训练',
      score: 84,
      durationText: '18 分钟',
      finishedAt: '今天 10:20',
      tag: '框架基础',
    },
    {
      id: 'session-2',
      title: '浏览器缓存与网络追问',
      score: 76,
      durationText: '22 分钟',
      finishedAt: '昨天 20:10',
      tag: '浏览器原理',
    },
    {
      id: 'session-3',
      title: 'TypeScript 类型体操快问快答',
      score: 88,
      durationText: '15 分钟',
      finishedAt: '03-19 19:40',
      tag: 'TypeScript',
    },
    {
      id: 'session-4',
      title: '前端工程化场景问答',
      score: 79,
      durationText: '25 分钟',
      finishedAt: '03-18 21:05',
      tag: '工程化',
    },
  ]);

  const weakPoints = ref<WeakPointItem[]>([
    { name: '性能优化', value: 91 },
    { name: '浏览器原理', value: 84 },
    { name: '工程化', value: 76 },
    { name: '手写题', value: 72 },
    { name: '网络安全', value: 68 },
  ]);

  const recommendTasks = ref<RecommendTaskItem[]>([
    {
      id: 'task-1',
      title: '性能优化专题训练',
      description: '聚焦首屏加载、缓存策略和长任务治理，建议完成 3 轮专项面试。',
      level: '高优先级',
    },
    {
      id: 'task-2',
      title: '浏览器渲染机制复盘',
      description: '补强关键渲染路径、回流重绘和事件循环的表达完整度。',
      level: '中优先级',
    },
    {
      id: 'task-3',
      title: '手写 Promise 与防抖节流',
      description: '适合面试前快速巩固，提升代码题稳定性。',
      level: '巩固建议',
    },
  ]);

  const statCards = computed<DashboardStatItem[]>(() => [
    {
      label: '总训练次数',
      value: String(totalSessions.value),
      hint: '累计完成的训练场次',
    },
    {
      label: '平均分',
      value: `${averageScore.value} 分`,
      hint: '最近阶段综合表现',
    },
    {
      label: '本周训练次数',
      value: String(weeklySessions.value),
      hint: '近 7 天训练活跃度',
    },
    {
      label: '当前最大薄弱项',
      value: weakestSkill.value,
      hint: '建议优先安排专项强化',
    },
  ]);

  return {
    statCards,
    recentSessions,
    weakPoints,
    recommendTasks,
  };
});