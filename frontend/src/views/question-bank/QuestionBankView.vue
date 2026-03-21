<template>
  <div class="question-bank-view">
    <el-alert
      v-if="error"
      :title="error"
      type="error"
      show-icon
      closable="false"
    />

    <QuestionFilterBar
      :filters="filters"
      :categories="categoryOptions"
      @change="updateFilters"
      @reset="resetFilters"
    />

    <QuestionList
      :items="filteredQuestions"
      :loading="loading"
      :favorite-loading-id="favoriteLoadingId"
      @view-detail="handleViewDetail"
      @toggle-favorite="handleToggleFavorite"
    />

    <QuestionDetailDrawer
      :visible="detailVisible"
      :question="selectedQuestion"
      :loading="detailLoading"
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
  detailLoading,
  favoriteLoadingId,
  error,
  filters,
  categoryOptions,
  filteredQuestions,
  selectedQuestion,
  detailVisible,
  updateFilters,
  resetFilters,
  handleViewDetail,
  handleToggleFavorite,
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