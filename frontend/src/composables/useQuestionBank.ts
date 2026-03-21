import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';

import { useQuestionBankStore } from '@/stores/modules/question-bank';
import type { QuestionBankFilters, QuestionItem } from '@/types/question-bank';

export const useQuestionBank = () => {
  const store = useQuestionBankStore();
  const {
    loading,
    detailLoading,
    favoriteLoadingId,
    error,
    filters,
    categoryOptions,
    filteredQuestions,
    selectedQuestion,
    detailVisible,
  } = storeToRefs(store);

  const updateFilters = async (payload: Partial<QuestionBankFilters>) => {
    await store.setFilters(payload);
  };

  const handleViewDetail = async (questionId: string) => {
    await store.openQuestionDetail(questionId);
  };

  const handleToggleFavorite = async (question: QuestionItem) => {
    await store.toggleFavorite(question);
  };

  onMounted(() => {
    if (!store.questions.length) {
      void store.fetchQuestions();
    }
  });

  return {
    loading,
    detailLoading,
    favoriteLoadingId,
    error,
    filters,
    categoryOptions,
    filteredQuestions,
    selectedQuestion,
    detailVisible,
    updateFilters,
    resetFilters: store.resetFilters,
    handleViewDetail,
    handleToggleFavorite,
    closeDetail: store.closeQuestionDetail,
  };
};