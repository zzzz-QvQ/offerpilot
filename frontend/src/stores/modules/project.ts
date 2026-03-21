import { defineStore } from 'pinia';
import { ref } from 'vue';

import { projectApi } from '@/api/project';
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
  const historyLoading = ref(false);
  const detailLoading = ref(false);
  const error = ref('');
  const selectedHistoryId = ref('');
  const form = ref<ProjectPolishFormData>(defaultForm());
  const result = ref<ProjectPolishOutput | null>(null);
  const historyList = ref<ProjectHistoryItem[]>([]);

  const setForm = (payload: ProjectPolishFormData) => {
    form.value = payload;
  };

  const fetchHistory = async () => {
    historyLoading.value = true;

    try {
      const response = await projectApi.getHistory();
      historyList.value = response.data;
    } finally {
      historyLoading.value = false;
    }
  };

  const fetchDetail = async (id: string) => {
    detailLoading.value = true;
    error.value = '';

    try {
      const response = await projectApi.getDetail(id);
      result.value = response.data;
      selectedHistoryId.value = id;
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载项目包装详情失败';
      throw err;
    } finally {
      detailLoading.value = false;
    }
  };

  const generate = async () => {
    loading.value = true;
    error.value = '';

    try {
      const response = await projectApi.generate(form.value);
      result.value = response.data;
      selectedHistoryId.value = response.data.id;
      await fetchHistory();
    } catch (err) {
      error.value = err instanceof Error ? err.message : '生成项目包装失败';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const initialize = async () => {
    error.value = '';

    try {
      await fetchHistory();

      if (historyList.value.length > 0) {
        await fetchDetail(historyList.value[0].id);
      } else {
        result.value = null;
        selectedHistoryId.value = '';
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '初始化项目包装数据失败';
    }
  };

  return {
    loading,
    historyLoading,
    detailLoading,
    error,
    selectedHistoryId,
    form,
    result,
    historyList,
    setForm,
    fetchHistory,
    fetchDetail,
    generate,
    initialize,
  };
});