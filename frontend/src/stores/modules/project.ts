import { defineStore } from 'pinia';
import { ref } from 'vue';

import { projectApi } from '@/services/modules/project';
import type { ProjectHistoryItem, ProjectPolishFormData, ProjectPolishOutput } from '@/types/project';

const defaultForm = (): ProjectPolishFormData => ({
  projectName: 'OfferPilot',
  role: '前端负责人',
  techStack: 'Vue 3 + TypeScript + Vite + Pinia',
  projectBackground: '一个面向前端求职者的智能训练平台',
  responsibility: '前端工程架构搭建、核心页面设计与状态管理方案落地',
});

export const useProjectStore = defineStore('project', () => {
  const loading = ref(false);
  const form = ref<ProjectPolishFormData>(defaultForm());
  const result = ref<ProjectPolishOutput | null>(null);
  const historyList = ref<ProjectHistoryItem[]>([
    {
      id: 'history-1',
      projectName: 'OfferPilot',
      createdAt: '今天 16:30',
      summary: '前端面试训练平台包装版本',
    },
    {
      id: 'history-2',
      projectName: '企业级低代码平台',
      createdAt: '昨天 20:15',
      summary: '强调组件体系与可配置能力',
    },
  ]);

  const setForm = (payload: ProjectPolishFormData) => {
    form.value = payload;
  };

  const generate = async () => {
    loading.value = true;

    try {
      const data = await projectApi.generateProjectPolish(form.value);
      result.value = data;

      historyList.value = [
        {
          id: `history-${Date.now()}`,
          projectName: form.value.projectName,
          createdAt: '刚刚',
          summary: data.resumeDescription.slice(0, 28),
        },
        ...historyList.value,
      ];
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    form,
    result,
    historyList,
    setForm,
    generate,
  };
});