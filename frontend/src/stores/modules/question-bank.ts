import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { questionApi } from '@/api/question';
import type { QuestionBankFilters, QuestionCategoryOption, QuestionItem } from '@/types/question-bank';

const defaultFilters = (): QuestionBankFilters => ({
  keyword: '',
  category: '',
  frequency: '',
});

export const useQuestionBankStore = defineStore('questionBank', () => {
  const loading = ref(false);
  const detailLoading = ref(false);
  const favoriteLoadingId = ref('');
  const error = ref('');
  const questions = ref<QuestionItem[]>([]);
  const filters = ref<QuestionBankFilters>(defaultFilters());
  const selectedQuestionId = ref('');
  const selectedQuestion = ref<QuestionItem | null>(null);
  const detailVisible = ref(false);

  const categoryOptions = computed<QuestionCategoryOption[]>(() => {
    const categories = Array.from(new Set(questions.value.map((item) => item.category)));
    return categories.map((category) => ({ label: category, value: category }));
  });

  const filteredQuestions = computed(() => questions.value);

  const fetchQuestions = async () => {
    loading.value = true;
    error.value = '';

    try {
      const response = await questionApi.getList({
        keyword: filters.value.keyword || undefined,
        category: filters.value.category || undefined,
        frequency: filters.value.frequency || undefined,
        page: 1,
        pageSize: 100,
      });

      questions.value = response.data.list;

      if (selectedQuestionId.value) {
        const current = response.data.list.find((item) => item.id === selectedQuestionId.value) ?? null;
        if (current) {
          selectedQuestion.value = current;
        }
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '加载题库失败';
      questions.value = [];
    } finally {
      loading.value = false;
    }
  };

  const fetchQuestionDetail = async (questionId: string) => {
    detailLoading.value = true;

    try {
      const response = await questionApi.getDetail(questionId);
      selectedQuestion.value = response.data;
      selectedQuestionId.value = questionId;
    } finally {
      detailLoading.value = false;
    }
  };

  const setFilters = async (payload: Partial<QuestionBankFilters>) => {
    filters.value = {
      ...filters.value,
      ...payload,
    };

    await fetchQuestions();
  };

  const resetFilters = async () => {
    filters.value = defaultFilters();
    await fetchQuestions();
  };

  const openQuestionDetail = async (questionId: string) => {
    detailVisible.value = true;
    await fetchQuestionDetail(questionId);
  };

  const closeQuestionDetail = () => {
    detailVisible.value = false;
  };

  const toggleFavorite = async (question: QuestionItem) => {
    favoriteLoadingId.value = question.id;

    try {
      const response = question.isFavorite
        ? await questionApi.removeFavorite(question.id)
        : await questionApi.addFavorite(question.id);

      const nextFavorite = response.data.isFavorite;
      questions.value = questions.value.map((item) => (
        item.id === question.id ? { ...item, isFavorite: nextFavorite } : item
      ));

      if (selectedQuestion.value?.id === question.id) {
        selectedQuestion.value = {
          ...selectedQuestion.value,
          isFavorite: nextFavorite,
        };
      }
    } finally {
      favoriteLoadingId.value = '';
    }
  };

  return {
    loading,
    detailLoading,
    favoriteLoadingId,
    error,
    questions,
    filters,
    detailVisible,
    categoryOptions,
    filteredQuestions,
    selectedQuestion,
    fetchQuestions,
    setFilters,
    resetFilters,
    openQuestionDetail,
    closeQuestionDetail,
    toggleFavorite,
  };
});