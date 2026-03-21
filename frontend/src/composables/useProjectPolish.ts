import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { ElMessage } from 'element-plus';

import { useProjectStore } from '@/stores/modules/project';
import type { ProjectPolishFormData } from '@/types/project';

export const useProjectPolish = () => {
  const store = useProjectStore();
  const { loading, historyLoading, detailLoading, error, selectedHistoryId, form, result, historyList } = storeToRefs(store);

  const submitForm = async (payload: ProjectPolishFormData) => {
    store.setForm(payload);
    await store.generate();
    ElMessage.success('项目包装内容已生成');
  };

  const viewHistoryDetail = async (id: string) => {
    await store.fetchDetail(id);
  };

  const copyText = async (text: string) => {
    await navigator.clipboard.writeText(text);
    ElMessage.success('已复制到剪贴板');
  };

  onMounted(() => {
    if (!store.historyList.length && !store.loading && !store.historyLoading) {
      void store.initialize();
    }
  });

  return {
    loading,
    historyLoading,
    detailLoading,
    error,
    selectedHistoryId,
    form,
    result,
    historyList,
    submitForm,
    viewHistoryDetail,
    copyText,
  };
};