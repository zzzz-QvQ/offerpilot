import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { questionBankApi } from '@/services/modules/question-bank';
import type { QuestionBankFilters, QuestionCategoryOption, QuestionItem } from '@/types/question-bank';

const defaultFilters = (): QuestionBankFilters => ({
  keyword: '',
  category: '',
  frequency: '',
});

export const useQuestionBankStore = defineStore('questionBank', () => {
  const loading = ref(false);
  const questions = ref<QuestionItem[]>([]);
  const filters = ref<QuestionBankFilters>(defaultFilters());
  const selectedQuestionId = ref('');
  const detailVisible = ref(false);

  const categoryOptions = computed<QuestionCategoryOption[]>(() => {
    const categories = Array.from(new Set(questions.value.map((item) => item.category)));
    return categories.map((category) => ({
      label: category,
      value: category,
    }));
  });

  const filteredQuestions = computed(() => {
    return questions.value.filter((item) => {
      const matchKeyword = !filters.value.keyword || item.title.toLowerCase().includes(filters.value.keyword.toLowerCase());
      const matchCategory = !filters.value.category || item.category === filters.value.category;
      const matchFrequency = !filters.value.frequency || item.frequency === filters.value.frequency;

      return matchKeyword && matchCategory && matchFrequency;
    });
  });

  const selectedQuestion = computed(() => {
    return questions.value.find((item) => item.id === selectedQuestionId.value) ?? null;
  });

  const fetchQuestions = async () => {
    loading.value = true;

    try {
      const data = await questionBankApi.fetchQuestionList();
      questions.value = data;

      if (!selectedQuestionId.value && data.length > 0) {
        selectedQuestionId.value = data[0].id;
      }
    } finally {
      loading.value = false;
    }
  };

  const setFilters = (payload: Partial<QuestionBankFilters>) => {
    filters.value = {
      ...filters.value,
      ...payload,
    };
  };

  const resetFilters = () => {
    filters.value = defaultFilters();
  };

  const openQuestionDetail = (questionId: string) => {
    selectedQuestionId.value = questionId;
    detailVisible.value = true;
  };

  const closeQuestionDetail = () => {
    detailVisible.value = false;
  };

  return {
    loading,
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
  };
});