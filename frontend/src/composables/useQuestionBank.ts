import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';

import { useQuestionBankStore } from '@/stores/modules/question-bank';
import type { QuestionBankFilters } from '@/types/question-bank';

export const useQuestionBank = () => {
  const store = useQuestionBankStore();
  const { loading, filters, categoryOptions, filteredQuestions, selectedQuestion, detailVisible } = storeToRefs(store);

  const updateFilters = (payload: Partial<QuestionBankFilters>) => {
    store.setFilters(payload);
  };

  const handleViewDetail = (questionId: string) => {
    store.openQuestionDetail(questionId);
  };

  onMounted(() => {
    if (!store.questions.length) {
      void store.fetchQuestions();
    }
  });

  return {
    loading,
    filters,
    categoryOptions,
    filteredQuestions,
    selectedQuestion,
    detailVisible,
    updateFilters,
    resetFilters: store.resetFilters,
    handleViewDetail,
    closeDetail: store.closeQuestionDetail,
  };
};