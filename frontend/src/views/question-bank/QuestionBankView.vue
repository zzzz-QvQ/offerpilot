<template>
  <div class="question-bank-view">
    <QuestionFilterBar
      :filters="filters"
      :categories="categoryOptions"
      @change="updateFilters"
      @reset="resetFilters"
    />

    <QuestionList
      :items="filteredQuestions"
      :loading="loading"
      @view-detail="handleViewDetail"
    />

    <QuestionDetailDrawer
      :visible="detailVisible"
      :question="selectedQuestion"
      @close="closeDetail"
    />
  </div>
</template>

<script setup lang="ts">
import QuestionDetailDrawer from '@/components/question-bank/QuestionDetailDrawer.vue';
import QuestionFilterBar from '@/components/question-bank/QuestionFilterBar.vue';
import QuestionList from '@/components/question-bank/QuestionList.vue';
import { useQuestionBank } from '@/composables/useQuestionBank';

const {
  loading,
  filters,
  categoryOptions,
  filteredQuestions,
  selectedQuestion,
  detailVisible,
  updateFilters,
  resetFilters,
  handleViewDetail,
  closeDetail,
} = useQuestionBank();
</script>

<style scoped lang="scss">
.question-bank-view {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>