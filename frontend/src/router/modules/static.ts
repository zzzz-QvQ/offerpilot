import type { RouteRecordRaw } from 'vue-router';

import MainLayout from '@/layouts/MainLayout.vue';

export const staticRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: {
      title: '登录',
    },
  },
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: {
          title: '控制台',
        },
      },
      {
        path: 'question-bank',
        name: 'QuestionBank',
        component: () => import('@/views/question-bank/QuestionBankView.vue'),
        meta: {
          title: '题库学习',
        },
      },
      {
        path: 'interview',
        name: 'Interview',
        component: () => import('@/views/interview/InterviewView.vue'),
        meta: {
          title: '模拟面试',
        },
      },
      {
        path: 'review',
        name: 'Review',
        component: () => import('@/views/review/ReviewView.vue'),
        meta: {
          title: '复盘分析',
        },
      },
      {
        path: 'project-polish',
        name: 'ProjectPolish',
        component: () => import('@/views/project-polish/ProjectPolishView.vue'),
        meta: {
          title: '项目包装',
        },
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/settings/SettingsView.vue'),
        meta: {
          title: '系统设置',
        },
      },
    ],
  },
];
