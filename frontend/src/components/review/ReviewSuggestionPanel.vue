<template>
  <el-card class="review-suggestion-panel" shadow="never">
    <template #header>
      <div class="review-suggestion-panel__header">学习建议与薄弱点</div>
    </template>

    <div class="review-suggestion-panel__list">
      <div v-for="item in items" :key="`${item.type}-${item.title}`" class="review-suggestion-panel__item">
        <div class="review-suggestion-panel__top">
          <h3>{{ item.title }}</h3>
          <el-tag :type="item.type === '建议' ? 'success' : 'warning'" size="small">{{ item.type }}</el-tag>
        </div>
        <p>{{ item.description }}</p>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { ReviewSuggestionItem } from '@/types/review';

defineProps<{
  items: ReviewSuggestionItem[];
}>();
</script>

<style scoped lang="scss">
.review-suggestion-panel {
  height: 100%;
}

.review-suggestion-panel__header {
  font-weight: 600;
}

.review-suggestion-panel__list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.review-suggestion-panel__item {
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fcfdff;
}

.review-suggestion-panel__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.review-suggestion-panel__top h3 {
  margin: 0;
  font-size: 15px;
  color: var(--color-text-primary);
}

.review-suggestion-panel__item p {
  margin: 10px 0 0;
  line-height: 1.8;
  color: var(--color-text-secondary);
}
</style>