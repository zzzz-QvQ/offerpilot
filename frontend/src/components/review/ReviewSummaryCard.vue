<template>
  <el-card class="review-summary-card" shadow="never">
    <div class="review-summary-card__score">{{ score }}</div>
    <div class="review-summary-card__label">总体评分</div>
    <div class="review-summary-card__list">
      <div v-for="item in items" :key="item.label" class="review-summary-card__item">
        <div class="review-summary-card__item-label">{{ item.label }}</div>
        <div class="review-summary-card__item-value">{{ item.value }}</div>
        <div class="review-summary-card__item-hint">{{ item.hint }}</div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { ReviewSummaryItem } from '@/types/review';

defineProps<{
  score: number;
  items: ReviewSummaryItem[];
}>();
</script>

<style scoped lang="scss">
.review-summary-card {
  height: 100%;
  border: 1px solid rgba(37, 99, 235, 0.08);
}

.review-summary-card__score {
  font-size: 42px;
  font-weight: 700;
  line-height: 1;
  color: var(--color-brand);
}

.review-summary-card__label {
  margin-top: 8px;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.review-summary-card__list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 24px;
}

.review-summary-card__item {
  padding: 14px;
  border-radius: 12px;
  background: #f8fbff;
}

.review-summary-card__item-label {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.review-summary-card__item-value {
  margin-top: 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.review-summary-card__item-hint {
  margin-top: 6px;
  font-size: 13px;
  line-height: 1.7;
  color: var(--color-text-secondary);
}
</style>